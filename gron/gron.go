package main

import (
	"fmt"
	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
	"sync"
	"time"
)

/*
go get github.com/roylee0704/gron
定时任务管理库
*/
type CustomJob struct {
	Name string
}

// 通过实现Run 实现自定义Job
func (j *CustomJob) Run() {
	fmt.Println("Hello ", j.Name)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	// 创建一个任务管理器
	cron := gron.New()
	cron.AddFunc(gron.Every(5*time.Second), func() {
		fmt.Println("runs every 5s")
	})
	cron.AddFunc(gron.Every(10*xtime.Second), func() {
		fmt.Println("runs every 10s")
	})
	t, _ := time.ParseDuration("1m10s")
	cron.AddFunc(gron.Every(t), func() {
		fmt.Println("runs every 1 minutes 10 seconds.")
	})

	myJob := &CustomJob{Name: "CustomJob"}
	cron.Add(gron.Every(5*time.Second), myJob)
	// 指定时间触发
	//通过gron.Every()设置每隔多长时间执行一次任务。对于大于 1 天的时间间隔，我们还可以使用gron.Every().At()
	cron.Start()
	wg.Wait()
}
