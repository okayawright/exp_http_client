package misc

import "net/http"

/* Interface for both the http.Client and the mocks.Client */
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
