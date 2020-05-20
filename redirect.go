package quick

import (
	"net/http"
)

const (
	// default redirect num
	DefaultRedirectNum = 10
	// request context redirect num key name
	ContextRedirectNumKey = "redirectNum"
)

// redirectFunc get redirectNum from request context and check redirect number.
func redirectFunc(req *http.Request, via []*http.Request) error {
	redirectNum := req.Context().Value(ContextRedirectNumKey).(int)
	if len(via) > redirectNum {
		err := &RedirectError{redirectNum}
		return WrapErr(err, "RedirectError")
	}
	return nil
}
