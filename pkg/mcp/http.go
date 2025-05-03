package mcp

import (
	"fmt"
	"net/http"
)

type ToolHTTPClient struct {
	httpClient *http.Client
}

func InitHttpClient() *ToolHTTPClient {
	client := &http.Client{}
	return &ToolHTTPClient{
		httpClient: client,
	}
}

func (client *ToolHTTPClient) DoRequest(request *http.Request) (*http.Response, error) {
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GenerateRequest(httpRequest *TransformedRequest) (*http.Request, error) {
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

	return req, nil
}
