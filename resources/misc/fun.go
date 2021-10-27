package misc

import "strings"

type compFunc func(bait string, probe string) bool

/* Find the first element index matching needle in a slice haystack.
If partial is true then needle can only be a substring of an element of haystack to be considered a match.
A result of -1 means no match */
func Find(haystack []string, needle string, partial bool) int {
	//Build the comparison function based on the partial argument
	var call compFunc
	if partial {
		call = func(bait string, probe string) bool {
			//Fuzzy
			return strings.Contains(bait, needle)
		}
	} else {
		call = func(bait string, probe string) bool {
			return bait == probe
		}
	}
	//Find; linear search
	for i, n := range haystack {
		if call(n, needle) {
			return i
		}
	}
	return -1
}
