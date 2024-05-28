package queue

import "container/list"

type likedDeque struct {
	data *list.List
}

func newLinkedListDeque() *likedDeque {
	return &likedDeque{
		data: list.New(),
	}
}
