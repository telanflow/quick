package quick

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HandlerFunc func(r *http.Request)

// Session is a http.Client
type Session struct {
	BaseURL    string
	Header     http.Header
	Proxy      *url.URL
	Timeout    time.Duration
	transport  *http.Transport
	client     *http.Client
	middleware []HandlerFunc
	i          int
	log        Logger
	trace      bool
}

// NewSession create a session
func NewSession(options ...*SessionOptions) *Session {
	var sessionOptions *SessionOptions
	if len(options) > 0 {
		sessionOptions = options[0]
	} else {
		sessionOptions = DefaultSessionOptions()
	}

	// set transport parameters.
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   sessionOptions.DialTimeout,
			KeepAlive: sessionOptions.DialKeepAlive,
		}).DialContext,
		MaxIdleConns:          sessionOptions.MaxIdleConns,
		MaxIdleConnsPerHost:   sessionOptions.MaxIdleConnsPerHost,
		MaxConnsPerHost:       sessionOptions.MaxConnsPerHost,
		IdleConnTimeout:       sessionOptions.IdleConnTimeout,
		TLSHandshakeTimeout:   sessionOptions.TLSHandshakeTimeout,
		ExpectContinueTimeout: sessionOptions.ExpectContinueTimeout,
		Proxy:                 proxyFunc,
	}
	if sessionOptions.DisableDialKeepAlives {
		transport.DisableKeepAlives = true
	}

	client := &http.Client{
		Transport:     transport,
		CheckRedirect: redirectFunc,
	}

	// set CookieJar
	if sessionOptions.DisableCookieJar == false {
		jar, err := NewCookieJar()
		if err != nil {
			return nil
		}
		client.Jar = jar
	}

	// Set default user agent
	return &Session{
		Header:     make(http.Header),
		client:     client,
		transport:  transport,
		middleware: make([]HandlerFunc, 0),
		i:          0,
		log:        createLogger(), // Logger
		trace:      false,
	}
}

// SetBaseURL method is to set Base URL in the client instance. It will be used with request
// raised from this client with relative URL
//		// Setting HTTP address
//		session.SetBaseURL("http://myjeeva.com")
//
//		// Setting HTTPS address
//		session.SetBaseURL("https://myjeeva.com")
//
// Since v0.4.0
func (session *Session) SetBaseURL(url string) *Session {
	session.BaseURL = strings.TrimRight(url, "/")
	return session
}

// SetLogger method sets given writer for logging Quick request and response details.
//
// Compliant to interface `quick.Logger`.
func (session *Session) SetLogger(l Logger) *Session {
	session.log = l
	return session
}

// Get request
func (session *Session) Get(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodGet).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// Post request
func (session *Session) Post(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodPost).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// PostFormData postForm request
func (session *Session) PostFormData(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodPost).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// Put request
func (session *Session) Put(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodPut).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// Head request
func (session *Session) Head(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodHead).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// Patch request
func (session *Session) Patch(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodPatch).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// Options request
func (session *Session) Options(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodOptions).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// Delete request
func (session *Session) Delete(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodDelete).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// Connect request
func (session *Session) Connect(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodConnect).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// Download file
func (session *Session) Download(rawurl string, toFile string) error {
	req := NewRequest().SetMethod(http.MethodGet).SetUrl(rawurl)
	resp, err := session.Suck(req)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(toFile, resp.GetBody(), 0644)
}

// InsecureSkipVerify ssl skip verify
func (session *Session) InsecureSkipVerify(skip bool) *Session {
	if session.transport.TLSClientConfig != nil {
		session.transport.TLSClientConfig.InsecureSkipVerify = skip
	} else {
		session.transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: skip,
		}
	}
	return session
}

// SetHeaderSingle set session global header single
func (session *Session) SetHeaderSingle(key, val string) *Session {
	session.Header.Set(key, val)
	return session
}

// GetHeaderSingle get session global header single
func (session *Session) GetHeaderSingle(key string) string {
	return session.Header.Get(key)
}

// SetHeader set global header
func (session *Session) SetHeader(h http.Header) *Session {
	session.Header = h
	return session
}

// GetHeader get global header
func (session *Session) GetHeader() http.Header {
	return session.Header
}

