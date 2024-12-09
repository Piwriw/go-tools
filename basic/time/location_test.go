package main

import (
	"testing"
	"time"
)

func TestLocation(t *testing.T) {
	now := time.Now()
	t.Log("now:", now)
	t.Log("now:", now.UTC())
	t.Log("now:", now.UTC().Format("2006-01-02 15:04:05"))
	t.Log("now:", now.In(time.FixedZone("Asia/Shanghai", 8*3600)).Format("2006-01-02 15:04:05"))
}
