package main

import (
	"fmt"
	"github.com/telanflow/requests"
)

func main() {
	resp, err := requests.Get("http://www.baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
