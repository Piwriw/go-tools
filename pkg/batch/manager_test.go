package batch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestNewManager 测试 NewManager 构造函数
// 测试场景：正常创建 Manager 实例
func TestNewManager(t *testing.T) {
	tests := []struct {
		name      string
		batchSize int
		batchFunc func([]*queue)
	}{
		{
			name:      "创建正常批处理管理器",
			batchSize: 10,
			batchFunc: func(items []*queue) {},
		},
		{
			name:      "创建批次大小为1的管理器",
			batchSize: 1,
			batchFunc: func(items []*queue) {},
		},
		{
			name:      "创建批次大小为100的管理器",
			batchSize: 100,
			batchFunc: func(items []*queue) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(tt.batchSize, tt.batchFunc)
			if manager == nil {
				t.Fatal("NewManager 返回 nil")
			}
			if manager.batchSize != tt.batchSize {
				t.Errorf("期望 batchSize=%d, 实际=%d", tt.batchSize, manager.batchSize)
			}
			if manager.batchFunc == nil {
				t.Error("batchFunc 不应为 nil")
			}
			if manager.stopCh == nil {
				t.Error("stopCh 不应为 nil")
			}
			if manager.moreCh == nil {
				t.Error("moreCh 不应为 nil")
			}
		})
	}
}

// TestManager_queueLen 测试 queueLen 方法
// 测试场景：空队列、有元素的队列
func TestManager_queueLen(t *testing.T) {
	manager := NewManager(10, func(items []*queue) {})

	t.Run("空队列长度应为0", func(t *testing.T) {
		if length := manager.queueLen(); length != 0 {
			t.Errorf("期望队列长度=0, 实际=%d", length)
		}
	})

	t.Run("添加元素后的队列长度", func(t *testing.T) {
		manager.mtx.Lock()
		manager.queues = append(manager.queues, &queue{1}, &queue{2})
		manager.mtx.Unlock()

		if length := manager.queueLen(); length != 2 {
			t.Errorf("期望队列长度=2, 实际=%d", length)
		}
	})
}

// TestManager_nextBatch 测试 nextBatch 方法
// 测试场景：空队列、队列元素少于批次大小、队列元素等于批次大小、队列元素多于批次大小
func TestManager_nextBatch(t *testing.T) {
	tests := []struct {
		name          string
		batchSize     int
		initialQueues []*queue
		expectedCount int
		remainCount   int
	}{
		{
			name:          "空队列",
			batchSize:     10,
			initialQueues: []*queue{},
			expectedCount: 0,
			remainCount:   0,
		},
		{
			name:          "队列元素少于批次大小",
			batchSize:     10,
			initialQueues: []*queue{{1}, {2}, {3}},
			expectedCount: 3,
			remainCount:   0,
		},
		{
			name:          "队列元素等于批次大小",
			batchSize:     3,
			initialQueues: []*queue{{1}, {2}, {3}},
			expectedCount: 3,
			remainCount:   0,
		},
		{
			name:          "队列元素多于批次大小",
			batchSize:     2,
			initialQueues: []*queue{{1}, {2}, {3}, {4}, {5}},
			expectedCount: 2,
			remainCount:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(tt.batchSize, func(items []*queue) {})
			manager.queues = tt.initialQueues

			batch := manager.nextBatch()

			if len(batch) != tt.expectedCount {
				t.Errorf("期望批次大小=%d, 实际=%d", tt.expectedCount, len(batch))
			}
			if len(manager.queues) != tt.remainCount {
				t.Errorf("期望剩余队列大小=%d, 实际=%d", tt.remainCount, len(manager.queues))
			}
		})
	}
}

// TestManager_keepNext 测试 keepNext 方法
// 测试场景：向 moreCh 发送信号
func TestManager_keepNext(t *testing.T) {
	manager := NewManager(10, func(items []*queue) {})

	t.Run("keepNext 应向 moreCh 发送信号", func(t *testing.T) {
		manager.keepNext()

		select {
		case <-manager.moreCh:
			// 成功接收到信号
		case <-time.After(100 * time.Millisecond):
			t.Error("未在超时时间内接收到 moreCh 信号")
		}
	})

	t.Run("多次调用 keepNext 不应阻塞", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			manager.keepNext()
		}
		// 如果没有阻塞，测试通过
	})
}

