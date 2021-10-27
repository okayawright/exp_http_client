package retriers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/okayawright/exp_http_client/resources/mocks"
)

/* Nominal case, create an exponential retrier with customized options */
func TestNewCustomizedNominal(t *testing.T) {
	expectedJitter := false
	expectedMaxTries := uint(10)
	expectedRetryableCodes := []int{999, 998}

	observed := NewExponentialRetrier().WithJitter(expectedJitter).WithMaxTries(expectedMaxTries).WithRetryableCodes(expectedRetryableCodes)
	if observed.jitter != expectedJitter {
		t.Errorf("NewExponentialRetrier():jitter = %v, want %v", observed.jitter, expectedJitter)
	}
	if observed.maxTries != expectedMaxTries {
		t.Errorf("NewExponentialRetrier():maxTries = %v, want %v", observed.maxTries, expectedMaxTries)
	}
	if !reflect.DeepEqual(observed.retryableCodes, expectedRetryableCodes) {
		t.Errorf("NewExponentialRetrier():retryableCodes = %v, want %v", observed.retryableCodes, expectedRetryableCodes)
	}
}

/* Nominal case, call a REST API that sends back a 200 */
func TestExponentialRetrierTryNoRetryNominal(t *testing.T) {

	expectedNumberOfTries := uint(1)
	serviceNominalStatusCode := 200

	mockClient := mocks.Client{}
	// Send back a fake response with just a 200 HTTP code
	mockClient.MockedDo = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: serviceNominalStatusCode,
			Body:       io.NopCloser(strings.NewReader("{}")),
		}, nil
	}
	req, _ := http.NewRequest("GET", "http://nowhere", nil)
	_, n, err := NewExponentialRetrier().Try(&mockClient, req)
	if err != nil {
		t.Fatalf("Try() unexpected error %v", err)
	}
	if n != expectedNumberOfTries {
		t.Errorf("Try() number of tries = %v, want %v", n, expectedNumberOfTries)
	}
}

/* Nominal case, call a REST API that temporarily sends back an error 429 */
func TestExponentialRetrierTryBackPressureNominal(t *testing.T) {

	expectedNumberOfTries := uint(4)
	serviceErrorStatusCode := 500
	serviceNominalStatusCode := 200

	mockClient := mocks.Client{}
	// Send back a fake response with just a 429 HTTP code for a while, then a 200
	errorPass := uint(0)
	mockClient.MockedDo = func(req *http.Request) (*http.Response, error) {
		errorPass++
		var serviceStatusCode int
		if errorPass < expectedNumberOfTries {
			serviceStatusCode = serviceErrorStatusCode
		} else {
			serviceStatusCode = serviceNominalStatusCode
		}
		return &http.Response{
			StatusCode: serviceStatusCode,
			Body:       io.NopCloser(strings.NewReader("{}")),
		}, nil
	}
	req, _ := http.NewRequest("GET", "http://nowhere", nil)
	_, n, err := NewExponentialRetrier().WithMaxTries(4).Try(&mockClient, req)
	if err != nil {
		t.Fatalf("Try() unexpected error %v", err)
	}
	if n != expectedNumberOfTries {
		t.Errorf("Try() number of tries = %v, want %v", n, expectedNumberOfTries)
	}
}

/* Nominal case, call a REST API that temporarily make the client times out */
func TestExponentialRetrierTryTimeoutNominal(t *testing.T) {

	expectedNumberOfTries := uint(3)
	serviceTimeoutError := context.DeadlineExceeded
	serviceNominalStatusCode := 200

	mockClient := mocks.Client{}
	// Send back a timeout error code for a while, then operate back to normal
	errorPass := uint(0)
	mockClient.MockedDo = func(req *http.Request) (*http.Response, error) {
		errorPass++
		var serviceStatusError error
		if errorPass < expectedNumberOfTries {
			serviceStatusError = serviceTimeoutError
		} else {
			serviceStatusError = nil
		}
		return &http.Response{
			StatusCode: serviceNominalStatusCode,
			Body:       io.NopCloser(strings.NewReader("{}")),
		}, serviceStatusError
	}
	req, _ := http.NewRequest("GET", "http://nowhere", nil)
	_, n, err := NewExponentialRetrier().Try(&mockClient, req)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Try() unexpected error %v", err)
	}
	if n != expectedNumberOfTries {
		t.Errorf("Try() number of tries = %v, want %v", n, expectedNumberOfTries)
	}
}
