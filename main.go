package main

import (
	"fmt"
	"os"
	"testapi/conf"
)

func main() {
	suite, err := conf.Load(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(suite.Cases[os.Args[2]])
}
