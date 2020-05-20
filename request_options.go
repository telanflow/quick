package quick

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

// request option func
type OptionFunc func(*Request)

// set request header
func OptionHeader(v interface{}) OptionFunc {
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

// set an http header to request
func OptionHeaderSingle(k, v string) OptionFunc {
	return func(req *Request) {
		req.SetHeaderSingle(k, v)
	}
}

// request query string for get
func OptionQueryString(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetQueryString(v)
	}
}

// request body for post
func OptionBody(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBody(v)
	}
}

// request body for post (FormData)
func OptionBodyFormData(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBodyFormData(v)
	}
}

// set proxy for request
func OptionProxy(v interface{}) OptionFunc {
	switch v.(type) {
	case string:
		rawurl := v.(string)
		return func(req *Request) {
			req.SetProxy(rawurl)
		}
	case *url.URL:
		u := v.(*url.URL)
		return func(req *Request) {
			req.Proxy = u
		}
	case url.URL:
		u := v.(url.URL)
		return func(req *Request) {
			req.Proxy = &u
		}
	default:
		panic("Proxy: parameter types are not supported")
	}
}

// set timeout to request
func OptionTimeout(v time.Duration) OptionFunc {
	return func(req *Request) {
		req.Timeout = v
	}
}

// set redirect num to request
func OptionRedirectNum(num int) OptionFunc {
	return func(req *Request) {
		req.RedirectNum = num
	}
}

// set cookies to request
func OptionCookies(cookies Cookies) OptionFunc {
	return func(req *Request) {
		req.Cookies = cookies
	}
}
