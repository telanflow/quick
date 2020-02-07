package requests

import (
	"net/http"
	"net/url"
	"strings"
)

// mergeHeaders merge Request headers and Session Headers.
// Request has higher priority.
func MergeHeaders(h1, h2 http.Header) http.Header {
	h := http.Header{}
	for key, values := range h1 {
		for _, value := range values {
			h.Set(key, value)
		}
	}
	for key, values := range h2 {
		for _, value := range values {
			h.Set(key, value)
		}
	}
	return h
}

// copy headers
func CopyHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// Get request merge url and query string encode.
func MergeQueryParams(parsedURL *url.URL, parsedQuery string) (*url.URL, error) {
	rawurl := strings.Join([]string{strings.Replace(parsedURL.String(), "?"+parsedURL.RawQuery, "", -1), parsedQuery}, "?")
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// get first interface
func getParam(params ...interface{}) interface{} {
	if params == nil || len(params) == 0 {
		return nil
	}
	return params[0]
}
