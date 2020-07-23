# Quick

[![GoDoc](https://godoc.org/github.com/telanflow/quick?status.svg)](https://godoc.org/github.com/telanflow/quick)

Simple and efficient HTTP request library

简单高效的Http请求库

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
    quick.SetUserAgent("A go request libraries for quick")
    
    // You should init it by using NewCookiesWithString like this:
    // 	cookies := quick.NewCookiesWithString(
    //		"key1=value1; key2=value2; key3=value3"
    // 	)
    // Note: param is cookie string
    cookies := quick.NewCookiesWithString("sessionid=11111")
    
    // request
    resp, err := quick.Get(
        "http://example.com?bb=1", 
        quick.OptionQueryString("name=quick&aa=11"),   // set Get params   eg. "example.com?bb=1&name=quick&aa=11"
        quick.OptionHeaderSingle("User-Agent", ""),    // set http header
        quick.OptionHeader(http.Header{}),             // set http header  eg. http.Header || map[string]string || []string
        quick.OptionRedirectNum(10),                   // set redirect num
        quick.OptionCookies(cookies),                  // set cookies to request
        // quick.OptionProxy("http://127.0.0.1:8080"), // set proxy address
        // quick.OptionBody(""),                       // POST body
        // quick.OptionBasicAuth("username", "password"), // HTTP Basic Authentication
        // ... quick.Option
    )
    if err != nil {
        panic(err)
    }

    fmt.Println(resp)
}
```

### Session (会话)

Request based session

所有Request都基于session（http.Client）

```go
func main() {
    // cookieJar
    cookieJar, err := quick.NewCookieJar()
    if err != nil {
        panic(err)
    }
    
    // quick use default global session
    // create session
    session := quick.NewSession()
    // https ssl skip verify 取消https验证
    session.InsecureSkipVerify(true)
    // set cookieJar
    session.SetCookieJar(cookieJar) 
    resp, err := session.Get("http://example.com")
    if err != nil {
        panic(err)
    }
    //resp.Status       e.g. "200 OK"
    //resp.StatusCode   e.g. 200
    //... 
    fmt.Println(resp)
}
```


Other example:
```go
func main() {
    // new Request
    req := quick.NewRequest().SetUrl("http://example.com").SetMethod(http.MethodGet)

    // send Request
    session := quick.NewSession()
    resp, err := session.Suck(
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

## Middleware（中间件）
```go
func main() {
    // new Request
    req := quick.NewRequest().SetUrl("http://example.com").SetMethod(http.MethodGet)

    // create session
    session := quick.NewSession()

    // use middleware
    session.Use(
        // middleware 1
        func(r *http.Request) {
            log.Printf(
                "Middleware: %v RedirectNum: %v Proxy: %v \n",
                r.URL,
                r.Context().Value(quick.ContextRedirectNumKey),
                r.Context().Value(quick.ContextProxyKey),
            )
        },
        // middleware 2
        func(r *http.Request) {
            log.Printf(
                "Middleware2: %v RedirectNum: %v Proxy: %v \n",
                r.URL,
                r.Context().Value(quick.ContextRedirectNumKey),
                r.Context().Value(quick.ContextProxyKey),
            )
        },
    )

    // send Request
    resp, err := session.Suck(
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

## License

[MIT](LICENSE)