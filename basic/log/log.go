package main

import (
	"bytes"
	"fmt"
	"log"
)

/*
Print/Printf/Println：正常输出日志；
Panic/Panicf/Panicln：输出日志后，以拼装好的字符串为参数调用panic；
Fatal/Fatalf/Fatalln：输出日志后，调用os.Exit(1)退出程序。
*/

/*
SetPrefix
Ldate：输出当地时区的日期，如2020/02/07；
Ltime：输出当地时区的时间，如11:45:45；
Lmicroseconds：输出的时间精确到微秒，设置了该选项就不用设置Ltime了。如11:45:45.123123；
Llongfile：输出长文件名+行号，含包名，如github.com/darjun/go-daily-lib/log/flag/main.go:50；
Lshortfile：输出短文件名+行号，不含包名，如main.go:50；
LUTC：如果设置了Ldate或Ltime，将输出 UTC 时间，而非当地时区。
*/
func SetPrefix() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

}
func NewLogger() {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", log.Lshortfile|log.LstdFlags)
	logger.Printf("loggerName:%s", "mylogger")
	fmt.Println(buf.String())
}
func main() {
	NewLogger()
}
