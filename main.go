package main

import (
	"fmt"
	"log"
	"testapi/state"
)

func main() {
	token, err := state.GetToken("test01")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
	// suite, err := conf.Load(os.Args[1])
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// req.ExecRequest(suite, os.Args[2])
}
