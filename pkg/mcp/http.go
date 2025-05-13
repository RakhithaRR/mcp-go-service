package mcp

import (
	"crypto/tls"
	"fmt"
	"mcp-server/pkg/service"
	"net/http"
	"sync"
	"time"
)

var (
	syncOnce        sync.Once
	httpClient      *MCPHTTPClient
	skipVerifying   bool
	maxIdleConns    int
	idleConnTimeout int
)

type MCPHTTPClient struct {
	httpClient *http.Client
	UserAgent  string
}

func InitHttpClient() *MCPHTTPClient {
	skipVerifying = service.GetConfig().Http.Insecure
	maxIdleConns = service.GetConfig().Http.MaxIdleConns
	idleConnTimeout = service.GetConfig().Http.IdleConnTimeout
	syncOnce.Do(func() {
		if httpClient == nil {
			client := http.Client{
				Transport: &http.Transport{
					MaxIdleConns:    maxIdleConns,
					IdleConnTimeout: time.Duration(idleConnTimeout) * time.Second,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: skipVerifying,
					},
				},
			}
			httpClient = &MCPHTTPClient{
				httpClient: &client,
				UserAgent:  "Bijira-MCP-Client-Go/0.1",
			}
		}
	})
	return httpClient
}

func (client *MCPHTTPClient) DoRequest(request *http.Request) (*http.Response, error) {
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (client *MCPHTTPClient) GenerateRequest(httpRequest *TransformedRequest) (*http.Request, error) {
	var req *http.Request
	var err error
	if httpRequest.Body == nil {
		req, err = http.NewRequest(httpRequest.Method, httpRequest.URL, nil)
	} else {
		req, err = http.NewRequest(httpRequest.Method, httpRequest.URL, httpRequest.Body)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	for k, v := range httpRequest.Headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", client.UserAgent)

	return req, nil
}
