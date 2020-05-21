package quick

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCookiesWithString(t *testing.T) {
	rawstr := "_octo=GH1.1.1392344520.1564456933; _ga=GA1.2.1229933325.1564456966;"
	cookies := NewCookiesWithString(rawstr)

	asserts := assert.New(t)
	asserts.Equal(len(cookies), 2)

	asserts.Equal(cookies[0].Name, "_octo")
	asserts.Equal(cookies[0].Value, "GH1.1.1392344520.1564456933")

	asserts.Equal(cookies[1].Name, "_ga")
	asserts.Equal(cookies[1].Value, "GA1.2.1229933325.1564456966")
}

func BenchmarkNewCookiesWithString(b *testing.B) {
	rawstr := "_octo=GH1.1.1392344520.1564456933; _ga=GA1.2.1229933325.1564456966;"
	for i := 0; i < b.N; i++ {
		NewCookiesWithString(rawstr)
	}
}
