package utils

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

// Define constant values for timeouts (default values)
const (
	defaultConnectTimeout        = 300 * time.Second
	defaultMaxIdleConn           = 5
	defaultIdleConnTimeout       = 5 * time.Second
	defaultResponseTimeout       = 300 * time.Second
	defaultTLSHandshakeTimeout   = 5 * time.Second
	defaultExpectContinueTimeout = 5 * time.Second
)

// HTTPClientConfig holds the configuration options for the HTTP client and transport
type HTTPClientConfig struct {
	ConnectTimeout        time.Duration
	MaxIdleConn           int
	MaxIdleConnPerHost    int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	ExpectContinueTimeout time.Duration
}

// DefaultHTTPClientConfig returns a default configuration with predefined values
func DefaultHTTPClientConfig() *HTTPClientConfig {
	return &HTTPClientConfig{
		ConnectTimeout:        defaultConnectTimeout,
		MaxIdleConn:           defaultMaxIdleConn,
		MaxIdleConnPerHost:    defaultMaxIdleConn,
		IdleConnTimeout:       defaultIdleConnTimeout,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ResponseHeaderTimeout: defaultResponseTimeout,
		ExpectContinueTimeout: defaultExpectContinueTimeout,
	}
}

// Global HTTP client (singleton)
var (
	client     *http.Client
	clientOnce sync.Once
)

// initClient initializes the HTTP client with custom configurations
func initClient(config *HTTPClientConfig) {
	client = &http.Client{
		Timeout: config.ConnectTimeout,
		Transport: &http.Transport{
			// Using DialContext to replace Dial
			DialContext: (&net.Dialer{
				Timeout: 300 * time.Second,
			}).DialContext,
			DisableKeepAlives:     true,
			MaxIdleConnsPerHost:   config.MaxIdleConnPerHost,
			MaxIdleConns:          config.MaxIdleConn,
			IdleConnTimeout:       config.IdleConnTimeout,
			TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
			ResponseHeaderTimeout: config.ResponseHeaderTimeout,
			ExpectContinueTimeout: config.ExpectContinueTimeout,
		},
	}
}

// DoHTTPRequest sends a POST request with the provided data to the specified URL
func DoHTTPRequest(data []byte, url string, config *HTTPClientConfig) error {
	// Ensure that the client is initialized only once
	clientOnce.Do(func() {
		initClient(config)
	})
	slog.Info("HTTP Info", slog.String("url", url), slog.String("data", string(data)))
	// Create the HTTP request and check for errors
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set appropriate headers
	request.Header.Set("Content-Type", "application/json;charset=utf-8")

	// Perform the HTTP request using the global client
	res, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer res.Body.Close()

	// Check for non-2xx status codes
	// Check for non-2xx status codes
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// Read response body for additional error context
		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return fmt.Errorf("HTTP request failed with status code %d, but failed to read response body: %w", res.StatusCode, readErr)
		}

		return fmt.Errorf("HTTP request failed with status code %d, response: %s", res.StatusCode, string(body))
	}

	return nil
}
