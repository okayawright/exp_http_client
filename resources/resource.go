package resources

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"strings"
	"time"

	"github.com/okayawright/exp_http_client/resources/misc"
	"github.com/okayawright/exp_http_client/resources/retriers"
	"github.com/okayawright/exp_http_client/resources/serializers"
)

// default request timeout if unspecified, in seconds
const defaultTimeout = 30

/* A NewJsonMarshaller() is a re-usable and configurable HTTP REST client tied to a specific endpoint and a serializer */
type resource struct {
	//HTTP client engine
	client misc.HttpClient
	//Request endpoint
	endpoint *netUrl.URL
	//Request and response (un)marshaller
	marshaller serializers.Marshaller
	//Retry handler
	retrier retriers.Retrier
	//Request timeout in seconds, 0 means no limit
	timeout uint
}

/* Make an HTTP request for a prepared Request.
Returns an instance of the struct we expect augmented with the response data, and the HTTP status code (0 means unknown) */
type CallFunc func() (interface{}, int, error)

/* Use a specific marshaller/unmarshaller for this resource.
Returns the updated resource */
func (resource *resource) WithMarshaller(marshaller serializers.Marshaller) *resource {
	if marshaller != nil {
		resource.marshaller = marshaller
	}
	return resource
}

/* Use a specific marshaller/unmarshaller for this resource.
Returns the updated resource */
func (resource *resource) WithTimeout(timeout uint) *resource {
	resource.timeout = timeout
	return resource
}

/* Use a specific HTTP client for this resource.
Returns the updated resource */
func (resource *resource) WithClient(client misc.HttpClient) *resource {
	if client != nil {
		resource.client = client
	}
	return resource
}

/* Use a specific retrying implementation for this resource.
Returns the updated resource */
func (resource *resource) WithRetrier(retrier retriers.Retrier) *resource {
	if retrier != nil {
		resource.retrier = retrier
	}
	return resource
}

/* resource constructor.
url is a mandatory parameterized URL template with parameters with the path, querystring or fragment enclosed between curly braces.
By default, the marshaller can read and write JSON, the HTTP client is http.DefaultClient, and some selected failed requests will be retried using an exponentila backoff.
Returns the newly built resource
*/
func NewResource(endpoint *netUrl.URL) *resource {
	return &resource{
		//Use the default HTTP client, can be replaced afterward; beware request timeouts will be handled through the context not within the client options
		client:   http.DefaultClient,
		endpoint: endpoint,
		//Use the default marshaller, can be replaced afterward
		marshaller: serializers.NewJsonMarshaller(),
		//use the default retrier, can be replaced afterward
		retrier: retriers.NewExponentialRetrier(),
		//Use a default timeout, can be replaced afterward
		timeout: defaultTimeout,
	}
}

/* encode the body if needed, using the provided marshaller.
Returns the encoding MIME type, and the actual serialized body */
func encodeRequestBody(body interface{}, marshaller serializers.Marshaller) (string, io.Reader, error) {
	if body != nil {
		//TODO right now we only handle JSON structured REST APIs
		encodedBody, err := marshaller.Serialize(body)
		if err != nil {
			return "", nil, err
		}
		return marshaller.SerializationCompatibleMimetype(), bytes.NewBuffer(encodedBody), nil
	} else {
		return "", nil, nil
	}
}

/*
Prepare a request for a given action, with optional values for named parameters within the URL and the body struct if required.
actionName is the case-sensitive name of a registered action on this resource, an undefined action is fatal,
urlParameters is an optional set of named parameters values to replace within the url to call,
body is the optional body to send in the request, only meaningful for verbs that usually send request bodies (e.g. POST, PUT, PATCH).
Returns a function to make the actual HTTP call, and a request cancelling function that can be used to abort the execution of the first returned function
*/
func (resource *resource) Request(verb string, urlParameters *map[string]string, body interface{}) (CallFunc, context.CancelFunc, error) {

	//Clone a new pristine context in order to control the request once sent
	//and make the request cancellable and expirable
	actualContext, cancel := context.WithTimeout(context.Background(), time.Duration(resource.timeout)*time.Second)

	//Resolve the template URL if needed
	url := misc.Resolve(resource.endpoint, urlParameters)
	//Shouldn't happen, panic
	if url == nil {
		panic("The endpoint to query cannot be nil")
	}

	//Prepare the body, if needed
	contentType, encodedBodyReader, err := encodeRequestBody(body, resource.marshaller)
	if err != nil {
		return nil, cancel, err
	}

	request, err := http.NewRequestWithContext(actualContext, strings.ToUpper(verb), url.String(), encodedBodyReader)
	//Rethrow without doing anything if an error occurs, we do not know what to do with it right here
	if err != nil {
		return nil, cancel, err
	}

	//Augment the request with metadata if needed
	if len(contentType) > 0 {
		request.Header.Set("Content-Type", contentType)
	}
	//TODO right now we only handle JSON structured REST APIs
	request.Header.Set("Accept", strings.Join(resource.marshaller.DeserializationCompatibleMimetypes(), ","))

	return func() (interface{}, int, error) {
		return call(resource.client, request, &resource.marshaller, &resource.retrier)
	}, cancel, nil

}

/* decode the body if needed, using thcalle provided marshaller, if compatible.
contentTypes are the optional MIME types returned in the response and are verified for compatibility,
Returns the actual serialized body */
func decodeResponseBody(rawBody io.Reader, contentTypes []string, marshaller serializers.Marshaller) (interface{}, error) {
	responseBody, err := ioutil.ReadAll(rawBody)
	if err != nil {
		return nil, err
	}

	//Just to be sure, we should check whether the API sent us back a format we understand
	if contentTypes != nil && len(contentTypes) > 0 {
		match := false
		for _, s := range marshaller.DeserializationCompatibleMimetypes() {
			if misc.Find(contentTypes, s, true) != -1 {
				match = true
				break
			}
		}
		if !match {
			return nil, errors.New("Cannot deserialize the response")
		}
	}

	//Transform the binary body into a structured map, if not empty
	if len(responseBody) > 0 {
		return marshaller.Deserialize(responseBody)
	} else {
		return nil, nil
	}
}

/* Make an HTTP request with the given client for the specified prepared request.
Returns the structured map corresponding to the response body, the HTTP status code, 0 means we don't have one to provide */
func call(client misc.HttpClient, request *http.Request, marshaller *serializers.Marshaller, retrier *retriers.Retrier) (interface{}, int, error) {

	//Actual HTTP request
	response, _, err := (*retrier).Try(client, request)
	if err != nil {
		return nil, 0, err
	}

	//Deserialize the response body
	//In order to reuse a TCP connection from the pool you need to read the response body till the end @see https://golang.cafe/blog/how-to-reuse-http-connections-in-go.html
	defer response.Body.Close()
	bodyStruct, err := decodeResponseBody(response.Body, response.Header["Content-Type"], *marshaller)

	return bodyStruct, response.StatusCode, err

}
