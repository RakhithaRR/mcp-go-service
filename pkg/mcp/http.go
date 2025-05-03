package mcp

import (
	"fmt"
	"net/http"
	"sync"
)

var syncOnce sync.Once

var httpClient *MCPHTTPClient

type MCPHTTPClient struct {
	httpClient *http.Client
	UserAgent  string
}

func InitHttpClient() *MCPHTTPClient {
	syncOnce.Do(func() {
		if httpClient == nil {
			client := http.Client{
				Transport: &http.Transport{
					MaxIdleConns:    20,
					IdleConnTimeout: 45,
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
