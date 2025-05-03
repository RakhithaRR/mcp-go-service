package mcp

import (
	"context"
	"io"
	"mcp-server/pkg/service"
	"net/http"
	"strings"
)

var logger = service.GetLogger()

func CallUnderlyingAPI(ctx context.Context, payload *MCPRequest) (string, int, error) {
	httpClient := InitHttpClient()
	httpRequest, err := transformMCPRequest(payload)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to transform request", "error", err)
		return "", http.StatusInternalServerError, err
	}
	request, err := httpClient.GenerateRequest(httpRequest)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to generate request", "error", err)
		return "", http.StatusInternalServerError, err
	}
	resp, err := httpClient.DoRequest(request)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to send request", "error", err)
		return "", http.StatusInternalServerError, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to read response body", "error", err)
		return "", http.StatusInternalServerError, err
	}
	response := string(body)
	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		response, err = processJsonResponse(response)
		if err != nil {
			logger.ErrorContext(ctx, "Failed to process JSON response", "error", err)
			return "", http.StatusInternalServerError, err
		}
	}
	return response, resp.StatusCode, nil
}
