package quick

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestMergeHeaders(t *testing.T) {
	h1 := make(http.Header)
	h1.Set("Context-Type", "text/html")

	h2 := make(http.Header)
	h2.Set("Context-Type", "text/xml")
	h2.Set("Server", "nginx")

	h3 := MergeHeaders(h1, h2)

	asserts := assert.New(t)
	asserts.Equal(h3.Get("Context-type"), "text/xml")
	asserts.Equal(h3.Get("Server"), "nginx")
}

func TestCopyHeader(t *testing.T) {
	h1 := make(http.Header)
	h1.Set("Context-Type", "text/html")

	h2 := CopyHeader(h1)

	asserts := assert.New(t)
	asserts.Equal(h1.Get("Context-Type"), h2.Get("Context-Type"))
}

func TestMergeQueryString(t *testing.T) {
	asserts := assert.New(t)

	u, _ := url.Parse("http://example.com")
	u2, _ := MergeQueryString(u, "c=3&d=4")
	asserts.Equal(u2.Query().Encode(), "c=3&d=4")

	u, _ = url.Parse("http://example.com?a=1&b=2")
	u2, _ = MergeQueryString(u, "c=3&d=4")
	asserts.Equal(u2.Query().Encode(), "a=1&b=2&c=3&d=4")

	u, _ = url.Parse("http://example.com?a=1")
	u2, _ = MergeQueryString(u, "a=2")
	asserts.Equal(u2.Query().Encode(), "a=2")

	u, _ = url.Parse("http://example.com?a=1")
	u2, _ = MergeQueryString(u, "a=2&a=3")
	asserts.Equal(u2.Query().Encode(), "a=3")
}

func TestReplaceQueryString(t *testing.T) {
	asserts := assert.New(t)

	u, _ := url.Parse("http://example.com")
	u2, _ := ReplaceQueryString(u, "c=3&d=4")
	asserts.Equal(u2.Query().Encode(), "c=3&d=4")

	u, _ = url.Parse("http://example.com?a=1")
	u2, _ = ReplaceQueryString(u, "c=3&d=4")
	asserts.Equal(u2.Query().Encode(), "c=3&d=4")
}
