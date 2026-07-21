package main

import (
	"fmt"
	"testapi/auth"
)

func main() {
	token, err := auth.GetSystemToken("sys_3.23")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(token)
	}
}
