package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestStrconv(t *testing.T) {
	n, err := strconv.ParseInt("128", 10, 8)
	t.Log(n, err)

	/*
		strconv.Itoa(i) 性能更好
		Sprintf 传入interface 有反射过程
	*/
	startTime := time.Now()
	for i := 0; i < 10000; i++ {
		fmt.Sprintf("%d", i)
	}

	fmt.Println(time.Now().Sub(startTime))
	startTime = time.Now()
	for i := 0; i < 10000; i++ {
		strconv.Itoa(i)
	}
	fmt.Println(time.Now().Sub(startTime))
}

/*
	输出一些特殊富豪

This is “studygolang.com” website.
*/
func TestPrintSpe(t *testing.T) {
	t.Log(fmt.Println(`This is "studygolang.com" website`))
	t.Log(fmt.Println("This is \"studygolang.com\" website"))
	t.Log(fmt.Println("This is", strconv.Quote("studygolang.com"), "website"))
}
