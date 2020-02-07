package requests

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type Session struct {
	Header     http.Header
	Proxy      *url.URL
	Timeout    time.Duration
	transport  *http.Transport
	client     *http.Client
	middleware []func()
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
		cookieJarOptions := cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}
		jar, err := cookiejar.New(&cookieJarOptions)
		if err != nil {
			return nil
		}
		client.Jar = jar
	}

	// Set default user agent
	header := make(http.Header)
	header.Add("User-Agent", UA)

	return &Session{
		Header:    header,
		client:    client,
		transport: transport,
	}
}

// get request
func (session *Session) Get(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodGet, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetQueryString(getParam(params...))
	}
	return session.Suck(req)
}

// post request
func (session *Session) Post(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodPost, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBody(getParam(params...))
	}
	return session.Suck(req)
}

// postForm request
func (session *Session) PostFormData(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodPost, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBodyFormData(getParam(params...))
	}
	return session.Suck(req)
}

// put request
func (session *Session) Put(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodPut, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBody(getParam(params...))
	}
	return session.Suck(req)
}

// head request
func (session *Session) Head(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodHead, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBody(getParam(params...))
	}
	return session.Suck(req)
}

// patch request
func (session *Session) Patch(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodPatch, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBody(getParam(params...))
	}
	return session.Suck(req)
}

// options request
func (session *Session) Options(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodOptions, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBody(getParam(params...))
	}
	return session.Suck(req)
}

// delete request
func (session *Session) Delete(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodDelete, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBody(getParam(params...))
	}
	return session.Suck(req)
}

// connect request
func (session *Session) Connect(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodConnect, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBody(getParam(params...))
	}
	return session.Suck(req)
}

// trace request
func (session *Session) Trace(rawurl string, params ...interface{}) (*Response, error) {
	req, err := NewRequest(http.MethodTrace, rawurl)
	if err != nil {
		return nil, err
	}
	if params != nil {
		req.SetBody(getParam(params...))
	}
	return session.Suck(req)
}

// download file
func (session *Session) Download(rawurl string, toFile string) error {
	req, err := NewRequest(http.MethodGet, rawurl)
	if err != nil {
		return err
	}

	resp, err := session.Suck(req)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(toFile, resp.Body, 0644)
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

// set session global user-agent
func (session *Session) SetUserAgent(ua string) *Session {
	session.SetHeaderSingle("User-Agent", ua)
	return session
}

// get session global user-agent
func (session *Session) GetUserAgent() string {
	return session.GetHeaderSingle("User-Agent")
}

// set session global proxy url
func (session *Session) SetProxy(rawurl string) *Session {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
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

// request suck data
func (session *Session) Suck(req *Request) (*Response, error) {
	return transmission(session, req)
}

// send the request
func transmission(session *Session, req *Request) (*Response, error) {
	// Set timeout to request context.
	// Default timeout is 30s.
	timeout := time.Second * 30
	if req.Timeout > 0 {
		timeout = req.Timeout
	} else if session.Timeout > 0 {
		timeout = session.Timeout
	}
	ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)

	// set proxy to request context.
	if req.Proxy != nil {
		ctx = context.WithValue(ctx, ContextProxyKey, req.Proxy)
	} else if session.Proxy != nil {
		ctx = context.WithValue(ctx, ContextProxyKey, session.Proxy)
	}

	// set redirectNum to request context.
	ctx = context.WithValue(ctx, ContextRedirectNumKey, req.RedirectNum)

	httpRequest, err := http.NewRequestWithContext(ctx, req.Method, req.Url, req.Body)
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
	httpRequest.Header = MergeHeaders(session.Header, httpRequest.Header)

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
