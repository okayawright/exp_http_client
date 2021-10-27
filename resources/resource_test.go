package resources

import (
	"bytes"
	"io"
	"net/http"
	netUrl "net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/okayawright/exp_http_client/resources/mocks"
	"github.com/okayawright/exp_http_client/resources/serializers"
)

/* Nominal case, create a resource with a customized timeout */
func TestNewCustomizedNominal(t *testing.T) {
	url, _ := netUrl.Parse("http://localhost:8080/api/julien/info")
	var timeout uint = 10

	observed := NewResource(url).WithTimeout(timeout)
	if url.String() != observed.endpoint.String() {
		t.Errorf("NewResource():endpoint = %v, want %v", observed.endpoint.String(), url)
	}
	if observed.timeout != timeout {
		t.Errorf("NewResource():timeout = %v, want %v", observed.timeout, timeout)
	}
}

/* Nominal case, no error, body correctly serialized and metadata added */
func TestRequestNominal(t *testing.T) {
	type secret struct {
		Token string `json:"token"`
	}
	type body struct {
		Data secret `json:"data"`
	}
	url, _ := netUrl.Parse("http://localhost:8080/api/julien/info")
	res := NewResource(url)
	save := body{
		Data: secret{
			Token: "zufeb5e1b6e1b6eb",
		},
	}

	call, cancel, err := res.Request("POST", nil, &save)

	if err != nil {
		t.Fatalf("Prepare() unexpected error %v", err)
	}
	if call == nil {
		t.Errorf("Prepare() unexpected nil call")
	}
	if cancel == nil {
		t.Errorf("Prepare() unexpected nil cancel")
	}

}

/* Nominal case, if there's a body encode it */
func TestResourceEncodeRequestBodyNominal(t *testing.T) {
	type secret struct {
		Token string `json:"token"`
	}
	type body struct {
		Data secret `json:"data"`
	}
	save := body{
		Data: secret{
			Token: "zufeb5e1b6e1b6eb",
		},
	}

	jsonMarshaller := serializers.NewJsonMarshaller()
	expectedMimeType := jsonMarshaller.SerializationCompatibleMimetype()
	expectedEncodedBody, _ := jsonMarshaller.Serialize(save)
	expectedRawEncodedBody := bytes.NewBuffer(expectedEncodedBody)

	observedMimeType, observedRawEncodedBody, err := encodeRequestBody(save, jsonMarshaller)
	if err != nil {
		t.Fatalf("encodeBody() unexpected error %v", err)
	}
	if expectedMimeType != observedMimeType {
		t.Errorf("encodeBody() = %v, want %v", observedMimeType, expectedMimeType)
	}
	if !reflect.DeepEqual(expectedRawEncodedBody, observedRawEncodedBody) {
		t.Errorf("encodeBody() = %v, want %v", observedRawEncodedBody, expectedRawEncodedBody)
	}
}

/* Nominal case, if there's a body decode it */
func TestResourceDecodeResponseBodyNominal(t *testing.T) {
	type secret struct {
		Token string `json:"token"`
	}
	type body struct {
		Data secret `json:"data"`
	}
	expectedBody := body{
		Data: secret{
			Token: "zufeb5e1b6e1b6eb",
		},
	}

	jsonMarshaller := serializers.NewJsonMarshaller()
	mimeTypes := []string{jsonMarshaller.SerializationCompatibleMimetype()}
	encodedBody, _ := jsonMarshaller.Serialize(expectedBody)
	rawEncodedBody := bytes.NewBuffer(encodedBody)

	decodedBody, err := decodeResponseBody(rawEncodedBody, mimeTypes, jsonMarshaller)
	if err != nil {
		t.Fatalf("encodeBody() unexpected error %v", err)
	}
	//Convert the interface{} into a proper map for reflect.DeepEqual() to properly work
	var decodedBodyStruct body
	err = mapstructure.Decode(decodedBody, &decodedBodyStruct)
	if err != nil {
		t.Fatalf("encodeBody() unexpected error %v", err)
	}
	if !reflect.DeepEqual(decodedBodyStruct, expectedBody) {
		t.Errorf("encodeBody() = %v, want %v", decodedBodyStruct, expectedBody)
	}
}

/* Nominal case, call the REST API and get back a 200 */
func TestResourceCallNominal(t *testing.T) {
	type secret struct {
		Token string `json:"token"`
	}
	type body struct {
		Data secret `json:"data"`
	}
	expectedStatusCode := 200

	url, _ := netUrl.Parse("http://localhost:8080/api/julien/info")
	mockClient := mocks.Client{}
	// Send back a fake response with just a 200 HTTP code
	mockClient.MockedDo = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: expectedStatusCode,
			Body:       io.NopCloser(strings.NewReader("{}")),
		}, nil
	}
	res := NewResource(url).WithClient(&mockClient).WithTimeout(60)

	call, _, err := res.Request("GET", nil, nil)
	if err != nil {
		t.Fatalf("Call() unexpected error %v", err)
	}

	_, observedStatusCode, err := call()
	if err != nil {
		t.Fatalf("Call() unexpected error %v", err)
	}
	if observedStatusCode != expectedStatusCode {
		t.Errorf("Call() = %v, want %v", observedStatusCode, expectedStatusCode)
	}
}

/* Nominal case, cancel a long request */
func TestResourceCancelNominal(t *testing.T) {
	type secret struct {
		Token string `json:"token"`
	}
	type body struct {
		Data secret `json:"data"`
	}
	url, _ := netUrl.Parse("http://localhost:8080/api/julien/info")
	mockClient := mocks.Client{}
	// Wait for 10 long seconds before sending back a dummy response
	mockClient.MockedDo = func(req *http.Request) (*http.Response, error) {
		time.Sleep(10 * time.Second)
		return &http.Response{}, nil
	}
	res := NewResource(url).WithClient(&mockClient).WithTimeout(60)

	call, cancel, err := res.Request("GET", nil, nil)
	if err != nil {
		t.Fatalf("Call() unexpected error %v", err)
	}

	go func() {
		call()
		//The call above should be stopped by the cancel order below in the mean time, before a response is produced
		t.Errorf("Call(), was not expecting the call to go through")
	}()

	time.Sleep(time.Second * 1)
	cancel()
	//If we're it means the request has been successfully cancelled on time
}
