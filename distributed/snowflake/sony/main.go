package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sony/sonyflake"
)

func main() {
	// 初始化 Sonyflake，手动设置 Machine ID（可选）
	sf := sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) { return 1, nil },
	})
	if sf == nil {
		log.Fatal("Sonyflake not created")
	}

	// 生成 Sonyflake ID
	id, err := sf.NextID()
	if err != nil {
		log.Fatal("Failed to generate ID:", err)
	}

	// 解析 ID
	decomposed := sonyflake.Decompose(id)

	// Sonyflake 的时间戳基准
	startTime := time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC)

	// Sonyflake 的时间戳是 10ms 级别，需要乘以 10ms
	idTime := startTime.Add(time.Duration(decomposed["time"]) * 10 * time.Millisecond)

	// 打印信息
	fmt.Println("Sonyflake ID:", id)
	fmt.Println("时间戳:", decomposed["time"])
	fmt.Println("实际时间:", idTime)
	fmt.Println("机器 ID:", decomposed["machine-id"])
	fmt.Println("序列号:", decomposed["sequence"])
}