// SetUserAgent set session global user-agent
func (session *Session) SetUserAgent(ua string) *Session {
	session.SetHeaderSingle("User-Agent", ua)
	return session
}

// GetUserAgent get session global user-agent
func (session *Session) GetUserAgent() string {
	return session.GetHeaderSingle("User-Agent")
}

// GetProxyUrl get session global proxy url
func (session *Session) GetProxyUrl() string {
	if session.Proxy == nil {
		return ""
	}
	return session.Proxy.String()
}

// SetProxyUrl set session global proxy url
func (session *Session) SetProxyUrl(rawurl string) *Session {
	u, err := url.Parse(rawurl)
	if err != nil {
		session.log.Errorf("url parse fail: %s", err)
	}
	session.Proxy = u
	return session
}

// GetProxyURL get session global proxy url
func (session *Session) GetProxyURL() *url.URL {
	return session.Proxy
}

// SetProxyURL set session global proxy url
func (session *Session) SetProxyURL(u *url.URL) *Session {
	session.Proxy = u
	return session
}

// SetTimeout set session global request timeout
// example: time.Second * 30
func (session *Session) SetTimeout(t time.Duration) *Session {
	session.Timeout = t
	return session
}

// SetProxyHandler set session global proxy handler.
// handler: func(req *http.Request) (*url.URL, error)
func (session *Session) SetProxyHandler(handler func(req *http.Request) (*url.URL, error)) *Session {
	session.transport.Proxy = handler
	return session
}

// SetCheckRedirectHandler set session global checkRedirect handler.
// handler: func(req *http.Request, via []*http.Request) error
func (session *Session) SetCheckRedirectHandler(handler func(req *http.Request, via []*http.Request) error) *Session {
	session.client.CheckRedirect = handler
	return session
}

// SetCookieJar set session global cookieJar.
func (session *Session) SetCookieJar(jar http.CookieJar) *Session {
	session.client.Jar = jar
	return session
}

// Cookies returns the cookies of the given url in Session.
func (session *Session) Cookies(rawurl string) Cookies {
	if session.client.Jar == nil {
		return nil
	}
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return nil
	}
	return session.client.Jar.Cookies(parsedURL)
}

// SetCookies set cookies of the url in Session.
func (session *Session) SetCookies(rawurl string, cookies Cookies) {
	if session.client.Jar == nil {
		return
	}
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return
	}
	session.client.Jar.SetCookies(parsedURL, cookies)
}

// Use use middleware handler.
func (session *Session) Use(middleware ...HandlerFunc) *Session {
	session.middleware = append(session.middleware, middleware...)
	return session
}

// next middleware
func (session *Session) next(r *http.Request) {
	current := session.i
	n := len(session.middleware)
	session.i++
	if current >= n {
		return
	}

	task := session.middleware[current]
	if task == nil {
		return
	}

	// handler
	task(r)

	session.next(r)
}

// EnableTrace method enables the Quick client trace for the requests fired from
// the client using `httptrace.ClientTrace` and provides insights.
//
// 		session := quick.NewSession().EnableTrace()
//
// 		resp, err := session.Get("https://httpbin.org/get")
// 		fmt.Println("Error:", err)
// 		fmt.Println("Trace Info:", resp.TraceInfo())
//
// Also `Request.EnableTrace` available too to get trace info for single request.
//
// Since v0.4.0
func (session *Session) EnableTrace() *Session {
	session.trace = true
	return session
}

// DisableTrace method disables the Quick client trace. Refer to `Session.EnableTrace`.
//
// Since v0.4.0
func (session *Session) DisableTrace() *Session {
	session.trace = false
	return session
}

