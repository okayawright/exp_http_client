package mocks

import "net/http"

// A bare bone mock HTTP client
type Client struct {
	MockedDo func(req *http.Request) (*http.Response, error)
}

// Implement with the mocked operation
//var MockedDo func(req *http.Request) (*http.Response, error)

func (m *Client) Do(req *http.Request) (*http.Response, error) {
	return m.MockedDo(req)
}
