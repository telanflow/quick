package main

import (
	"fmt"
	"requests"
)

func main() {

	resp, err := requests.Get("http://www.telan.me")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.ExecTime.Seconds())
}
