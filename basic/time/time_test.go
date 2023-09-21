package main

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Logf("Now:%v\n", time.Now())
	t.Logf("Unix:%v\n", time.Now().Unix())
	format := time.Now().Format("2006-04-02 15-04-05")
	t.Log(format)
	location, err := time.ParseInLocation("2006-04-02 15-04-05", time.Now().Format("2006-04-02 15-04-05"), time.Local)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(location)
}
