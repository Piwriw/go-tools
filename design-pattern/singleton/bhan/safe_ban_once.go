package bhan

import (
	"fmt"
	"sync"
)

/*
源码中的Once 其实就是safe_bhan
*/
var once sync.Once

type safeOnceSingleton struct{}

// 定义锁
var safeOnceInstance *safeOnceSingleton

func GetSafeOnceInstance() *safeOnceSingleton {
	once.Do(func() {
		safeOnceInstance = new(safeOnceSingleton)
	})
	return safeOnceInstance
}

func (s *safeOnceSingleton) SomeThing() {
	fmt.Println("单例对象的某方法")
}

//单例模式的优缺点
//优点：
//(1) 单例模式提供了对唯一实例的受控访问。
//(2) 节约系统资源。由于在系统内存中只存在一个对象。
//
//缺点：
//(1) 扩展略难。单例模式中没有抽象层。
//(2) 单例类的职责过重。

// 适用场景
// (1) 系统只需要一个实例对象，如系统要求提供一个唯一的序列号生成器或资源管理器，或者需要考虑资源消耗太大而只允许创建一个对象。
// (2) 客户调用类的单个实例只允许使用一个公共访问点，除了该公共访问点，不能通过其他途径访问该实例
func main() {
	s := GetSafeOnceInstance()
	s.SomeThing()
}
