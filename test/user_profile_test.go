package main

import (
	"fmt"
	"testapi/conf"
	"testapi/req"
	"testing"
)

var suiteFile = "../test.json"

func TestGetUserSettings(t *testing.T) {
	suite, err := conf.Load(suiteFile)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	req.ExecRequest(suite, "01")
}
