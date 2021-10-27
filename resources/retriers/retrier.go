package retriers

import (
	"net/http"

	"github.com/okayawright/exp_http_client/resources/misc"
)

/* Retry sending the request with the given client if needed using an implementation-specific frequency until it succeeds or a specific condition is met */
type Retrier interface {
	Try(client misc.HttpClient, request *http.Request) (*http.Response, uint, error)
}
