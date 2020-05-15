package quick

import (
	"net/http"
	"testing"
)

func TestMergeHeaders(t *testing.T) {
	h1 := make(http.Header)
	h1.Set("Context-Type", "text/html")
	h1.Set("User-Agent", UA)

	h2 := make(http.Header)
	h2.Set("Context-Type", "text/xml")
	h2.Set("Server", "nginx")

	h3 := MergeHeaders(h1, h2)
	if h3.Get("User-Agent") != UA {
		t.Fail()
	}
	if h3.Get("Context-type") != "text/xml" {
		t.Fail()
	}
	if h3.Get("Server") != "nginx" {
		t.Fail()
	}

	t.Log(h3)
}
