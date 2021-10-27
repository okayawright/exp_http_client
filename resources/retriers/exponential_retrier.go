package retriers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/okayawright/exp_http_client/resources/misc"
)

// default maximum number of times we can resend a request
const defaultMaxTries = 3

// default HTTP error codes that trigger an automatic retry
func defaultRetryableCodes() []int { return []int{429, 500, 503, 504} }

// Do we allow some random timing kew between each retry in order to avoid potential synchronized peaks of requests
const defaultJittering = true

/* Retry sending the request if needed using an exponential back-off timing */
type exponentialRetrier struct {
	/* Number of times, at most, we can try to send the request.*/
	maxTries uint
	/* HTTP status codes that trigger a retry */
	retryableCodes []int
	/* Allow some jitter for each retry timing */
	jitter bool
}

/* Set a maximum number of tries. 1 is the minimum.
Returns the updated resource */
func (retrier *exponentialRetrier) WithMaxTries(maxTries uint) *exponentialRetrier {
	if maxTries > 0 {
		retrier.maxTries = maxTries
	}
	return retrier
}

/* Set the HTTP error codes that can initiate a retry.
Returns the updated resource */
func (retrier *exponentialRetrier) WithRetryableCodes(retryableCodes []int) *exponentialRetrier {
	if retryableCodes != nil && len(retryableCodes) > 0 {
		retrier.retryableCodes = retryableCodes
	}
	return retrier
}

/* Do we allow the time between each retry not to be exact.
Returns the updated resource */
func (retrier *exponentialRetrier) WithJitter(jitter bool) *exponentialRetrier {
	retrier.jitter = jitter
	return retrier
}

/* exponentialRetrier c'tor.
Will try at max. 3 times to send a request, with some jitter between each retry, if the HTTP response error is 429, 500, 503, or 504.*/
func NewExponentialRetrier() *exponentialRetrier {
	return &exponentialRetrier{
		maxTries:       defaultMaxTries,
		retryableCodes: defaultRetryableCodes(),
		jitter:         defaultJittering,
	}
}

/* Try to make an HTTP request with the given client for the specified prepared request.
If unsuccessful, retry after some time if the HTTP error code allows it, or if the client timed out, and we don't go beyond the maximum allowed number of tries yet.
Returns the response if successful, and the actual number of tries */
func (retrier *exponentialRetrier) Try(client misc.HttpClient, request *http.Request) (*http.Response, uint, error) {
	var response *http.Response
	var err error
	var previousDelay int64 = 0
	var try uint
	for try = 1; try <= retrier.maxTries; try++ {

		//Timestamp the call, mandatory
		request.Header.Set("Date", time.Now().UTC().Format(time.RFC1123))

		//Actual HTTP request
		response, err = client.Do(request)
		//Wait and analyse the result
		canRetry := false
		//One can only retry calling the service if the client timed out or if the service returned a compatible HTTP error code
		if err != nil && errors.Is(err, context.DeadlineExceeded) {
			canRetry = true
		} else if response != nil {
			for _, v := range retrier.retryableCodes {
				if int(v) == response.StatusCode {
					canRetry = true
					break
				}
			}
		}
		//If we need to retry then wait with an exponential back-off
		if canRetry {
			delay := int64(math.Floor((math.Pow(2, float64(try)) - 1) * 0.5))
			//How much jitter should we apply?
			//The delay can be be reduced or increased by 25% at most, compared to the expected
			maxJitter := int64(float32(delay-previousDelay) * 0.25)
			var jitter int64
			if maxJitter > 0 {
				rand.Seed(time.Now().UnixNano())
				//Will pick a number within the [-maxJitter,+maxJitter[ interval
				jitter = rand.Int63n(2*maxJitter) - maxJitter
			} else {
				jitter = 0
			}
			previousDelay = delay
			time.Sleep(time.Duration(delay+jitter) * time.Second)
		} else {
			break
		}
	}
	fmt.Printf("err %v\n", err)
	return response, try, err
}
