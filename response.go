package quick

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"html"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	RequestId        uint64 // request id
	HttpRequest      *http.Request
	Status           string // e.g. "200 OK"
	StatusCode       int    // e.g. 200
	Proto            string // e.g. "HTTP/1.0"
	ProtoMajor       int    // e.g. 1
	ProtoMinor       int    // e.g. 0
	Header           http.Header
	Body             *bytes.Buffer
	ContentLength    int64
	ExecTime         time.Duration // request exec time
	TLS              *tls.ConnectionState
	TransferEncoding []string
	Encoding         encoding.Encoding // Response body encoding
	clientTrace      *clientTrace
}

func BuildResponse(resp *http.Response) (*Response, error) {
	if resp == nil {
		return nil, errors.New("http response is nil")
	}

	coding := unicode.UTF8 // HTML default encoding UTF8
	buffReader := bufio.NewReader(resp.Body)
	buff, err := buffReader.Peek(1024)
	if err == nil {
		coding, _, _ = charset.DetermineEncoding(buff, "")
	}

	reader := transform.NewReader(buffReader, coding.NewDecoder())
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &Response{
		HttpRequest:      resp.Request,
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Header:           CopyHeader(resp.Header),
		Body:             bytes.NewBuffer(body),
		ContentLength:    resp.ContentLength,
		TLS:              resp.TLS,
		TransferEncoding: resp.TransferEncoding,
		Encoding:         coding,
		clientTrace:      nil,
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
	return json.Unmarshal(r.Body.Bytes(), v)
}

func (r *Response) GetHtml() string {
	return html.UnescapeString(r.Body.String())
}

func (r *Response) GetXml(v interface{}) error {
	return xml.Unmarshal(r.Body.Bytes(), v)
}

func (r *Response) GetBody() []byte {
	return r.Body.Bytes()
}

// TraceInfo method returns the trace info for the request.
// If either the Client or Request EnableTrace function has not been called
// prior to the request being made, an empty TraceInfo object will be returned.
//
// Since v0.4.0
func (r *Response) TraceInfo() TraceInfo {
	ct := r.clientTrace
	if ct == nil {
		return TraceInfo{}
	}

	ti := TraceInfo{
		DNSLookup:     ct.dnsDone.Sub(ct.dnsStart),
		TLSHandshake:  ct.tlsHandshakeDone.Sub(ct.tlsHandshakeStart),
		ServerTime:    ct.gotFirstResponseByte.Sub(ct.gotConn),
		IsConnReused:  ct.gotConnInfo.Reused,
		IsConnWasIdle: ct.gotConnInfo.WasIdle,
		ConnIdleTime:  ct.gotConnInfo.IdleTime,
	}

	// Calculate the total time accordingly,
	// when connection is reused
	if ct.gotConnInfo.Reused {
		ti.TotalTime = ct.endTime.Sub(ct.getConn)
	} else {
		ti.TotalTime = ct.endTime.Sub(ct.dnsStart)
	}

	// Only calculate on successful connections
	if !ct.connectDone.IsZero() {
		ti.TCPConnTime = ct.connectDone.Sub(ct.dnsDone)
	}

	// Only calculate on successful connections
	if !ct.gotConn.IsZero() {
		ti.ConnTime = ct.gotConn.Sub(ct.getConn)
	}

	// Only calculate on successful connections
	if !ct.gotFirstResponseByte.IsZero() {
		ti.ResponseTime = ct.endTime.Sub(ct.gotFirstResponseByte)
	}

	// Capture remote address info when connection is non-nil
	if ct.gotConnInfo.Conn != nil {
		ti.RemoteAddr = ct.gotConnInfo.Conn.RemoteAddr()
	}

	return ti
}

func (r *Response) String() string {
	return r.Body.String()
}

func (r *Response) Read(p []byte) (n int, err error) {
	return r.Body.Read(p)
}
