package requests

import (
	"net/http"
	"testing"
)

const BaiDuUrl = "http://www.baidu.com"

func TestNewSession(t *testing.T) {
	session := NewSession()
	//session.SetProxy("http://127.0.0.1:8080")
	resp, err := session.Get(BaiDuUrl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func BenchmarkSession_Get(b *testing.B) {
	session := NewSession()
	for i := 0; i < b.N; i++ {
		_, _ = session.Get(BaiDuUrl)
	}
}

func BenchmarkSession_Get_Parallel(b *testing.B) {
	session := NewSession()
	//session.SetProxy("http://127.0.0.1:8080")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = session.Get(BaiDuUrl)
		}
	})
}

func TestSession_Post(t *testing.T) {
	cookieJar, err := NewCookieJar()
	if err != nil {
		t.Fatal(err)
	}

	session := NewSession()
	session.SetCookieJar(cookieJar)
	//session.SetProxy("http://127.0.0.1:8080")

	req, err := NewRequest(http.MethodGet, BaiDuUrl)
	if err != nil {
		t.Fatal(err)
	}

	req.SetCookies(NewCookiesWithString("test=111111"))
	resp, err := session.Suck(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}
