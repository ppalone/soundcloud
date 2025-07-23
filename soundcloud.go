package soundcloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	liburl "net/url"
	"strings"
	"sync"

	"github.com/grafov/m3u8"
	"golang.org/x/net/html"
)

// Soundcloud Client.
type Client struct {
	httpClient *http.Client
	clientId   string
}

const (
	webURL = "https://soundcloud.com"
	apiURL = "https://api-v2.soundcloud.com"
)

var (
	ErrScrapingClientId = errors.New("error while scrapping client id")
)

// NewClient returns a new Soundcloud client.
func NewClient(opts ...ClientOption) (*Client, error) {
	options := defaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}

	clientId := strings.TrimSpace(options.clientId)
	if len(clientId) == 0 {
		id, err := scrapClientID(options.httpClient)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrScrapingClientId, err)
		}
		clientId = id
	}

	return &Client{
		httpClient: options.httpClient,
		clientId:   clientId,
	}, nil
}

// SearchTracks
func (c *Client) SearchTracks(ctx context.Context, q string, opts ...SearchOption) (SearchTracksResults, error) {
	options := defaultSearchOptions()
	for _, opt := range opts {
		opt(options)
	}

	return c.searchTracks(ctx, q, options)
}

// GetTrackById
func (c *Client) GetTrackById(ctx context.Context, id int) (Track, error) {
	return c.getTrackById(ctx, id)
}

// GetStream
func (c *Client) GetStream(ctx context.Context, transcoding Transcoding) (io.ReadCloser, error) {
	return c.getStream(ctx, transcoding)
}

// GetStreamById
func (c *Client) GetStreamById(ctx context.Context, id int, opts ...StreamOption) (io.ReadCloser, error) {
	options := defaultStreamOptions()
	for _, opt := range opts {
		opt(options)
	}
	return c.getStreamById(ctx, id, options)
}

