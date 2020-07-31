package quick

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HandlerFunc func(r *http.Request)

type Session struct {
	Header     http.Header
	Proxy      *url.URL
	Timeout    time.Duration
	transport  *http.Transport
	client     *http.Client
	middleware []HandlerFunc
	i          int
}

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
	}
}

// get request
func (session *Session) Get(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodGet).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// post request
func (session *Session) Post(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodPost).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// postForm request
func (session *Session) PostFormData(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodPost).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// put request
func (session *Session) Put(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodPut).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// head request
func (session *Session) Head(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodHead).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// patch request
func (session *Session) Patch(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodPatch).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// options request
func (session *Session) Options(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodOptions).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// delete request
func (session *Session) Delete(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodDelete).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// connect request
func (session *Session) Connect(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodConnect).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// trace request
func (session *Session) Trace(rawurl string, ops ...OptionFunc) (*Response, error) {
	req := NewRequest().SetMethod(http.MethodTrace).SetUrl(rawurl)
	return session.Suck(req, ops...)
}

// download file
func (session *Session) Download(rawurl string, toFile string) error {
	req := NewRequest().SetMethod(http.MethodGet).SetUrl(rawurl)
	resp, err := session.Suck(req)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(toFile, resp.GetBody(), 0644)
}

// ssl skip verify
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

// set session global header single
func (session *Session) SetHeaderSingle(key, val string) *Session {
	session.Header.Set(key, val)
	return session
}

// get session global header single
func (session *Session) GetHeaderSingle(key string) string {
	return session.Header.Get(key)
}

// set global header
func (session *Session) SetHeader(h http.Header) *Session {
	session.Header = h
	return session
}

// get global header
func (session *Session) GetHeader() http.Header {
	return session.Header
}

// set session global user-agent
func (session *Session) SetUserAgent(ua string) *Session {
	session.SetHeaderSingle("User-Agent", ua)
	return session
}

// get session global user-agent
func (session *Session) GetUserAgent() string {
	return session.GetHeaderSingle("User-Agent")
}

// get session global proxy url
func (session *Session) GetProxyUrl() string {
	if session.Proxy == nil {
		return ""
	}
	return session.Proxy.String()
}

// set session global proxy url
func (session *Session) SetProxyUrl(rawurl string) *Session {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	session.Proxy = u
	return session
}

// get session global proxy url
func (session *Session) GetProxyURL() *url.URL {
	return session.Proxy
}

// set session global proxy url
func (session *Session) SetProxyURL(u *url.URL) *Session {
	session.Proxy = u
	return session
}

// set session global request timeout
// example: time.Second * 30
func (session *Session) SetTimeout(t time.Duration) *Session {
	session.Timeout = t
	return session
}

// set session global proxy handler
// handler: func(req *http.Request) (*url.URL, error)
func (session *Session) SetProxyHandler(handler func(req *http.Request) (*url.URL, error)) *Session {
	session.transport.Proxy = handler
	return session
}

// set session global checkRedirect handler
// handler: func(req *http.Request, via []*http.Request) error
func (session *Session) SetCheckRedirectHandler(handler func(req *http.Request, via []*http.Request) error) *Session {
	session.client.CheckRedirect = handler
	return session
}

// set session global cookieJar
func (session *Session) SetCookieJar(jar http.CookieJar) *Session {
	session.client.Jar = jar
	return session
}

// returns the cookies of the given url in Session.
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

// set cookies of the url in Session.
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

// use middleware handler
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

// request suck data
func (session *Session) Suck(req *Request, ops ...OptionFunc) (*Response, error) {
	// Apply the HTTP request options
	for _, option := range ops {
		option(req)
	}

	return transmission(session, req)
}

// send http.Request
func (session *Session) Do(req *http.Request) (*Response, error) {
	// Set timeout to request context.
	// Default timeout is 30s.
	timeout := time.Second * 30
	if session.Timeout > 0 {
		timeout = session.Timeout
	}
	ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)

	if session.Proxy != nil {
		ctx = context.WithValue(ctx, ContextProxyKey, session.Proxy)
	}

	// set redirectNum to request context.
	ctx = context.WithValue(ctx, ContextRedirectNumKey, DefaultRedirectNum)

	req = req.WithContext(ctx)

	// merge request header and session header
	req.Header = MergeHeaders(session.Header, req.Header)

	// middleware
	session.next(req)

	// start request time
	startTime := time.Now()

	// http.Client send request
	httpResponse, err := session.client.Do(req)
	if err != nil {
		// check timeout error
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return nil, WrapErr(ErrTimeout, err.Error())
		}
		return nil, WrapErr(err, "Request Error")
	}
	defer func() {
		if err := httpResponse.Body.Close(); err != nil {
			panic(err)
		}
	}()

	resp, err := BuildResponse(httpResponse)
	if err != nil {
		return nil, WrapErr(err, "build Response Error")
	}

	// request exec time
	resp.ExecTime = time.Now().Sub(startTime)

	// cancel the timeout context after request successed.
	timeoutCancel()

	return resp, nil
}

// send the request
func transmission(session *Session, req *Request) (*Response, error) {
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

	if err != nil {
		// check timeout error
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return nil, WrapErr(ErrTimeout, err.Error())
		}
		return nil, WrapErr(err, "Request Error")
	}
	defer func() {
		if err := httpResponse.Body.Close(); err != nil {
			panic(err)
		}
	}()

	resp, err := BuildResponse(httpResponse)
	if err != nil {
		return nil, WrapErr(err, "build Response Error")
	}

	// request id
	resp.RequestId = req.Id
	// request exec time
	resp.ExecTime = time.Now().Sub(startTime)

	// cancel the timeout context after request successed.
	timeoutCancel()

	return resp, nil
}
