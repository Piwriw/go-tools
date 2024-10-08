package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func warpPrint(err error) {
	fmt.Println(errors.Wrap(err, "no node need to be have some tasks"))
}

func main() {
	err := errors.New("sss")
	warpPrint(err)
	fmt.Println("all", errors.Cause(err))
}
