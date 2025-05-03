package mcp

import "bytes"

type MCPRequest struct {
	ToolName  string  `json:"tool_name"`
	API       APIInfo `json:"api"`
	Arguments string  `json:"arguments,omitempty"`
	Schema    string  `json:"schema,omitempty"`
}

type APIInfo struct {
	APIName  string `json:"api_name"`
	Endpoint string `json:"endpoint"`
	Context  string `json:"context"`
	Version  string `json:"version"`
	Path     string `json:"path"`
	Verb     string `json:"verb"`
	Auth     string `json:"auth,omitempty"`
}

type TransformedRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    *bytes.Reader
}

type SchemaMapping struct {
	PathParameters   []string `json:"pathParameters"`
	QueryParameters  []string `json:"queryParameters"`
	HeaderParameters []string `json:"headerParameters"`
	HasBody          bool     `json:"hasBody"`
	ContentType      string   `json:"contentType,omitempty"`
}

type MCPSchema struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema *MCPInputSchema `json:"inputSchema"`
	Required    []string        `json:"required"`
}

type MCPInputSchema struct {
	Type           string         `json:"type"`
	Properties     map[string]any `json:"properties"`
	RequiredFields []string       `json:"requiredFields,omitempty"`
	ContentType    string         `json:"contentType,omitempty"`
}
