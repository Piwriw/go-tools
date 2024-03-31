package workergroup

import (
	"fmt"
	"time"
)

type worker struct {
	id  int
	err error
}

func (wk *worker) work(workerChan chan<- *worker) error {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				wk.err = err
			} else {
				wk.err = fmt.Errorf("Panic with happend with [%v]", r)
			}
		}
		workerChan <- wk
	}()
	fmt.Println("Start Worker...ID", wk.id)

	// 每个Worker睡眠一定时间之后，panic退出或者Goexit()退出
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
	}
	panic("worker panic...")
	return wk.err
}
