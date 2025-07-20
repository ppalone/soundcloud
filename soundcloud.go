package soundcloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

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
