package requests

import (
	"testing"
)

func TestNewCookiesWithString(t *testing.T) {
	rawstr := "_octo=GH1.1.1392344520.1564456933; _ga=GA1.2.1229933325.1564456966;"
	cookies := NewCookiesWithString(rawstr)
	if len(cookies) != 2 {
		t.Fatal("NewCookiesWithString fail")
	}
	t.Log(cookies)
}

func BenchmarkNewCookiesWithString(b *testing.B) {
	rawstr := "_octo=GH1.1.1392344520.1564456933; _ga=GA1.2.1229933325.1564456966;"
	for i := 0; i < b.N; i++ {
		NewCookiesWithString(rawstr)
	}
}
