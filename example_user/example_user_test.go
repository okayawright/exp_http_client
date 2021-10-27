package example_user

import (
	"reflect"
	"testing"
)

/* Nominal case, replace a given host into a given URL */
func TestMergeUrlHostNominal(t *testing.T) {

	expected := "https://newhost.com/api/v1/info/user1?withPhoto=true"

	newHost := "https://newhost.com/"
	defaultUrl := "http://localhost:8080/api/v1/info/user1?withPhoto=true"
	mergedUrl := mergeUrlHost(newHost, defaultUrl)

	if mergedUrl.String() != expected {
		t.Errorf("mergeUrlHost() = %v, want %v", mergedUrl.String(), expected)
	}

}

/* Nominal case, transform a map into a struct */
func TestConvertResponseNominal(t *testing.T) {
	expected := Data{
		Data: &User{
			ID:        "012",
			FirstName: "first",
			LastName:  "last",
		},
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         "012",
			"first_name": "first",
			"last_name":  "last",
		},
	}
	observed, err := convertResponse(data)
	if err != nil {
		t.Fatalf("ConvertResponse() unexpected error %v", err)
	}
	if !reflect.DeepEqual(*observed, expected) {
		t.Errorf("ConvertResponse() = %v, want %v", *observed, expected)
	}
}
