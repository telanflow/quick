package main

import (
	"fmt"
	"github.com/telanflow/quick"
	"log"
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

	// use middleware
	quick.Use(
		func(r *http.Request) {
			log.Printf(
				"Middleware: %v RedirectNum: %v Proxy: %v \n",
				r.URL,
				r.Context().Value(quick.ContextRedirectNumKey),
				r.Context().Value(quick.ContextProxyKey),
			)
		},

		func(r *http.Request) {
			log.Printf(
				"Middleware2: %v RedirectNum: %v Proxy: %v \n",
				r.URL,
				r.Context().Value(quick.ContextRedirectNumKey),
				r.Context().Value(quick.ContextProxyKey),
			)
		},
	)

	// You should init it by using NewCookiesWithString like this:
	// 	cookies := quick.NewCookiesWithString(
	//		"key1=value1; key2=value2; key3=value3"
	// 	)
	// Note: param is cookie string
	cookies := quick.NewCookiesWithString("sessionid=11111")

	// request
	resp, err := quick.Get(
		"http://www.baidu.com?bb=1",
		quick.OptionQueryString("name=quick&aa=11"),   // set Get params   eg. "example.com?bb=1&name=quick&aa=11"
		//quick.OptionProxy("http://127.0.0.1:8080"),  // set proxy
		//quick.OptionHeaderSingle("User-Agent", ""),  // set http header
		//quick.OptionHeader(http.Header{}),           // set http header  eg. http.Header || map[string]string || []string
		//quick.OptionRedirectNum(10),                 // set redirect num
		quick.OptionCookies(cookies), // set cookies to request
		// quick.OptionBody(""),                       // POST body
		// quick.OptionBasicAuth("username", "password"), // HTTP Basic Authentication
		// ... quick.Option
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.ExecTime)
}
