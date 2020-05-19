package quick

import (
	"net/http"
	"net/url"
	"strings"
)

// request option func
type OptionFunc func(*Request)

// set request header
func HeaderOption(v interface{}) OptionFunc {
	hd := make(http.Header)

	switch v.(type) {
	case http.Header:
		hd = v.(http.Header)
	case map[string]string:
		dict := v.(map[string]string)
		for k, v := range dict {
			hd.Set(k, v)
		}
	case []string:
		l := v.([]string)
		for _, v := range l {
			arr := strings.Split(v, ":")
			if len(arr) == 2 {
				hd.Set(strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1]))
			} else if len(arr) == 1 {
				hd.Set(strings.TrimSpace(arr[0]), "")
			}
		}
	default:
		panic("Header: parameter types are not supported")
	}

	return func(req *Request) {
		req.SetHeader(hd)
	}
}

// request query string for get
func QueryStringOption(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetQueryString(v)
	}
}

// request body for post
func BodyOption(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBody(v)
	}
}

// request body for post (FormData)
func BodyFormDataOption(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBodyFormData(v)
	}
}

// set proxy for request
func ProxyOption(v interface{}) OptionFunc {
	var rawurl string
	switch v.(type) {
	case string:
		rawurl = v.(string)
	case *url.URL:
		u := v.(*url.URL)
		rawurl = u.String()
	default:
		panic("Proxy: parameter types are not supported")
	}
	return func(req *Request) {
		req.SetProxy(rawurl)
	}
}