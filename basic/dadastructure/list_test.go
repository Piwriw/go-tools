package main

import (
	"container/list"
	"testing"
)

func TestList(t *testing.T) {
	list := list.New()
	list.PushFront(1)
	list.PushFront(2)
	list.PushBack(3)
}
