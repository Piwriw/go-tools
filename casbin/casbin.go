package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"log"
)

/*
go get github.com/casbin/casbin/v2
WEB URL:https://casbin.org/zh/
*/
func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, err := e.Enforce(sub, obj, act)
	if err != nil {
		fmt.Println(err)
		return
	}
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}
func main() {
	e, err := casbin.NewEnforcer("./casbin/model.conf", "./casbin/policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}
	//
	check(e, "piwriw", "data1", "read")
	check(e, "joohwan", "data2", "write")
	check(e, "piwriw", "data1", "write")
	check(e, "joohwan", "data2", "read")
	check(e, "root", "data2", "read")
	check(e, "piwriw2", "data3", "read")
	check(e, "piwriw2", "data3", "write")

}
