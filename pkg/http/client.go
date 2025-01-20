package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultConnectTimeout = 300 * time.Second
	defaultContentType    = JSON
)

type ContentType string

const (
	JSON ContentType = "application/json;charset=utf-8"
)

type Client struct {
	contentType ContentType
	client      *http.Client
}

func NewHTTPClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: defaultConnectTimeout,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 300 * time.Second,
				}).Dial,
				DisableKeepAlives:     true,
				MaxIdleConnsPerHost:   5,
				MaxIdleConns:          5,
				IdleConnTimeout:       5 * time.Second,
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 300 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
			},
		},
	}
}

func (h *Client) JSON() *Client {
	h.contentType = JSON
	return h
}

func (h *Client) Post(url string, data []byte) ([]byte, error) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	if h.contentType == "" {
		h.contentType = defaultContentType
	}
	request.Header.Set("Content-Type", string(h.contentType))

	res, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// 读取响应
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// PostFile sends a file in a POST request using multipart/form-data
func (h *Client) PostFile(url, fieldName, filePath string) ([]byte, error) {
	// 创建一个缓冲区和 multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建文件字段
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	// 将文件内容写入 multipart
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// 关闭 multipart writer 以完成请求体
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// 创建 HTTP 请求
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// 设置 Content-Type 为 multipart/form-data
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	res, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
