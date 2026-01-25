package batch

import "sync"

type Manager struct {
	queues    []*queue
	mtx       sync.RWMutex
	stopCh    chan struct{}
	moreCh    chan struct{}
	batchFunc func([]*queue)
	batchSize int
}

func NewManager(batchSize int, batchFunc func([]*queue)) *Manager {
	return &Manager{
		batchSize: batchSize,
		batchFunc: batchFunc,
		stopCh:    make(chan struct{}),
		moreCh:    make(chan struct{}, 1),
	}
}

func (n *Manager) Start() {
	go n.sendLoop()
}

type queue interface{}

func (n *Manager) queueLen() int {
	n.mtx.RLock()
	defer n.mtx.RUnlock()

	return len(n.queues)
}
func (n *Manager) nextBatch() []*queue {
	// 加锁保护队列数据，防止并发访问
	n.mtx.Lock()
	// 确保在函数退出时解锁
	defer n.mtx.Unlock()

	var queues []*queue

	// 检查队列长度是否超过最大批次大小
	if maxBatchSize := n.batchSize; len(n.queues) > n.batchSize {
		// 如果超过，只取前maxBatchSize个警报
		queues = append(make([]*queue, 0, maxBatchSize), n.queues[:maxBatchSize]...)
		// 从队列中移除已取出的警报
		n.queues = n.queues[maxBatchSize:]
	} else {
		// 如果不超过，取出全部警报
		queues = append(make([]*queue, 0, len(n.queues)), n.queues...)
		// 清空队列
		n.queues = n.queues[:0]
	}

	return queues
}

func (n *Manager) sendOneBatch() {
	queues := n.nextBatch()
	n.batchFunc(queues)
}
func (n *Manager) sendLoop() {
	for {
		select {
		case <-n.stopCh:
			return
		default:
			select {
			case <-n.stopCh:
				return

			case <-n.moreCh:
				n.sendOneBatch()
				if n.queueLen() > 0 {
					n.keepNext()
				}
			}
		}
	}
}

func (n *Manager) keepNext() {
	select {
	case n.moreCh <- struct{}{}:
	default:
	}
}
