package requests

import (
	"encoding/xml"
	"errors"
	"html"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	RequestId     uint64 // request id
	Status        string // e.g. "200 OK"
	StatusCode    int    // e.g. 200
	Proto         string // e.g. "HTTP/1.0"
	ProtoMajor    int    // e.g. 1
	ProtoMinor    int    // e.g. 0
	Header        http.Header
	Body          []byte
	ContentLength int64
	ExecTime      time.Duration // request exec time
}

func BuildResponse(resp *http.Response) (*Response, error) {
	if resp == nil {
		return nil, errors.New("http response is nil")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		Proto:         resp.Proto,
		ProtoMajor:    resp.ProtoMajor,
		ProtoMinor:    resp.ProtoMinor,
		Header:        CopyHeader(resp.Header),
		Body:          body,
		ContentLength: resp.ContentLength,
	}, nil
}

func (r *Response) GetHeader() http.Header {
	return r.Header
}

func (r *Response) GetHeaderSingle(key string) string {
	return r.Header.Get(key)
}

func (r *Response) GetContextType() string {
	return r.GetHeaderSingle("Content-Type")
}

func (r *Response) GetJson(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

func (r *Response) GetHtml() string {
	return html.UnescapeString(string(r.Body))
}

func (r *Response) GetXml(v interface{}) error {
	return xml.Unmarshal(r.Body, v)
}

func (r *Response) String() string {
	return string(r.Body)
}

func (r *Response) GetBody() []byte {
	return r.Body
}
