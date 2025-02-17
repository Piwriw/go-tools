package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

var cnt int

func uploadFile(url, deviceID, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 创建一个缓冲区和 multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加 device_id 字段
	err := writer.WriteField("device_id", deviceID)
	if err != nil {
		fmt.Printf("Error adding device_id: %v\n", err)
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// 添加 img 文件字段
	part, err := writer.CreateFormFile("img", file.Name())
	if err != nil {
		fmt.Printf("Error adding file field: %v\n", err)
		return
	}

	// 将文件内容写入 multipart
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		return
	}

	// 关闭 multipart writer 以完成请求体
	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing writer: %v\n", err)
		return
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// 设置 Content-Type 为 multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}
	cnt++
	// 打印响应
	fmt.Printf("Response: %s\n %d", respBody, cnt)
}

func main() {
	err := errors.New("aa")
	slog.Info("xxx", slog.Any("err", err), "sss", "xx")
}
