package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func ReadFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}
func StuRead() {
	// 控制台读取
	data, err := ReadFrom(os.Stdin, 11)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", data)
}

// ReadAt 接口使得可以从指定偏移量处开始读取数据
func ReadAt() {
	reader := strings.NewReader("Go is good")
	p := make([]byte, 6)
	n, err := reader.ReadAt(p, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s,%d\n", p, n)
}

// WriteAt  在文件流的 offset=n 处写入 内容
func WriteAt() {
	file, err := os.Create("writeAt.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	n, err := file.WriteString("Go is writing something")
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
	file.WriteAt([]byte("ssss"), 2)
}

// TestReadio 使用NewWriter
func TestReadio(t *testing.T) {
	file, err := os.Open("writeAt.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(os.Stdout)
	writer.ReadFrom(file)
	writer.Flush()
}

/*
 PipeReader 和 PipeWriter 类型
 管道读法
*/
// PipeWrite Write
func PipeWrite(write *io.PipeWriter) {
	data := []byte("sending pipe")
	for i := 0; i < 2; i++ {
		n, err := write.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("write bytes:%d", n)
	}
	write.CloseWithError(errors.New("writes is closed"))
}

// PipeRead  Read
func PipeRead(reader *io.PipeReader) {
	buf := make([]byte, 128)
	for {
		fmt.Println("Read will waiting for 5s")
		time.Sleep(5 * time.Second)
		fmt.Println("Reader is ok,will starting")
		n, err := reader.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("bytes:%d \n res:%s", n, buf)
	}
}
func TestPipe(t *testing.T) {
	reader, writer := io.Pipe()
	go PipeWrite(writer)
	go PipeRead(reader)
	time.Sleep(30 * time.Second)
}
func TestCopy(t *testing.T) {
	io.Copy(os.Stdout, strings.NewReader("Go is good"))
	io.CopyN(os.Stdout, strings.NewReader("Go语言中文网"), 8)

}
