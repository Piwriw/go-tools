package main

import (
	"flag"
	"fmt"
)

var helps string

func main() {
	flag.StringVar(&helps, "help", "default", "string flag value")
	flag.Parse()
	fmt.Println(helps)
}
