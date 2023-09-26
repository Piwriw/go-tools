package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

/*
go get github.com/fsnotify/fsnotify
viper 可以监听文件修改进而自动重新加载。 其内部使用的就是fsnotify这个库，它是跨平台的
Name表示发生变化的文件或目录名，Op表示具体的变化
Chmod事件在文件或目录的属性发生变化时
*/
func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed:", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case events, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("%s %s \n", events.Name, events.Op)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("./")
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
}
