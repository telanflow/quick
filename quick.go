package quick

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// default global session
var defaultSession = NewSession()

// json library
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// SetBaseURL method is to set Base URL in the client instance. It will be used with request
// raised from this client with relative URL
//		// Setting HTTP address
//		quick.SetBaseURL("http://myjeeva.com")
//
//		// Setting HTTPS address
//		quick.SetBaseURL("https://myjeeva.com")
//
// Since v0.4.0
func SetBaseURL(url string) *Session {
	defaultSession.BaseURL = strings.TrimRight(url, "/")
	return defaultSession
}

// SetLogger method sets given writer for logging Quick request and response details.
//
// Compliant to interface `quick.Logger`.
func SetLogger(l Logger) *Session {
	defaultSession.log = l
	return defaultSession
}

// Get request
func Get(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Get(rawurl, ops...)
}

// Post request
func Post(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Post(rawurl, ops...)
}

// PostFormData request
func PostFormData(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.PostFormData(rawurl, ops...)
}

// Put request
func Put(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Put(rawurl, ops...)
}

// Head request
func Head(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Head(rawurl, ops...)
}

// Patch request
func Patch(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Patch(rawurl, ops...)
}

// Options request
func Options(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Options(rawurl, ops...)
}

// Delete request
func Delete(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Delete(rawurl, ops...)
}

// Connect request
func Connect(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Connect(rawurl, ops...)
}

// Download download file to save hard disk
func Download(rawurl string, toFile string) error {
	return defaultSession.Download(rawurl, toFile)
}

// InsecureSkipVerify ssl skip verify
func InsecureSkipVerify(skip bool) *Session {
	return defaultSession.InsecureSkipVerify(skip)
}

// SetHeaderSingle set global header single
func SetHeaderSingle(key, val string) *Session {
	return defaultSession.SetHeaderSingle(key, val)
}

// GetHeaderSingle get global header single
func GetHeaderSingle(key string) string {
	return defaultSession.GetHeaderSingle(key)
}

// SetHeader set global header
func SetHeader(h http.Header) *Session {
	defaultSession.Header = h
	return defaultSession
}

// GetHeader get global header
func GetHeader() http.Header {
	return defaultSession.Header
}

// SetUserAgent set global user-agent
func SetUserAgent(ua string) *Session {
	return defaultSession.SetHeaderSingle("User-Agent", ua)
}

// GetUserAgent get global user-agent
func GetUserAgent() string {
	return defaultSession.GetHeaderSingle("User-Agent")
}

// GetProxyUrl get session global proxy url
func GetProxyUrl() string {
	return defaultSession.GetProxyUrl()
}

// SetProxyUrl set global proxy url
func SetProxyUrl(rawurl string) *Session {
	return defaultSession.SetProxyUrl(rawurl)
}

// GetProxyURL get session global proxy url
func GetProxyURL() *url.URL {
	return defaultSession.GetProxyURL()
}

// SetProxyURL set global proxy url
func SetProxyURL(u *url.URL) *Session {
	return defaultSession.SetProxyURL(u)
}

// SetTimeout set global request timeout
// example: time.Second * 30
func SetTimeout(t time.Duration) *Session {
	return defaultSession.SetTimeout(t)
}

// SetProxyHandler set global proxy handler
// handler: func(req *http.Request) (*url.URL, error)
func SetProxyHandler(handler func(req *http.Request) (*url.URL, error)) *Session {
	return defaultSession.SetProxyHandler(handler)
}

// SetCheckRedirectHandler set global checkRedirect handler
// handler: func(req *http.Request, via []*http.Request) error
func SetCheckRedirectHandler(handler func(req *http.Request, via []*http.Request) error) *Session {
	return defaultSession.SetCheckRedirectHandler(handler)
}

// SetCookieJar set global cookieJar
func SetCookieJar(jar http.CookieJar) *Session {
	return defaultSession.SetCookieJar(jar)
}

// Use use middleware handler
func Use(middleware ...HandlerFunc) *Session {
	defaultSession.Use(middleware...)
	return defaultSession
}

// EnableTrace method enables the Quick client trace for the requests fired from
// the client using `httptrace.ClientTrace` and provides insights.
//
// 		resp, err := quick.EnableTrace().Get("https://httpbin.org/get")
// 		fmt.Println("Error:", err)
// 		fmt.Println("Trace Info:", resp.TraceInfo())
//
// Also `Request.EnableTrace` available too to get trace info for single request.
//
// Since v0.4.0
func EnableTrace() *Session {
	defaultSession.trace = true
	return defaultSession
}

// DisableTrace method disables the Quick client trace. Refer to `quick.EnableTrace`.
//
// Since v0.4.0
func DisableTrace() *Session {
	defaultSession.trace = false
	return defaultSession
}

// Suck request suck data
func Suck(req *Request, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Suck(req, ops...)
}

// Do send http.Request
func Do(req *http.Request) (*Response, error) {
	return defaultSession.Do(req)
}
