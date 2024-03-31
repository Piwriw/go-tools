package workergroup

import "fmt"

type WorkerManager struct {
	workerChan chan *worker
	numWorker  int
}

func NewWorkerManager(numWorkers int) *WorkerManager {
	return &WorkerManager{
		numWorker:  numWorkers,
		workerChan: make(chan *worker, numWorkers),
	}
}
func (wm *WorkerManager) StartWorkerPool() {
	// 开启一定数量的worker
	for i := 0; i < wm.numWorker; i++ {
		wk := &worker{id: i}
		go wk.work(wm.workerChan)
	}
	wm.KeepLiveWorkers()
}
func (wm *WorkerManager) KeepLiveWorkers() {
	// 如果有worker dead ，workChan获取到这个worker，输出异常，重启
	for wk := range wm.workerChan {
		fmt.Printf("Worker %d stopped with err:[%v]\n", wk.id, wk.err)
		// reset err
		wk.err = nil
		go wk.work(wm.workerChan)
	}
}
