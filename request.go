package quick

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
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
	URL         *url.URL
	Method      string
	Header      http.Header   // request headers
	Body        io.Reader     // request encode
	RedirectNum int           // Number of redirects requested. default 5
	Timeout     time.Duration // request timeout
	Proxy       *url.URL      // request proxy url
	Cookies     Cookies       // request cookies

	host string // customize the request Host field
	ctx  context.Context
}

// create a request instance
func NewRequest() *Request {
	return NewRequestWithContext(nil)
}

// create a request instance with context.Context
func NewRequestWithContext(ctx context.Context) *Request {
	return &Request{
		Id:          atomic.AddUint64(&sequenceNo, 1),
		URL:         nil,
		Method:      http.MethodGet,
		Header:      make(http.Header),
		Body:        nil,
		RedirectNum: DefaultRedirectNum, // set request redirect num. default 10.
		Timeout:     30 * time.Second,
		Proxy:       nil,
		Cookies:     nil,
		ctx: 		 ctx,
	}
}

// Convert http.Request To Request
func ConvertHttpRequest(r *http.Request) *Request {
	// copy the URL
	newURL, _ := CopyURL(r.URL)

	// copy Cookies
	var copyBody io.Reader
	if r.Body != nil {
		copyBody := new(bytes.Buffer)
		_, _ = io.Copy(copyBody, r.Body)
	}

	// Generate a new request.Id
	newReq := NewRequest()
	newReq.URL 			= newURL
	newReq.Method 		= r.Method
	newReq.Header 		= CopyHeader(r.Header)
	newReq.Body 		= copyBody
	newReq.RedirectNum 	= DefaultRedirectNum
	newReq.Timeout 		= 30 * time.Second
	newReq.ctx 			= r.Context()
	return newReq
}

// with context.Context for Request
func (req *Request) WithContext(ctx context.Context) *Request {
	req.ctx = ctx
	return req
}

// get context.Context for Request
func (req *Request) Context() context.Context {
	return req.ctx
}

// set request url
func (req *Request) SetUrl(rawurl string) *Request {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	req.URL = u
	return req
}

// get request url
func (req *Request) GetUrl() string {
	return req.URL.String()
}

// set request url
func (req *Request) SetURL(u *url.URL) *Request {
	req.URL = u
	return req
}

// get request url
func (req *Request) GetURL() *url.URL {
	return req.URL
}

func (req *Request) SetMethod(method string) *Request {
	req.Method = strings.ToUpper(method)
	return req
}

func (req *Request) GetMethod() string {
	return req.Method
}

func (req *Request) SetTimeout(t time.Duration) *Request {
	req.Timeout = t
	return req
}

func (req *Request) GetTimeout() time.Duration {
	return req.Timeout
}

// Custom request host field
// GET /index HTTP/1.1
// Host: domain
// ....
func (req *Request) SetHost(host string) *Request {
	req.host = host
	return req
}

// BasicAuth returns the username and password provided in the request's
// Authorization header, if the request uses HTTP Basic Authentication.
// See RFC 2617, Section 2.
func (r *Request) BasicAuth() (username, password string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}
	return parseBasicAuth(auth)
}

// SetBasicAuth sets the request's Authorization header to use HTTP
// Basic Authentication with the provided username and password.
//
// With HTTP Basic Authentication the provided username and password
// are not encrypted.
//
// Some protocols may impose additional requirements on pre-escaping the
// username and password. For instance, when used with OAuth2, both arguments
// must be URL encoded first with url.QueryEscape.
func (r *Request) SetBasicAuth(username, password string) {
	r.Header.Set("Authorization", "Basic "+basicAuth(username, password))
}

// set GET parameters to request
func (req *Request) SetQueryString(params interface{}) *Request {
	buff := new(bytes.Buffer)
	form := new(encode.XWwwFormUrlencoded)
	form.SetValue(params)
	if err := form.Encode(buff); err != nil {
		panic(err)
	}

	// format get request parameters
	u, err := MergeQueryString(req.URL, buff.String())
	if err != nil {
		panic(err)
	}

	req.URL = u
	req.Body = nil
	return req
}

// set POST body to request
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

// set POST body (FormData) to request
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

// get proxy url for this request
func (req *Request) GetProxyUrl() string {
	if req.Proxy == nil {
		return ""
	}
	return req.Proxy.String()
}

// set the proxy for this request
// eg. "http://127.0.0.1:8080" "http://username:password@127.0.0.1:8080"
func (req *Request) SetProxyUrl(rawurl string) *Request {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	req.Proxy = u
	return req
}

// get proxy url.URL for this request
func (req *Request) GetProxyURL() *url.URL {
	return req.Proxy
}

// set the proxy url for this request
func (req *Request) SetProxyURL(u *url.URL) *Request {
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

// Copy request
func (req *Request) Copy() *Request {
	// copy the URL
	newURL, _ := CopyURL(req.URL)

	var copyBody io.Reader
	if req.Body != nil {
		copyBody := new(bytes.Buffer)
		_, _ = io.Copy(copyBody, req.Body)
	}

	// copy the proxy url
	copyProxy, _ := CopyURL(req.Proxy)

	var copyCookies Cookies
	if req.Cookies != nil {
		copyCookies = make(Cookies, len(req.Cookies))
		copy(copyCookies, req.Cookies)
	}

	// Generate a new request.Id
	newReq := NewRequest()
	newReq.URL 			= newURL
	newReq.Method 		= req.Method
	newReq.Header 		= CopyHeader(req.Header)
	newReq.Body 		= copyBody
	newReq.RedirectNum 	= req.RedirectNum
	newReq.Timeout 		= req.Timeout
	newReq.Proxy 		= copyProxy
	newReq.Cookies 		= copyCookies
	newReq.host 		= req.host
	newReq.ctx 			= req.ctx
	return newReq
}

// copy url.URL
func CopyURL(u *url.URL) (URL *url.URL, err error) {
	if u == nil {
		err = errors.New("copy url.URL is nil")
		return
	}

	// copy basic authentication username,password
	var user *url.Userinfo
	if u.User != nil {
		password, _ := u.User.Password()
		user = url.UserPassword(u.User.Username(), password)
	}

	URL = &url.URL{
		Scheme:     u.Scheme,
		Opaque:     u.Opaque,
		User:       user,
		Host:       u.Host,
		Path:       u.Path,
		RawPath:    u.RawPath,
		ForceQuery: u.ForceQuery,
		RawQuery:   u.RawQuery,
		Fragment:   u.Fragment,
	}
	return
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// See 2 (end of page 4) https://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
