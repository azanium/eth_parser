package httpclient

import "net/http"

// This is here, so I can easily switch mock http.Client with real one.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
