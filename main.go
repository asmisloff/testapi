package main

import (
	"fmt"
	"testapi/state"
)

func main() {
	token, err := state.GetToken("test01")
	if err != nil {
		fmt.Println(err)
		state.SaveState("test01", &state.StateEntry{Token: "12345", RefreshToken: "54321"})
	}
	fmt.Println(token)
	// suite, err := conf.Load(os.Args[1])
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// req.ExecRequest(suite, os.Args[2])
}
