package misc

import (
	netUrl "net/url"
	"strings"
)

/* Try to match and replace the given named parameters, which are enclosed in curly brackets, within the path, fragment and query of the specified URL.
If none matched we append those parameters as querystrings.
Any remaining unresolved named parameter are not removed in the url and kept as-is.
TODO: handle arrays in querystrings, not just as single-value items
The URL is not modified by side-effect, read the returned URL the see the actual chnages */
func Resolve(url *netUrl.URL, values *map[string]string) *netUrl.URL {
	//Shortcut
	if url == nil {
		return nil
	}

	//Shallow cloning of the url in order not to modify the original one
	copiedUrl := *url

	//Shortcut
	if values == nil || len(*values) == 0 {
		return &copiedUrl
	}

	for k, v := range *values {
		tag := "{" + k + "}"

		//Try to find and replace all occurences of the named parameter within the path, then within the fragment, and finally within the query string
		hasMatched := false
		for _, vv := range [3]*string{&copiedUrl.Path, &copiedUrl.Fragment, &copiedUrl.RawQuery} {
			if strings.Contains(*vv, tag) {
				*vv = strings.ReplaceAll(*vv, tag, v)
				hasMatched = true
			}
		}

		//If we couldn't find the named parameter it means we need to add it as a new query string
		if !hasMatched {
			q := copiedUrl.Query()
			q.Add(k, v)
			copiedUrl.RawQuery = q.Encode()
		}
	}

	return &copiedUrl
}
