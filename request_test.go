package quick

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRequest_Copy(t *testing.T) {
	asserts := assert.New(t)

	req1 := NewRequest()
	req1.SetUrl("http://example.com:8080")
	req1.SetMethod(http.MethodGet)
	req1.SetHeaderSingle("User-Agent", "quick")
	req1.SetProxyUrl("http://127.0.0.1:8080")
	req1.SetCookies(NewCookiesWithString("sessionid=2222222"))

	req2 := req1.Copy()

	p1 := fmt.Sprintf("%p", req1)
	p2 := fmt.Sprintf("%p", req2)
	asserts.NotEqual(p1, p2)

	p1 = fmt.Sprintf("%p", &req1.Id)
	p2 = fmt.Sprintf("%p", &req2.Id)
	asserts.NotEqual(p1, p2)
	asserts.NotEqual(req1.Id, req2.Id)

	p1 = fmt.Sprintf("%p", req1.URL)
	p2 = fmt.Sprintf("%p", req2.URL)
	asserts.NotEqual(p1, p2)
	asserts.Equal(req1.URL, req2.URL)

	p1 = fmt.Sprintf("%p", &req1.Method)
	p2 = fmt.Sprintf("%p", &req2.Method)
	asserts.NotEqual(p1, p2)
	asserts.Equal(req1.Method, req2.Method)

	p1 = fmt.Sprintf("%p", &req1.Header)
	p2 = fmt.Sprintf("%p", &req2.Header)
	asserts.NotEqual(p1, p2)
	asserts.Equal(req1.Header, req2.Header)

	p1 = fmt.Sprintf("%p", &req1.RedirectNum)
	p2 = fmt.Sprintf("%p", &req2.RedirectNum)
	asserts.NotEqual(p1, p2)
	asserts.Equal(req1.RedirectNum, req2.RedirectNum)

	p1 = fmt.Sprintf("%p", &req1.Timeout)
	p2 = fmt.Sprintf("%p", &req2.Timeout)
	asserts.NotEqual(p1, p2)
	asserts.Equal(req1.Timeout, req2.Timeout)

	p1 = fmt.Sprintf("%p", req1.Proxy)
	p2 = fmt.Sprintf("%p", req2.Proxy)
	asserts.NotEqual(p1, p2)
	asserts.Equal(req1.Proxy, req2.Proxy)

	p1 = fmt.Sprintf("%p", &req1.Cookies)
	p2 = fmt.Sprintf("%p", &req2.Cookies)
	p3 := fmt.Sprintf("%p", req1.Cookies)
	p4 := fmt.Sprintf("%p", req2.Cookies)
	asserts.NotEqual(p1, p2)
	asserts.NotEqual(p3, p4)
	asserts.Equal(req1.Cookies, req2.Cookies)
}
