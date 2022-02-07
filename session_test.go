package quick

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RunServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("quick", "hd")
		_, _ = w.Write([]byte("quick"))
	}))
}

func TestSession_Get(t *testing.T) {
	ser := RunServer()
	defer ser.Close()

	resp, err := NewSession().Get(ser.URL)
	if err != nil {
		t.Fatal(err)
	}

	asserts := assert.New(t)
	asserts.Equal(resp.StatusCode, 200)
	asserts.Equal(resp.Body.String(), "quick")
	asserts.Equal(resp.Header.Get("quick"), "hd")
}

func BenchmarkSession_Get(b *testing.B) {
	ser := RunServer()
	defer ser.Close()

	session := NewSession()

	for i := 0; i < b.N; i++ {
		_, _ = session.Get(ser.URL)
	}
}

func BenchmarkSession_Get_Parallel(b *testing.B) {
	ser := RunServer()
	defer ser.Close()
	session := NewSession()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = session.Get(ser.URL)
		}
	})
}

func TestSession_Post(t *testing.T) {
	asserts := assert.New(t)

	ser := RunServer()
	defer ser.Close()

	cookieJar, _ := NewCookieJar()
	session := NewSession()
	session.SetCookieJar(cookieJar)

	req := NewRequest().SetMethod(http.MethodGet).SetUrl(ser.URL)
	req.SetCookies(NewCookiesWithString("test=111111"))

	resp, err := session.Suck(req)
	if err != nil {
		t.Fatal(err)
	}

	asserts.Equal(resp.StatusCode, 200)
}

func TestSession_SetBaseURL(t *testing.T) {
	session := NewSession().SetBaseURL("https://httpbin.org/")

	resp, err := session.Get("/get?a=1#12")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.Status)

	req, _ := http.NewRequest("GET", "/get?a=1#12", nil)
	resp, err = session.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestSession_EnableTrace(t *testing.T) {
	session := NewSession()
	resp, _ := session.EnableTrace().Get("https://httpbin.org/get")
	fmt.Println(resp.TraceInfo())
}
