package main

import (
	"bufio"
	"fmt"
	"os"
)

func write() {
	file, err := os.Create("./basic/bufio/output.txt")
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// 写入数据到缓冲区
	_, err = writer.WriteString("Hello, World!")
	if err != nil {
		fmt.Println("写入数据失败:", err)
		return
	}

	// 刷新缓冲区到文件
	err = writer.Flush()
	if err != nil {
		fmt.Println("刷新缓冲区失败:", err)
		return
	}

	fmt.Println("数据写入完成。")
}
func read() {
	file, err := os.Open("./basic/bufio/output.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
func main() {
	write()
	read()

}
