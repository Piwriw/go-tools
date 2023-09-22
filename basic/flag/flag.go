package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help       bool
	version, V bool
	q          *bool
	s          string
	p          string
	c          string
	g          string
)

func FlagInit() {
	flag.BoolVar(&help, "h", false, "this is help")
	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&V, "V", false, "show version moreinfo")

	flag.StringVar(&s, "s", "", "send `signal` to a master process: stop, quit, reopen, reload")
	flag.StringVar(&p, "p", "/usr/local/nginx/", "set `prefix` path")
	flag.StringVar(&c, "c", "conf/nginx.conf", "set configuration `file`")
	flag.Usage = usage
}

func main() {
	FlagInit()
	flag.Parse()
	if help {
		flag.Usage()
	}
}
func usage() {
	fmt.Fprintf(os.Stderr, `nginx version: nginx/1.10.0
Usage: nginx [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]
Options:
`)
	flag.PrintDefaults()
}
