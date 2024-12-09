package main

import "fmt"

type Subject interface {
	Subscribe(observer Observer)
	Notify(msg string)
}

type Observer interface {
	Update(string)
}

type SubjectImpl struct {
	subscribers []Observer
}

func (s *SubjectImpl) Subscribe(observer Observer) {
	s.subscribers = append(s.subscribers, observer)
}

func (s *SubjectImpl) Notify(msg string) {
	for _, observer := range s.subscribers {
		observer.Update(msg)
	}
}

// Observer1 Observer1
type Observer1 struct{}

// Update 实现观察者接口
func (Observer1) Update(msg string) {
	fmt.Printf("Observer1: %s\n", msg)
}

// Observer2 Observer2
type Observer2 struct{}

// Update 实现观察者接口
func (Observer2) Update(msg string) {
	fmt.Printf("Observer2: %s\n", msg)
}

/*
Observer1: Hello World
Observer2: Hello World
*/
func main() {
	sub := &SubjectImpl{}
	sub.Subscribe(&Observer1{})
	sub.Subscribe(&Observer2{})
	sub.Notify("Hello World")
}
