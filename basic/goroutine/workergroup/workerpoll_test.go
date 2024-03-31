package workergroup

import "testing"

func TestWorkerPool(t *testing.T) {
	wm := NewWorkerManager(10)
	wm.StartWorkerPool()
}
