package http

import (
	"bytes"
	"io"
	"net"
	"net/http"
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