func (c *Client) getStream(ctx context.Context, transcoding Transcoding) (io.ReadCloser, error) {
	req, err := c.buildRequest(ctx, strings.TrimPrefix(transcoding.URL, fmt.Sprintf("%s/", apiURL)), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	v := &struct {
		URL string `json:"url"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}

	if len(v.URL) == 0 {
		return nil, fmt.Errorf("url is empty")
	}

	pr, pw := io.Pipe()
	p := transcoding.Format.Protocol
	switch p {
	case HLS.String():
		go c.downloadHLS(ctx, v.URL, pw)
	case PROGRESSIVE.String():
		go c.downloadProgressive(ctx, v.URL, pw)
	default:
		errProtocolNotHandled := fmt.Errorf("protocol not handled: %s", p)
		pw.CloseWithError(errProtocolNotHandled)
		return nil, errProtocolNotHandled
	}

	return pr, nil
}

func (c *Client) getStreamById(ctx context.Context, id int, opts *streamOptions) (io.ReadCloser, error) {
	track, err := c.getTrackById(ctx, id)
	if err != nil {
		return nil, err
	}

	t, ok := findTranscoding(track.Transcodings, opts.preset, opts.protocol)
	if !ok {
		return nil, fmt.Errorf("transcoding with preset %v and protocol %v not found for track", opts.preset.String(), opts.protocol.String())
	}

	return c.getStream(ctx, t)
}

func findTranscoding(transcodings []Transcoding, preset Preset, protocol Protocol) (Transcoding, bool) {
	for _, t := range transcodings {
		if strings.HasPrefix(t.Preset, preset.String()) && t.Format.Protocol == protocol.String() {
			return t, true
		}
	}
	return Transcoding{}, false
}

func (c *Client) downloadHLS(ctx context.Context, url string, pw *io.PipeWriter) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		pw.CloseWithError(err)
		return
	}

	playlistURL, err := liburl.Parse(url)
	if err != nil {
		pw.CloseWithError(err)
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		pw.CloseWithError(err)
		return
	}
	defer resp.Body.Close()

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		pw.CloseWithError(err)
		return
	}

	if listType != m3u8.MEDIA {
		pw.CloseWithError(fmt.Errorf("unexpected list type: %v", listType))
		return
	}

	mediaPlaylist, ok := playlist.(*m3u8.MediaPlaylist)
	if !ok {
		pw.CloseWithError(fmt.Errorf("not a valid media playlist"))
		return
	}

	type result struct {
		idx  int
		data []byte
		err  error
	}

	count := 0
	for _, s := range mediaPlaylist.Segments {
		if s != nil {
			count += 1
		}
	}

	// limit concurrency to 10
	// TODO: make it configurable via options
	limit := make(chan struct{}, 10)
	results := make([]result, count)
	wg := &sync.WaitGroup{}

	for i, seg := range mediaPlaylist.Segments {
		if seg == nil {
			continue
		}

		segURI := seg.URI
		wg.Add(1)
		limit <- struct{}{}
		go func(ctx context.Context, idx int, url string, wg *sync.WaitGroup, limit chan struct{}) {
			defer wg.Done()
			defer func() {
				<-limit
			}()

			u, err := playlistURL.Parse(url)
			if err != nil {
				results[idx] = result{idx, nil, err}
				return
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				results[idx] = result{idx, nil, err}
				return
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				results[idx] = result{idx, nil, err}
				return
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			results[idx] = result{idx, data, err}
		}(ctx, i, segURI, wg, limit)
	}

	// wait
	wg.Wait()

	for _, result := range results {
		if result.err != nil {
			pw.CloseWithError(result.err)
			return
		}

		_, err := pw.Write(result.data)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
	}

	// close
	pw.Close()
}

func (c *Client) downloadProgressive(ctx context.Context, url string, pw *io.PipeWriter) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		pw.CloseWithError(err)
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		pw.CloseWithError(err)
		return
	}
	defer resp.Body.Close()

	// 200 - 207
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultiStatus {
		pw.CloseWithError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return
	}

	_, copyErr := io.Copy(pw, resp.Body)
	if copyErr != nil {
		pw.CloseWithError(copyErr)
		return
	}

	pw.Close()
}

func (c *Client) searchTracks(ctx context.Context, q string, opts *searchOptions) (SearchTracksResults, error) {
	q = strings.TrimSpace(q)
	if len(q) == 0 {
		return SearchTracksResults{}, fmt.Errorf("search query is required")
	}
	opts.q = q

	req, err := c.buildRequest(ctx, "search/tracks", opts.build())
	if err != nil {
		return SearchTracksResults{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchTracksResults{}, err
	}
	defer resp.Body.Close()

	apiResponse := new(searchTracksAPIResponse)
	err = json.NewDecoder(resp.Body).Decode(apiResponse)
	if err != nil {
		return SearchTracksResults{}, err
	}

	return apiResponse.toResults(), nil
}

func (c *Client) getTrackById(ctx context.Context, id int) (Track, error) {
	req, err := c.buildRequest(ctx, fmt.Sprintf("tracks/%d", id), nil)
	if err != nil {
		return Track{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Track{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Track{}, fmt.Errorf("invalid id")
	}

	apiResponse := new(trackAPIResponse)
	err = json.NewDecoder(resp.Body).Decode(apiResponse)
	if err != nil {
		return Track{}, err
	}

	return apiResponse.toTrack(), nil
}

func scrapClientID(httpClient *http.Client) (string, error) {
	resp, err := httpClient.Get(webURL)
	if err != nil {
		return "", fmt.Errorf("error requesting url: %w", err)
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing html: %w", err)
	}

	urls := []string{}
	var traverse func(n *html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "src" && strings.Contains(attr.Val, "a-v2.sndcdn.com/assets") {
					urls = append(urls, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(node)

	if len(urls) == 0 {
		return "", fmt.Errorf("no urls found in response")
	}

	resp, err = httpClient.Get(urls[len(urls)-1])
	if err != nil {
		return "", fmt.Errorf("error requesting url: %w", err)
	}
	defer resp.Body.Close()

	h, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response: %w", err)
	}
	body := string(h)

	if strings.Contains(body, `,client_id:"`) {
		split := strings.Split(body, `,client_id:"`)
		if len(split) > 1 {
			t := strings.Split(split[1], `"`)
			if len(t) > 0 {
				return t[0], nil
			}
		}
	}

	return "", fmt.Errorf("error getting client id")
}

func (c *Client) buildRequest(ctx context.Context, path string, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", apiURL, path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	p := req.URL.Query()
	p.Set("client_id", c.clientId)

	for k, v := range params {
		p.Set(k, v)
	}
	req.URL.RawQuery = p.Encode()

	return req, nil
}
