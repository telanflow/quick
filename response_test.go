package quick

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBuildResponse(t *testing.T) {
	asserts := assert.New(t)

	ser := RunServer()
	defer ser.Close()

	httpResp, err := http.Get(ser.URL)
	asserts.Equal(err, nil)

	resp, err := BuildResponse(httpResp)
	asserts.Equal(err, nil)

	asserts.Equal(resp.StatusCode, 200)
	asserts.Equal(resp.Body.String(), "quick")
}
