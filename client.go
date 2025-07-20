package soundcloud

import "net/http"

// client options.
type clientOptions struct {
	httpClient *http.Client
	clientId   string
}

type ClientOption func(o *clientOptions)

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		httpClient: &http.Client{},
		clientId:   "",
	}
}

func WithHTTPClient(c *http.Client) ClientOption {
	return func(o *clientOptions) {
		o.httpClient = c
	}
}

func WithClientID(clientId string) ClientOption {
	return func(o *clientOptions) {
		o.clientId = clientId
	}
}
