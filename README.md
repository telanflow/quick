# Quick

[![GoDoc](https://godoc.org/github.com/telanflow/quick?status.svg)](https://godoc.org/github.com/telanflow/quick)

Http request library for Go

快速简单的Http请求库

## examples

```go
package main

import (
    "fmt"
    "github.com/telanflow/quick"
    "net/http"
)

func main() {
	// quick.Post("example.com")
	// quick.PostFormData("example.com")
	// quick.Put("example.com")
	// quick.Head("example.com")
	// quick.Delete("example.com")
	// quick.Patch("example.com")
	// quick.Options("example.com")
	// quick.Trace("example.com")

    // https ssl skip verify 取消https验证
    quick.InsecureSkipVerify(true)

    // set header
    quick.SetHeader(http.Header{
        "Context-Type": []string{"text/html"},
    })
    
    // set UserAgent to request
    quick.SetHeaderSingle("User-Agent", "A go request libraries for quick")

    // or
    quick.SetUserAgent("A go request libraries for quick")
    
    // request
    resp, err := quick.Get(
        "example.com?bb=1", 
        quick.OptionQueryString("name=quick&aa=11"),// set Get params   eg. "example.com?bb=1&name=quick&aa=11"
        quick.OptionProxy("http://127.0.0.1:8080"), // set proxy 
        quick.OptionHeaderSingle("User-Agent", ""), // set http header
        quick.OptionHeader(http.Header{}),          // set http header  eg. http.Header || map[string]string || []string
        quick.OptionRedirectNum(5),                 // set redirect num
        // quick.OptionBody(""),                    // POST body
        // ... quick.Option
    )
    if err != nil {
        panic(err)
    }

	fmt.Println(resp)
}
```

### Session
```go
package main

import (
    "fmt"
    "github.com/telanflow/quick"
    "net/http"
    "time"
)

func main() {
    //session := quick.NewSession(&quick.SessionOptions{
    //    DialTimeout:           30 * time.Second,
    //    DialKeepAlive:         30 * time.Second,
    //    MaxConnsPerHost:       0,
    //    MaxIdleConns:          100,
    //    MaxIdleConnsPerHost:   2,
    //    IdleConnTimeout:       90 * time.Second,
    //    TLSHandshakeTimeout:   10 * time.Second,
    //    ExpectContinueTimeout: 1 * time.Second,
    //    DisableCookieJar:      false,
    //    DisableDialKeepAlives: false,
    //})

    cookieJar, err := quick.NewCookieJar()
    if err != nil {
        panic(err)
    }
    
    // quick use default global session

    // create new session
    session := quick.NewSession()

    // https ssl skip verify 取消https验证
    session.InsecureSkipVerify(true)

    // set cookieJar
    session.SetCookieJar(cookieJar) 

    resp, err := session.Get("example.com")
    if err != nil {
        panic(err)
    }
    
    // Other example:
    // 
    // create new request
    req := quick.NewRequest().SetUrl("example.com").SetMethod(http.MethodGet)

    // send request
    resp, err = session.Suck(
        req, 
        quick.OptionHeaderSingle("User-Agent", ""), // set http header
        // ... quick.Option
    )
    if err != nil {
        panic(err)
    }

    //resp.Status       e.g. "200 OK"
    //resp.StatusCode   e.g. 200
    //... 
    fmt.Println(resp)
}
```

## 
