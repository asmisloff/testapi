package main

import (
	"fmt"
	"os"
	"testapi/conf"
	"testapi/req"
)

func main() {
	conf.Foo()
	args := os.Args
	status, body, err := req.Get(args[1], nil)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(status)
		fmt.Println(body)
	}
}