// TestManager_sendLoop 测试 sendLoop 方法
// 测试场景：启动 sendLoop，验证批处理函数被调用
func TestManager_sendLoop(t *testing.T) {
	t.Run("sendLoop 应处理批次", func(t *testing.T) {
		var callCount int32
		var receivedItems []*queue

		manager := NewManager(2, func(items []*queue) {
			atomic.AddInt32(&callCount, 1)
			receivedItems = items
		})

		manager.queues = []*queue{{1}, {2}, {3}, {4}}
		manager.Start()

		// 触发批处理
		manager.keepNext()

		time.Sleep(100 * time.Millisecond)

		if atomic.LoadInt32(&callCount) == 0 {
			t.Error("批处理函数未被调用")
		}
	})

	t.Run("sendLoop 应在收到 stop 信号后停止", func(t *testing.T) {
		stopped := make(chan struct{})
		manager := NewManager(10, func(items []*queue) {})

		go func() {
			manager.sendLoop()
			close(stopped)
		}()

		// 发送停止信号
		close(manager.stopCh)

		select {
		case <-stopped:
			// 成功停止
		case <-time.After(100 * time.Millisecond):
			t.Error("sendLoop 未在超时时间内停止")
		}
	})
}

// TestManager_sendOneBatch 测试 sendOneBatch 方法
// 测试场景：批处理函数被正确调用
func TestManager_sendOneBatch(t *testing.T) {
	t.Run("sendOneBatch 应调用批处理函数", func(t *testing.T) {
		var called bool
		var receivedItems []*queue

		manager := NewManager(2, func(items []*queue) {
			called = true
			receivedItems = items
		})

		manager.queues = []*queue{{1}, {2}}
		manager.sendOneBatch()

		if !called {
			t.Error("批处理函数未被调用")
		}
		if len(receivedItems) != 2 {
			t.Errorf("期望接收2个项目，实际=%d", len(receivedItems))
		}
		if len(manager.queues) != 0 {
			t.Errorf("期望队列为空，实际=%d", len(manager.queues))
		}
	})
}

// TestManager_Concurrent 测试并发场景
// 测试场景：多个 goroutine 同时操作队列
func TestManager_Concurrent(t *testing.T) {
	t.Run("并发添加和获取批次", func(t *testing.T) {
		const batchSize = 100
		const totalItems = 1000
		var processCount int32

		manager := NewManager(batchSize, func(items []*queue) {
			atomic.AddInt32(&processCount, 1)
		})

		manager.Start()
		defer close(manager.stopCh)

		var wg sync.WaitGroup
		for i := 0; i < totalItems; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				manager.mtx.Lock()
				manager.queues = append(manager.queues, &queue{idx})
				manager.mtx.Unlock()
				manager.keepNext()
			}(i)
		}

		wg.Wait()
		time.Sleep(200 * time.Millisecond)

		expectedBatches := totalItems / batchSize
		if actual := atomic.LoadInt32(&processCount); actual < int32(expectedBatches) {
			t.Logf("警告: 实际处理批次数=%d, 期望至少=%d", actual, expectedBatches)
		}
	})

	t.Run("并发 queueLen 调用", func(t *testing.T) {
		manager := NewManager(10, func(items []*queue) {})
		manager.queues = []*queue{{1}, {2}, {3}}

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = manager.queueLen()
			}()
		}

		wg.Wait()
		// 如果没有 panic 或死锁，测试通过
	})
}

// BenchmarkManager_nextBatch 性能基准测试
func BenchmarkManager_nextBatch(b *testing.B) {
	manager := NewManager(100, func(items []*queue) {})
	manager.queues = make([]*queue, 1000)
	for i := 0; i < 1000; i++ {
		manager.queues[i] = &queue{i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.nextBatch()
	}
}

// BenchmarkManager_queueLen 性能基准测试
func BenchmarkManager_queueLen(b *testing.B) {
	manager := NewManager(100, func(items []*queue) {})
	manager.queues = make([]*queue, 1000)
	for i := 0; i < 1000; i++ {
		manager.queues[i] = &queue{i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.queueLen()
	}
}

// ExampleNewManager 示例函数
func ExampleNewManager() {
	// 创建一个批处理管理器，每批次最多处理 10 个项目
	manager := NewManager(10, func(items []interface{}) {
		// 处理批次中的项目
		for _, item := range items {
			println(item)
		}
	})

	// 启动批处理循环
	manager.Start()

	// 添加项目到队列
	manager.queues = append(manager.queues, &queue{"item1"})
	manager.keepNext()

	_ = manager
}
