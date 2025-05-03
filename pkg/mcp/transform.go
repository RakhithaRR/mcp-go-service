package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func transformMCPRequest(mcpRequest *MCPRequest) (*TransformedRequest, error) {
	httpRequest := &TransformedRequest{
		Headers: make(map[string]string),
	}

	method, err := processHTTPMethod(mcpRequest.Verb)
	if err != nil {
		logger.Error("Failed to process HTTP method", "error", err)
		return nil, err
	}
	httpRequest.Method = method

	schemaMapping, err := processSchema(mcpRequest.Schema)
	if err != nil {
		logger.Error("Failed to process schema", "error", err)
		return nil, err
	}

	ep, err := processEndpoint(mcpRequest, schemaMapping)
	if err != nil {
		logger.Error("Failed to process endpoint", "error", err)
		return nil, err
	}
	httpRequest.URL = ep

	headers, err := processHeaderParameters(mcpRequest, schemaMapping)
	if err != nil {
		logger.Error("Failed to process header parameters", "error", err)
		return nil, err
	}
	httpRequest.Headers = headers

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		bytesReader, err := processRequestBody(mcpRequest, schemaMapping)
		if err != nil {
			logger.Error("Failed to process request body", "error", err)
			return nil, err
		}
		if bytesReader != nil {
			httpRequest.Body = bytesReader
		} else {
			logger.Warn("Request body is nil")
		}
	}

	return httpRequest, nil
}

func processHTTPMethod(verb string) (string, error) {
	method := strings.ToUpper(verb)
	switch method {
	case "GET":
		return http.MethodGet, nil
	case "POST":
		return http.MethodPost, nil
	case "PUT":
		return http.MethodPut, nil
	case "DELETE":
		return http.MethodDelete, nil
	case "PATCH":
		return http.MethodPatch, nil
	case "OPTIONS":
		return http.MethodOptions, nil
	default:
		return "", fmt.Errorf("unsupported HTTP method: %s", verb)
	}
}

// processEndpoint constructs the endpoint URL for the request.
// It processes path parameters and query parameters based on the schema mapping.
// Returns the transformed endpoint URL or an error if the endpoint is invalid.
func processEndpoint(mcpRequest *MCPRequest, schemaMapping *SchemaMapping) (string, error) {
	args, err := parseArgs(mcpRequest)
	if err != nil {
		logger.Error("Failed to parse arguments", "error", err)
		return "", err
	}

	endpoint := mcpRequest.Endpoint
	transformedEp := ""
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		transformedEp = fmt.Sprintf("%s/%s/%s/%s", endpoint, mcpRequest.Context, mcpRequest.Version, mcpRequest.Path)
		// Process path parameters
		transformedEp = processPathParameters(args, schemaMapping, transformedEp)
		// Process query parameters
		queryParams := processQueryParameters(args, schemaMapping)
		if queryParams != "" {
			transformedEp += queryParams
		}
		return transformedEp, nil
	}
	return "", fmt.Errorf("invalid endpoint: %s", endpoint)
}

// processQueryParameters generates a query string from the provided arguments and schema mapping.
// It URL-encodes parameter names and values and appends them to the query string.
// Returns the constructed query string.
func processQueryParameters(args map[string]any, schemaMapping *SchemaMapping) string {
	queryParams := schemaMapping.QueryParameters
	if len(queryParams) > 0 {
		queryString := "?"
		for _, param := range queryParams {
			paramValue := args[param]
			if paramValue == nil {
				logger.Warn("Query parameter value is not available", "parameter", param)
				continue
			}
			// URL encode the parameter name and value
			urlEncodedParam := url.QueryEscape(fmt.Sprintf("%v", param))
			urlEncodedValue := url.QueryEscape(fmt.Sprintf("%v", paramValue))
			queryString += fmt.Sprintf("%s=%s&", urlEncodedParam, urlEncodedValue)
		}
		queryString = strings.TrimSuffix(queryString, "&")
		return queryString
	}
	return ""
}

// processPathParameters replaces placeholders in the URL with actual values from the arguments.
// It URL-encodes parameter names and values before substitution.
// Returns the transformed URL with path parameters replaced.
func processPathParameters(args map[string]any, schemaMapping *SchemaMapping, unProcessedUrl string) string {
	pathParams := schemaMapping.PathParameters
	transformedUrl := unProcessedUrl
	if len(pathParams) > 0 {
		for _, param := range pathParams {
			paramValue := args[param]
			if paramValue == nil {
				logger.Warn("Path parameter value is not available", "parameter", param)
				continue
			}
			// URL encode the parameter name and value
			urlEncodedParam := url.QueryEscape(fmt.Sprintf("%v", param))
			urlEncodedValue := url.QueryEscape(fmt.Sprintf("%v", paramValue))
			processedUrl := strings.Replace(transformedUrl, "{"+urlEncodedParam+"}", urlEncodedValue, 1)
			transformedUrl = processedUrl
		}
	}
	return transformedUrl
}

// processHeaderParameters generates a map of header parameters from the provided arguments and schema mapping.
// Returns a map of header names and values.
func processHeaderParameters(mcpRequest *MCPRequest, schemaMapping *SchemaMapping) (map[string]string, error) {
	args, err := parseArgs(mcpRequest)
	if err != nil {
		logger.Error("Failed to parse arguments", "error", err)
		return nil, err
	}
	headers := make(map[string]string)
	headerParams := schemaMapping.HeaderParameters
	if len(headerParams) > 0 {
		for _, param := range headerParams {
			paramValue := args[param]
			if paramValue == nil {
				logger.Warn("Header parameter value is not available", "parameter", param)
				continue
			}
			headers[param] = fmt.Sprintf("%v", paramValue)
		}
	}
	return headers, nil
}

// processRequestBody processes the request body from the MCP request.
// Returns a bytes.Reader for the request body or an error if the body is invalid.
func processRequestBody(mcpRequest *MCPRequest, schemaMapping *SchemaMapping) (*bytes.Reader, error) {
	contentType := schemaMapping.ContentType
	args, err := parseArgs(mcpRequest)
	if err != nil {
		logger.Error("Failed to parse arguments", "error", err)
		return nil, err
	}

	body := args["requestBody"].(string)
	if body == "" {
		logger.Warn("Request body is empty")
		return nil, nil
	}
	if strings.HasPrefix(body, "{") && strings.HasSuffix(body, "}") {
		if contentType == "application/json" {
			byteArray := []byte(body)
			bodyReader := bytes.NewReader(byteArray)
			return bodyReader, nil
		} else {
			logger.Error("Unsupported content type for request body", "contentType", contentType)
			return nil, fmt.Errorf("unsupported content type for request body: %s", contentType)
		}
	} else {
		logger.Error("Invalid request body format")
		return nil, fmt.Errorf("invalid request body format")
	}
}

func parseArgs(mcpRequest *MCPRequest) (map[string]any, error) {
	var args map[string]any
	if mcpRequest.Arguments != "" {
		err := json.Unmarshal([]byte(mcpRequest.Arguments), &args)
		if err != nil {
			logger.Error("Failed to unmarshal arguments", "error", err)
			return nil, err
		}
	}
	return args, nil
}
