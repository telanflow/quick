package quick

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"time"
)

// default global session
var defaultSession = NewSession()

// json library
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// get request
func Get(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Get(rawurl, ops...)
}

// post request
func Post(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Post(rawurl, ops...)
}

// postForm request
func PostFormData(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.PostFormData(rawurl, ops...)
}

// put request
func Put(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Put(rawurl, ops...)
}

// head request
func Head(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Head(rawurl, ops...)
}

// patch request
func Patch(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Patch(rawurl, ops...)
}

// options request
func Options(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Options(rawurl, ops...)
}

// delete request
func Delete(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Delete(rawurl, ops...)
}

// connect request
func Connect(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Connect(rawurl, ops...)
}

// trace request
func Trace(rawurl string, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Trace(rawurl, ops...)
}

// download file to save hard disk
func Download(rawurl string, toFile string) error {
	return defaultSession.Download(rawurl, toFile)
}

// ssl skip verify
func InsecureSkipVerify(skip bool) *Session {
	return defaultSession.InsecureSkipVerify(skip)
}

// set global header single
func SetHeaderSingle(key, val string) *Session {
	return defaultSession.SetHeaderSingle(key, val)
}

// get global header single
func GetHeaderSingle(key string) string {
	return defaultSession.GetHeaderSingle(key)
}

// set global header
func SetHeader(h http.Header) *Session {
	defaultSession.Header = h
	return defaultSession
}

// get global header
func GetHeader() http.Header {
	return defaultSession.Header
}

// set global user-agent
func SetUserAgent(ua string) *Session {
	return defaultSession.SetHeaderSingle("User-Agent", ua)
}

// get global user-agent
func GetUserAgent() string {
	return defaultSession.GetHeaderSingle("User-Agent")
}

// get session global proxy url
func GetProxyUrl() string {
	return defaultSession.GetProxyUrl()
}

// set global proxy url
func SetProxyUrl(rawurl string) *Session {
	return defaultSession.SetProxyUrl(rawurl)
}

// get session global proxy url
func GetProxyURL() *url.URL {
	return defaultSession.GetProxyURL()
}

// set global proxy url
func SetProxyURL(u *url.URL) *Session {
	return defaultSession.SetProxyURL(u)
}

// set global request timeout
// example: time.Second * 30
func SetTimeout(t time.Duration) *Session {
	return defaultSession.SetTimeout(t)
}

// set global proxy handler
// handler: func(req *http.Request) (*url.URL, error)
func SetProxyHandler(handler func(req *http.Request) (*url.URL, error)) *Session {
	return defaultSession.SetProxyHandler(handler)
}

// set global checkRedirect handler
// handler: func(req *http.Request, via []*http.Request) error
func SetCheckRedirectHandler(handler func(req *http.Request, via []*http.Request) error) *Session {
	return defaultSession.SetCheckRedirectHandler(handler)
}

// set global cookieJar
func SetCookieJar(jar http.CookieJar) *Session {
	return defaultSession.SetCookieJar(jar)
}

// use middleware handler
func Use(middleware ...HandlerFunc) *Session {
	defaultSession.Use(middleware...)
	return defaultSession
}

// request suck data
func Suck(req *Request, ops ...OptionFunc) (*Response, error) {
	return defaultSession.Suck(req, ops...)
}

// send http.Request
func Do(req *http.Request) (*Response, error) {
	return defaultSession.Do(req)
}