// Suck request suck data
func (session *Session) Suck(req *Request, ops ...OptionFunc) (*Response, error) {
	// Apply the HTTP request options
	for _, option := range ops {
		option(req)
	}

	var (
		ctx           context.Context
		timeoutCancel context.CancelFunc
	)

	// Set timeout to request context.
	// Default timeout is 30s.
	timeout := time.Second * 30
	if req.Timeout > 0 {
		timeout = req.Timeout
	} else if session.Timeout > 0 {
		timeout = session.Timeout
	}

	if req.ctx == nil {
		ctx, timeoutCancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, timeoutCancel = context.WithTimeout(req.ctx, timeout)
	}

	// set proxy to request context.
	if req.Proxy != nil {
		ctx = context.WithValue(ctx, ContextProxyKey, req.Proxy)
	} else if session.Proxy != nil {
		ctx = context.WithValue(ctx, ContextProxyKey, session.Proxy)
	}

	// set redirectNum to request context.
	ctx = context.WithValue(ctx, ContextRedirectNumKey, req.RedirectNum)

	// Enable trace
	if session.trace || req.trace {
		req.clientTrace = &clientTrace{}
		ctx = req.clientTrace.createContext(ctx)
	}

	// request set base url
	if session.BaseURL != "" {
		rawurl := fmt.Sprintf("%s%s", session.BaseURL, req.URL)
		if u, err := url.Parse(rawurl); err == nil {
			req.URL = u
		}
	}

	httpRequest, err := http.NewRequestWithContext(ctx, req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, err
	}

	// customize the request Host field
	if len(req.host) > 0 {
		httpRequest.Host = req.host
	}

	// handle cookies
	if req.Cookies != nil {
		// if cookieJar is enabled, the requested cookies are merged
		if session.client.Jar != nil {
			session.client.Jar.SetCookies(httpRequest.URL, req.Cookies)
		} else {
			for _, cookie := range req.Cookies {
				httpRequest.AddCookie(cookie)
			}
		}
	}

	// merge request header and session header
	httpRequest.Header = MergeHeaders(session.Header, req.Header)

	// middleware
	session.next(httpRequest)

	// start request time
	startTime := time.Now()

	// do request
	httpResponse, err := session.client.Do(httpRequest)
	defer func() {
		if httpResponse != nil && httpResponse.Body != nil {
			_ = httpResponse.Body.Close()
		}
	}()

	if err != nil {
		// check timeout error
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return nil, WrapErr(ErrTimeout, err.Error())
		}
		return nil, WrapErr(err, "Request Error")
	}

	resp, err := BuildResponse(httpResponse)
	if err != nil {
		return nil, WrapErr(err, "build Response Error")
	}

	// request
	resp.RequestId = req.Id
	// request exec time
	resp.ExecTime = time.Now().Sub(startTime)
	// trace info
	resp.clientTrace = req.clientTrace

	// cancel the timeout context after request successful.
	timeoutCancel()

	return resp, nil
}

// Do send http.Request
func (session *Session) Do(req *http.Request) (*Response, error) {
	// Set timeout to request context.
	// Default timeout is 30s.
	timeout := time.Second * 30
	if session.Timeout > 0 {
		timeout = session.Timeout
	}

	// request set base url
	if session.BaseURL != "" {
		rawurl := fmt.Sprintf("%s%s", session.BaseURL, req.URL)
		if u, err := url.Parse(rawurl); err == nil {
			req.URL = u
			req.Host = u.Host
		}
	}

	ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)

	if session.Proxy != nil {
		ctx = context.WithValue(ctx, ContextProxyKey, session.Proxy)
	}

	// set redirectNum to request context.
	ctx = context.WithValue(ctx, ContextRedirectNumKey, DefaultRedirectNum)

	// Enable trace
	var ct *clientTrace
	if session.trace {
		ct = &clientTrace{}
		ctx = ct.createContext(ctx)
	}

	// merge request header and session header
	req.Header = MergeHeaders(session.Header, req.Header)

	req = req.WithContext(ctx)

	// middleware
	session.next(req)

	// start request time
	startTime := time.Now()

	// http.Client send request
	httpResponse, err := session.client.Do(req)
	defer func() {
		if httpResponse == nil {
			return
		}
		if httpResponse.Body == nil {
			return
		}
		if err := httpResponse.Body.Close(); err != nil {
			session.log.Warnf("response close body fail: %s", err)
		}
	}()

	if err != nil {
		// check timeout error
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return nil, WrapErr(ErrTimeout, err.Error())
		}
		return nil, WrapErr(err, "Request Error")
	}

	resp, err := BuildResponse(httpResponse)
	if err != nil {
		return nil, WrapErr(err, "build Response Error")
	}

	// request
	resp.clientTrace = ct
	// request exec time
	resp.ExecTime = time.Now().Sub(startTime)

	// cancel the timeout context after request successed.
	timeoutCancel()

	return resp, nil
}
