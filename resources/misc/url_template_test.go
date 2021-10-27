package misc

import (
	netUrl "net/url"
	"reflect"
	"testing"
)

/* Nominal case, features existing named parameters to replace and ones to add*/
func TestReplaceNominal(t *testing.T) {
	expected, _ := netUrl.Parse("http://localhost:8080/api/julien/info?withBio=true&withCredentials=true&withPhoto=false")
	u, _ := netUrl.Parse("http://localhost:8080/api/{user}/info?withCredentials={credentials}&withPhoto=false")
	observed := Resolve(u, &map[string]string{
		"user":        "julien",
		"credentials": "true",
		"withBio":     "true",
	})
	if !reflect.DeepEqual(observed.String(), expected.String()) {
		t.Errorf("Resolve() = %v, want %v", observed.String(), expected.String())
	}
}

/* Nominal case, url without any specified named parameter */
func TestReplaceNoReplacement(t *testing.T) {
	expected, _ := netUrl.Parse("http://localhost:8080/api/julien/info?withBio={withBio}&withCredentials=true&withPhoto=false")
	u, _ := netUrl.Parse("http://localhost:8080/api/julien/info?withBio={withBio}&withCredentials=true&withPhoto=false")
	observed := Resolve(u, nil)
	if !reflect.DeepEqual(observed.String(), expected.String()) {
		t.Errorf("Resolve() = %v, want %v", observed.String(), expected.String())
	}
}

/* Corner case, properly re-encode urls with prior escaped characters */
func TestReplaceProperEscaping(t *testing.T) {
	expected, _ := netUrl.Parse("http://localhost:8080/api/L%27orec/info?login=OAuth2")
	u, _ := netUrl.Parse("http://localhost:8080/api/{user}/info?login=OAuth2")
	observed := Resolve(u, &map[string]string{
		"user": "L'orec",
	})
	if !reflect.DeepEqual(observed.String(), expected.String()) {
		t.Errorf("Resolve() = %v, want %v", observed.String(), expected.String())
	}
}
