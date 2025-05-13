package mcp

import (
	"encoding/json"
	"slices"
	"strings"
)

func processSchema(schema string) (*SchemaMapping, error) {
	var inputSchema MCPInputSchema
	err := json.Unmarshal([]byte(schema), &inputSchema)
	if err != nil {
		logger.Error("Error processing the MCP input schema", "error", err)
		return nil, err
	}

	schemaMapping := processInputProperties(inputSchema.Properties, inputSchema.Required)
	if inputSchema.ContentType != "" {
		schemaMapping.ContentType = inputSchema.ContentType
	} else {
		schemaMapping.ContentType = ContentTypeJSON
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
