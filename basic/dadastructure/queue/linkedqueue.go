package queue

import "container/list"

type Queue interface {
	push(value int)
	pop() any
	peek(value int)
	toList()
	isEmpty() bool
}

type linkedQueue struct {
	data *list.List
}

func NewLinkedQueue() *linkedQueue {
	return &linkedQueue{
		data: list.New(),
	}
}
func (q linkedQueue) push(value int) {
	q.data.PushBack(value)
}
func (q linkedQueue) pop() any {
	if q.isEmpty() {
		return nil
	}
	e := q.data.Front()
	q.data.Remove(e)
	return e.Value
}
func (q linkedQueue) peek(value int) any {
	if q.isEmpty() {
		return nil
	}
	return q.data.Front().Value
}
func (q linkedQueue) isEmpty() bool {
	return q.data.Len() == 0
}

/* 获取 List 用于打印 */
func (s *linkedQueue) toList() *list.List {
	return s.data
}
