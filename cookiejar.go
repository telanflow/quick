package quick

import (
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)

// create a cookieJar
func NewCookieJar() (http.CookieJar, error) {
	cookieJarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	return cookiejar.New(&cookieJarOptions)
}
