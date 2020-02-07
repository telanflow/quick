package requests

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"time"
)

const UA = "Go library telanflow/requests"

// default global session
var defaultSession = NewSession()

// json library
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// get request
func Get(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Get(rawurl, params...)
}

// post request
func Post(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Post(rawurl, params...)
}

// postForm request
func PostFormData(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.PostFormData(rawurl, params)
}

// put request
func Put(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Put(rawurl, params...)
}

// head request
func Head(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Head(rawurl, params...)
}

// patch request
func Patch(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Patch(rawurl, params...)
}

// options request
func Options(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Options(rawurl, params...)
}

// delete request
func Delete(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Delete(rawurl, params...)
}

// connect request
func Connect(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Connect(rawurl, params...)
}

// trace request
func Trace(rawurl string, params ...interface{}) (*Response, error) {
	return defaultSession.Trace(rawurl, params...)
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

// set global user-agent
func SetUserAgent(ua string) *Session {
	return defaultSession.SetHeaderSingle("User-Agent", ua)
}

// get global user-agent
func GetUserAgent() string {
	return defaultSession.GetHeaderSingle("User-Agent")
}

// set global proxy url
func SetProxy(rawurl string) *Session {
	return defaultSession.SetProxy(rawurl)
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

// request suck data
func Suck(req *Request) (*Response, error) {
	return defaultSession.Suck(req)
}
