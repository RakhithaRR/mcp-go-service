package mcp

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

func processSchema(schema string) (*SchemaMapping, error) {
	var mcpSchema MCPSchema
	var inputSchema MCPInputSchema
	err := json.Unmarshal([]byte(schema), &mcpSchema)
	if err != nil {
		logger.Error("Error processing the MCP input schema", "error", err)
		return nil, err
	}
	if mcpSchema.InputSchema != nil {
		inputSchema = *mcpSchema.InputSchema
	} else {
		logger.Error("Input schema is nil")
		return nil, fmt.Errorf("input schema is nil")
	}

	schemaMapping := processInputProperties(inputSchema.Properties, mcpSchema.Required)
	if inputSchema.ContentType != "" {
		schemaMapping.ContentType = inputSchema.ContentType
	} else {
		schemaMapping.ContentType = "application/json"
	}

	return schemaMapping, nil

}

func processInputProperties(properties map[string]any, requiredParams []string) *SchemaMapping {
	var pathParameters []string
	var queryParameters []Param
	var headerParameters []Param
	var requestBody bool
	for k := range properties {
		name := k

		if strings.HasPrefix(name, "query_") {
			refinedName := strings.TrimPrefix(name, "query_")
			param := Param{
				Name:     refinedName,
				Required: false,
			}
			if slices.Contains(requiredParams, name) {
				param.Required = true
			}
			queryParameters = append(queryParameters, param)
		} else if strings.HasPrefix(name, "header_") {
			refinedName := strings.TrimPrefix(name, "header_")
			param := Param{
				Name:     refinedName,
				Required: false,
			}
			if slices.Contains(requiredParams, name) {
				param.Required = true
			}
			headerParameters = append(headerParameters, param)
		} else if strings.HasPrefix(name, "path_") {
			refinedName := strings.TrimPrefix(name, "path_")
			pathParameters = append(pathParameters, refinedName)
		} else if name == "requestBody" {
			requestBody = true
		} else {
			logger.Warn("Unknown property prefix", "name", name)
		}

	}

	schemaMapping := &SchemaMapping{
		PathParameters:   pathParameters,
		QueryParameters:  queryParameters,
		HeaderParameters: headerParameters,
		HasBody:          requestBody,
	}
	return schemaMapping
}
