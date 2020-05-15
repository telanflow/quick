package quick

import (
	"net/http"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest().SetMethod(http.MethodGet).SetUrl("http://www.baidu.com")
	t.Log(req)
}
