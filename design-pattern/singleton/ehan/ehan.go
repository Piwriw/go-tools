package ehan

import "fmt"

// 1、保证这个类非公有化，外界不能通过这个类直接创建一个对象
//
//	那么这个类就应该变得非公有访问 类名称首字母要小写
type singleton struct{}

// 2、但是还要有一个指针可以指向这个唯一对象，但是这个指针永远不能改变方向
//
//	Golang中没有常指针概念，所以只能通过将这个指针私有化不让外部模块访问
var instance *singleton = new(singleton)

func GetInstance() *singleton {
	return instance
}

func (s *singleton) SomeThing() {
	fmt.Println("单例对象的某方法")
}

func main() {
	s := GetInstance()
	s.SomeThing()
}
