# Quick
Http request library for Go

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

    // set header
    quick.SetHeader(http.Header{
        "Context-Type": []string{"text/html"},
    })
    
    // set UserAgent to request
    quick.SetHeaderSingle("UserAgent", "A go request libraries for quick")

    // or
    quick.SetUserAgent("A go request libraries for quick")

	resp, err := quick.Get("example.com")
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
    
    session := quick.NewSession()
    session.SetCookieJar(cookieJar)
	resp, err := session.Get("example.com")
	if err != nil {
		panic(err)
	}
    
    //resp.Status       e.g. "200 OK"
    //resp.StatusCode   e.g. 200
    //...
	fmt.Println(resp)
}
```
