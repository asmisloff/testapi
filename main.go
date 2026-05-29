package main

import (
	"fmt"
	"os"
	"testapi/conf"
	"testapi/req"
)

func main() {
	suite, err := conf.Load(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	req.ExecRequest(suite, os.Args[2])
}
