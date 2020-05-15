package quick

import (
	"errors"
	"golang.org/x/net/http/httpproxy"
	"net/http"
	"net/url"
	"sync"
)

// request context proxy key name
const ContextProxyKey = "proxy"

var (
	// proxyConfigOnce guards proxyConfig
	envProxyOnce      sync.Once
	envProxyFuncValue func(*url.URL) (*url.URL, error)
)

// proxyFunc get proxy from request context.
// If there is no proxy set, use default proxy from environment.
func proxyFunc(req *http.Request) (*url.URL, error) {
	proxyURL := req.Context().Value(ContextProxyKey) // get proxy *url.URL form context

	// If there is no proxy set, use default proxy from environment.
	// This mitigates expensive lookups on some platforms (e.g. Windows).
	envProxyOnce.Do(func() {
		envProxyFuncValue = httpproxy.FromEnvironment().ProxyFunc()
	})

	if proxyURL != nil {
		httpURL, ok := proxyURL.(*url.URL)
		if !ok {
			return nil, WrapErr(errors.New("proxy address is not of type *url.URL"), "HTTP Proxy error, please check proxy url")
		}
		return httpURL, nil
	}

	return envProxyFuncValue(req.URL)
}
