package bhan

import "fmt"

type singleton struct{}

var instance *singleton

func GetInstance() *singleton {
	// 初始化，让这个要使用的第一次分配内存
	// 但是存在了多个goroutines 抢占 导致多实例
	if instance == nil {
		instance = new(singleton)
	}
	return instance
}

func (s *singleton) SomeThing() {
	fmt.Println("单例对象的某方法")
}

func main() {
	s := GetInstance()
	s.SomeThing()
}
