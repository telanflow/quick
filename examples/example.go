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

	quick.SetHeader(http.Header{
		"Context-Type": []string{"text/html"},
	})

	// set UserAgent to request
	quick.SetHeaderSingle("UserAgent", "A go request libraries for quick")
	// or
	quick.SetUserAgent("A go request libraries for quick")

	resp, err := quick.Get("http://www.baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
