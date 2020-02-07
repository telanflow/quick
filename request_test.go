package requests

import (
	"net/http"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://www.baidu.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(req)
}
