package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	DefaultClient *Client
)

func init() {
	DefaultClient = NewHTTPClient()
}

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

// SetContentType 设置默认的 Content-Type
func SetContentType(ct ContentType) {
	DefaultClient.contentType = ct
}

// SetClient 设置默认的 HTTP 客户端
func SetClient(client *http.Client) {
	DefaultClient.client = client
}

// SetTimeout 设置默认的超时时间
func SetTimeout(t time.Duration) {
	DefaultClient.client.Timeout = t
}

// SetIdleConnTimeout 设置默认的空闲连接超时时间
func SetIdleConnTimeout(t time.Duration) {
	if transport, ok := DefaultClient.client.Transport.(*http.Transport); ok {
		transport.IdleConnTimeout = t
	}
}

// SetMaxIdleConns 设置默认的最大空闲连接数
func SetMaxIdleConns(n int) {
	if transport, ok := DefaultClient.client.Transport.(*http.Transport); ok {
		transport.MaxIdleConns = n
	}
}

// SetTLSHandshakeTimeout 设置默认的 TLS 握手超时时间
func SetTLSHandshakeTimeout(t time.Duration) {
	if transport, ok := DefaultClient.client.Transport.(*http.Transport); ok {
		transport.TLSHandshakeTimeout = t
	}
}

// SetCheckRedirect 设置默认的重定向检查函数
func SetCheckRedirect(fn func(req *http.Request, via []*http.Request) error) {
	DefaultClient.client.CheckRedirect = fn
}

// NewHTTPClient creates a new HTTP client with default options.
func NewHTTPClient(opts ...Options) *Client {
	client := &Client{
		client: &http.Client{
			Timeout: defaultConnectTimeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 300 * time.Second,
				}).DialContext,
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

	for _, opt := range opts {
		opt(client)
	}
	return client
}

// Options 是一个用于配置 Client 的函数类型
type Options func(*Client)

// WithConnectTimeout 设置客户端连接超时时间
func WithConnectTimeout(t time.Duration) Options {
	return func(c *Client) {
		c.client.Timeout = t
	}
}

// WithTimeOut 设置客户端超时时间
func WithTimeOut(t time.Duration) Options {
	return func(c *Client) {
		c.client.Timeout = t
	}
}

// WithIdleConnTimeout 设置空闲连接超时时间
func WithIdleConnTimeout(t time.Duration) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.IdleConnTimeout = t
		}
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(n int) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.MaxIdleConns = n
		}
	}
}

// WithContentType 设置请求的 Content-Type
func WithContentType(ct ContentType) Options {
	return func(c *Client) {
		c.contentType = ct
	}
}

// WithTransport 设置自定义的 HTTP 传输层（Transport）
func WithTransport(transport http.RoundTripper) Options {
	return func(c *Client) {
		c.client.Transport = transport
	}
}

// WithCheckRedirect 设置重定向检查函数，用于控制是否允许跳转
func WithCheckRedirect(fn func(req *http.Request, via []*http.Request) error) Options {
	return func(c *Client) {
		c.client.CheckRedirect = fn
	}
}

// WithJar 设置 CookieJar，用于自动处理 Cookie
func WithJar(jar http.CookieJar) Options {
	return func(c *Client) {
		c.client.Jar = jar
	}
}

// WithTimeout 设置客户端整体请求的超时时间
func WithTimeout(t time.Duration) Options {
	return func(c *Client) {
		c.client.Timeout = t
	}
}

// WithTLSHandshakeTimeout 设置 TLS 握手的超时时间（仅在使用 http.Transport 时有效）
func WithTLSHandshakeTimeout(t time.Duration) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.TLSHandshakeTimeout = t
		}
	}
}

// WithResponseHeaderTimeout 设置服务器响应头的超时时间（仅在使用 http.Transport 时有效）
func WithResponseHeaderTimeout(t time.Duration) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.ResponseHeaderTimeout = t
		}
	}
}

// WithExpectContinueTimeout 设置 Expect: 100-continue 超时时间（仅在使用 http.Transport 时有效）
func WithExpectContinueTimeout(t time.Duration) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.ExpectContinueTimeout = t
		}
	}
}

// WithProxy 设置代理服务器（仅在使用 http.Transport 时有效）
func WithProxy(proxy func(*http.Request) (*url.URL, error)) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.Proxy = proxy
		}
	}
}

// WithDisableKeepAlives 设置是否禁用 HTTP 连接的 Keep-Alive 机制（仅在使用 http.Transport 时有效）
func WithDisableKeepAlives(disable bool) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.DisableKeepAlives = disable
		}
	}
}

// WithDisableCompression 设置是否禁用 HTTP 压缩（仅在使用 http.Transport 时有效）
func WithDisableCompression(disable bool) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.DisableCompression = disable
		}
	}
}

// WithMaxIdleConnsPerHost 设置每个主机的最大空闲连接数（仅在使用 http.Transport 时有效）
func WithMaxIdleConnsPerHost(n int) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.MaxIdleConnsPerHost = n
		}
	}
}

// WithProxyConnectHeader 设置代理服务器的连接请求头（仅在使用 http.Transport 时有效）
func WithProxyConnectHeader(header http.Header) Options {
	return func(c *Client) {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.ProxyConnectHeader = header
		}
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
