package bhan

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type safeSingleton struct{}

// 定义锁
var lock sync.Mutex
var initialized uint32

var safeInstance *safeSingleton

func GetSafeInstance() *safeSingleton {
	// 初始化，让这个要使用的第一次分配内存
	// 通过锁实现 线程安全，通过atomic实现标记，防止重复上锁，浪费性能
	//如果标记为被设置，直接返回，不加锁
	if atomic.LoadUint32(&initialized) == 1 {
		return safeInstance
	}

	//如果没有，则加锁申请
	lock.Lock()
	defer lock.Unlock()

	if initialized == 0 {
		safeInstance = new(safeSingleton)
		//设置标记位
		atomic.StoreUint32(&initialized, 1)
	}

	return safeInstance
}

func (s *safeSingleton) SomeThing() {
	fmt.Println("单例对象的某方法")
}

func main() {
	s := GetSafeInstance()
	s.SomeThing()
}
