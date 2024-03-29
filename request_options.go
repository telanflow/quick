package quick

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OptionFunc request option func
type OptionFunc func(*Request)

// OptionHeader set request header
func OptionHeader(v interface{}) OptionFunc {
	hd := make(http.Header)

	switch t := v.(type) {
	case http.Header:
		hd = t
	case map[string]string:
		for k, v := range t {
			hd.Set(k, v)
		}
	case []string:
		for _, v := range t {
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

// OptionHeaderSingle set an http header to request
func OptionHeaderSingle(k, v string) OptionFunc {
	return func(req *Request) {
		req.SetHeaderSingle(k, v)
	}
}

// OptionQueryString request query string for get
func OptionQueryString(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetQueryString(v)
	}
}

// OptionBody request body for post
func OptionBody(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBody(v)
	}
}

// OptionBodyFormData request body for post (FormData)
func OptionBodyFormData(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBodyFormData(v)
	}
}

// OptionBodyXWwwFormUrlencoded set request body x-www-form-urlencoded
func OptionBodyXWwwFormUrlencoded(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBodyXWwwFormUrlencoded(v)
	}
}

// OptionBodyJSON request body for post
func OptionBodyJSON(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBodyJson(v)
	}
}

// OptionBodyXML request body for post
func OptionBodyXML(v interface{}) OptionFunc {
	return func(req *Request) {
		req.SetBodyXML(v)
	}
}

// OptionBasicAuth HTTP Basic Authentication
func OptionBasicAuth(username, password string) OptionFunc {
	return func(req *Request) {
		req.SetBasicAuth(username, password)
	}
}

// OptionProxy set proxy for request
func OptionProxy(v interface{}) OptionFunc {
	switch t := v.(type) {
	case string:
		return func(req *Request) {
			req.SetProxyUrl(t)
		}
	case *url.URL:
		return func(req *Request) {
			req.Proxy = t
		}
	case url.URL:
		return func(req *Request) {
			req.Proxy = &t
		}
	default:
		panic("Proxy: parameter types are not supported")
	}
}

// OptionTimeout set timeout to request
func OptionTimeout(v time.Duration) OptionFunc {
	return func(req *Request) {
		req.Timeout = v
	}
}

// OptionRedirectNum set redirect num to request
func OptionRedirectNum(num int) OptionFunc {
	return func(req *Request) {
		req.RedirectNum = num
	}
}

// OptionCookies set cookies to request
func OptionCookies(cookies Cookies) OptionFunc {
	return func(req *Request) {
		req.Cookies = cookies
	}
}
