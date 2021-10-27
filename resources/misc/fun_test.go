package misc

import (
	"testing"
)

/* Nominal case, finds an exact specific string */
func TestFindNominalExact(t *testing.T) {
	list := []string{
		"a",
		"b",
		"c",
	}
	observed := Find(list, "b", false)
	expected := 1
	if observed != expected {
		t.Errorf("Find() = %v, want %v", observed, expected)
	}
}

/* Nominal case, finds a partially matching string */
func TestFindNominalPartial(t *testing.T) {
	list := []string{
		"abc",
		"bcd",
		"cde",
	}
	observed := Find(list, "d", true)
	expected := 1
	if observed != expected {
		t.Errorf("Find() = %v, want %v", observed, expected)
	}
	observed = Find(list, "de", true)
	expected = 2
	if observed != expected {
		t.Errorf("Find() = %v, want %v", observed, expected)
	}
}
