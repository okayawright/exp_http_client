/* Integration tests for the exp_http_client module.
Private mock API not provided */
package example_user_test

import (
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/okayawright/exp_http_client/example_user"
)

var userId = "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

/* Nominal case, integration test for Create then Fetch and finally Delete for the same User.
In order to be able to leave the database clean after the test one is encouraged to run the three together in sequence
TODO don't use the service to reset the application nominal state, directly tap into the database, so that you can test create and delete independently */
func TestCreateFetchDeleteUserNominal(t *testing.T) {
	version := 0
	user := example_user.Data{
		Data: &example_user.User{
			ID:        userId,
			FirstName: "first",
			LastName:  "last",
		},
	}

	//Create
	observedResponseData, observedStatusCode, err := example_user.Create("http://sampleapi:8080/", &user)
	if err != nil {
		t.Fatalf("TestCreateFetchDeleteUserNominal:Create unexpected error %v", err)
	}
	if observedStatusCode != 201 {
		t.Errorf("TestCreateFetchDeleteUserNominal:Create:observedStatusCode = %v, want %v", observedStatusCode, 201)
	}
	if !reflect.DeepEqual(*observedResponseData, user) {
		t.Errorf("TestCreateFetchDeleteUserNominal:Create:observedResponseData = %v, want %v", *observedResponseData, user)
	}

	//Fetch
	observedResponseData, observedStatusCode, err = example_user.Fetch("http://sampleapi:8080/", userId)
	if err != nil {
		t.Fatalf("TestCreateFetchDeleteUserNominal:Fetch unexpected error %v", err)
	}
	if observedStatusCode != 200 {
		t.Errorf("TestCreateFetchDeleteUserNominal:Fetch:observedStatusCode = %v, want %v", observedStatusCode, 200)
	}
	if !reflect.DeepEqual(*observedResponseData, user) {
		t.Errorf("TestCreateFetchDeleteUserNominal:Fetch:observedResponseData = %v, want %v", *observedResponseData, user)
	}

	//Delete
	observedResponseData, observedStatusCode, err = example_user.Delete("http://sampleapi:8080/", userId, version)
	if err != nil {
		t.Fatalf("TestCreateFetchDeleteUserNominal:Delete unexpected error %v", err)
	}
	if observedStatusCode != 204 {
		t.Errorf("TestCreateFetchDeleteUserNominal:Delete:observedStatusCode = %v, want %v", observedStatusCode, 204)
	}
	empty := example_user.Data{}
	if !reflect.DeepEqual(*observedResponseData, empty) {
		t.Errorf("TestCreateFetchDeleteUserNominal:Fetch:observedResponseData = %v, want %v", *observedResponseData, empty)
	}

}

/* Error case, empty input for Create */
func TestCreateEmptyInputError(t *testing.T) {
	_, observedStatusCode, err := example_user.Create("http://sampleapi:8080/", nil)
	if err != nil {
		t.Fatalf("TestCreateEmptyInputError unexpected error %v", err)
	}
	if observedStatusCode != 500 {
		t.Errorf("TestCreateEmptyInputError:observedStatusCode = %v, want %v", observedStatusCode, 500)
	}
}

/* Error case, invalid input for Create */
func TestCreateInvalidInputError(t *testing.T) {
	invalidUser := example_user.Data{
		Data: &example_user.User{},
	}

	_, observedStatusCode, err := example_user.Create("http://sampleapi:8080/", &invalidUser)
	if err != nil {
		t.Fatalf("TestCreateInvalidInputError unexpected error %v", err)
	}
	if observedStatusCode != 400 {
		t.Errorf("TestCreateInvalidInputError:observedStatusCode = %v, want %v", observedStatusCode, 400)
	}
}

/* Error case, invalid input for Fetch */
func TestFetchInvalidInputError(t *testing.T) {
	_, observedStatusCode, err := example_user.Fetch("http://sampleapi:8080/", "InvalidData")
	if err != nil {
		t.Fatalf("TestFetchInvalidInputError unexpected error %v", err)
	}
	if observedStatusCode != 400 {
		t.Errorf("TestFetchInvalidInputError:observedStatusCode = %v, want %v", observedStatusCode, 400)
	}
}

/* Error case, non existing data for Fetch */
func TestFetchNotExistError(t *testing.T) {
	_, observedStatusCode, err := example_user.Fetch("http://sampleapi:8080/", "ad27e265-9605-4b4b-a0e5-3003ea9cc4dd")
	if err != nil {
		t.Errorf("TestFetchNotExistError unexpected error %v", err)
	}
	if observedStatusCode != 404 {
		t.Errorf("TestFetchNotExistError:observedStatusCode = %v, want %v", observedStatusCode, 404)
	}
}

/* Error case, invalid URL */
func TestFetchInvalidUrlError(t *testing.T) {
	_, observedStatusCode, err := example_user.Fetch("http://outside:9999/", "ad27e265-9605-4b4b-a0e5-3003ea9cc4dd")
	if err == nil {
		t.Fatalf("TestFetchInvalidUrlError unexpected success")
	}
	dnsErr := new(net.DNSError)
	if !errors.As(err, &dnsErr) {
		t.Errorf("TestFetchInvalidUrlError unexpected error type %T", err)
	}
	if observedStatusCode != 0 {
		t.Errorf("TestFetchInvalidUrlError:observedStatusCode = %v, want %v", observedStatusCode, 0)
	}
}
