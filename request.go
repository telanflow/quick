package quick

import (
	"bytes"
	"encoding/xml"
	"github.com/telanflow/quick/encode"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

// Sequence number is incremented and utilized for all request created.
var sequenceNo uint64

// http request payload
type Request struct {
	Id          uint64
	Url         string
	Method      string
	Header      http.Header   // request headers
	Body        io.Reader     // request encode
	RedirectNum int           // Number of redirects requested. default 5
	Timeout     time.Duration // request timeout
	Proxy       *url.URL      // request proxy url
	Cookies     Cookies       // request cookies

	host        string // customize the request Host field
	queryString string // customize the GET request query string
}

// create a request instance
func NewRequest() *Request {
	return &Request{
		Id:          atomic.AddUint64(&sequenceNo, 1),
		Url:         "",
		Method:      http.MethodGet,
		Header:      make(http.Header),
		Body:        nil,
		RedirectNum: 10, // set request redirect num. default 10.
		Timeout:     30 * time.Second,
		Proxy:       nil,
		Cookies:     nil,
	}
}

// set request url
func (req *Request) SetUrl(rawurl string) *Request {
	req.Url = rawurl
	return req
}

// get request url
func (req *Request) GetUrl() string {
	return req.Url
}

func (req *Request) SetMethod(method string) *Request {
	req.Method = strings.ToUpper(method)
	return req
}

func (req *Request) GetMethod() string {
	return req.Method
}

// Custom request host field
// GET /index HTTP/1.1
// Host: domain
// ....
func (req *Request) SetHost(host string) *Request {
	req.host = host
	return req
}

func (req *Request) SetQueryString(params interface{}) *Request {
	buff := new(bytes.Buffer)
	form := new(encode.XWwwFormUrlencoded)
	form.SetValue(params)
	if err := form.Encode(buff); err != nil {
		panic(err)
	}

	rawPointUrl, err := url.Parse(req.Url)
	if err != nil {
		panic(err)
	}

	// format get request parameters
	u, err := MergeQueryParams(rawPointUrl, buff.String())
	if err != nil {
		panic(err)
	}

	req.Url = u.String()
	req.Body = nil
	return req
}

func (req *Request) SetBody(params interface{}) *Request {
	buff := new(bytes.Buffer)
	form := new(encode.XWwwFormUrlencoded)
	form.SetValue(params)
	if err := form.Encode(buff); err != nil {
		panic(err)
	}
	req.Body = buff
	return req
}

func (req *Request) SetBodyFormData(params interface{}) *Request {
	buff := new(bytes.Buffer)
	form := new(encode.XWwwFormUrlencoded)
	form.SetValue(params)
	if err := form.Encode(buff); err != nil {
		panic(err)
	}
	req.SetHeaderSingle("Content-Type", "application/form-data")
	req.Body = buff
	return req
}

func (req *Request) SetBodyJson(params interface{}) *Request {
	buff, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	req.SetHeaderSingle("Content-Type", "application/json")
	req.Body = bytes.NewReader(buff)
	return req
}

func (req *Request) SetBodyXml(params interface{}) *Request {
	buff, err := xml.Marshal(params)
	if err != nil {
		panic(err)
	}
	req.SetHeaderSingle("Content-Type", "application/xml")
	req.Body = bytes.NewReader(buff)
	return req
}

func (req *Request) SetBodyXWwwFormUrlencoded(params interface{}) *Request {
	req.SetHeaderSingle("Content-Type", "application/x-www-form-urlencoded")
	return req.SetBody(params)
}

// set request header
func (req *Request) SetHeader(header http.Header) *Request {
	req.Header = header
	return req
}

// get request header
func (req *Request) GetHeader() http.Header {
	return req.Header
}

// set request header single
func (req *Request) SetHeaderSingle(key, val string) *Request {
	req.Header.Set(key, val)
	return req
}

// get request header single
func (req *Request) GetHeaderSingle(key string) string {
	return req.Header.Get(key)
}

// merge request origin header and header
func (req *Request) SetHeaders(header http.Header) *Request {
	for key, val := range header {
		for _, v := range val {
			req.Header.Set(key, v)
		}
	}
	return req
}

// set request referer
func (req *Request) SetReferer(referer string) *Request {
	req.Header.Set("Referer", referer)
	return req
}

// set request charset
func (req *Request) SetCharset(charset string) *Request {
	req.SetHeaderSingle("Accept-Charset", charset)
	return req
}

// set request user-agent
func (req *Request) SetUserAgent(ua string) *Request {
	req.SetHeaderSingle("User-Agent", ua)
	return req
}

// get request user-agent
func (req *Request) GetUserAgent() string {
	return req.GetHeaderSingle("User-Agent")
}

// set the proxy for this request
// example: http://127.0.0.1:8080
func (req *Request) SetProxy(rawurl string) *Request {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	req.Proxy = u
	return req
}

// set cookies to request
// sample:
// 		quick.SetCookies(
//			quick.NewCookiesWithString("key1=value1; key2=value2; key3=value3")
//		)
func (req *Request) SetCookies(cookies Cookies) *Request {
	req.Cookies = cookies
	return req
}
