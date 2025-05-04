package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	mcp "mcp-server/pkg/mcp"
	"mcp-server/pkg/service"
)

var logger = service.GetLogger()

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

func serveRequest(c *gin.Context) {
	var mcpRequest mcp.MCPRequest

	if err := c.BindJSON(&mcpRequest); err != nil {
		logger.Error("Failed to bind JSON", "error", err)
		return
	}
	// Validate the request
	if mcpRequest.ToolName == "" {
		logger.Error("Tool name is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tool name is required"})
		return
	} else if mcpRequest.API.APIName == "" {
		logger.Warn("API name is not proided")
	} else if mcpRequest.API.Endpoint == "" {
		logger.Error("API endpoint is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API endpoint is required"})
		return
	} else if mcpRequest.API.Context == "" {
		logger.Error("API context is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API context is required"})
		return
	} else if mcpRequest.API.Version == "" {
		logger.Error("API version is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API version is required"})
		return
	} else if mcpRequest.API.Path == "" {
		logger.Error("Resource path is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource path is required"})
		return
	} else if mcpRequest.API.Verb == "" {
		logger.Error("HTTP verb is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "HTTP verb is required"})
		return
	} else if mcpRequest.Arguments == "" {
		logger.Error("Arguments are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arguments are required"})
		return
	} else if mcpRequest.Schema == "" {
		logger.Warn("Input schema is not provided")
	}

	// Set logging context
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKey("toolName").String(), mcpRequest.ToolName)
	if mcpRequest.API.APIName != "" {
		ctx = context.WithValue(ctx, contextKey("apiName").String(), mcpRequest.API.APIName)
	}

	if mcpRequest.API.Auth == "" {
		logger.WarnContext(ctx, "Authentication is not provided for the underlying API. Assuming no authentication is required.")
	}

	// Call the underlying API
	resp, code, err := mcp.CallUnderlyingAPI(ctx, &mcpRequest)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to call underlying API", "error", err)
		c.JSON(code, gin.H{"error": "Failed to call underlying API", "details": err.Error()})
		return
	}
	c.SecureJSON(code, resp)
}

func main() {
	router := service.GetRouter()
	router.POST("/mcp", serveRequest)
	logger.Info("Service started on localhost:8080...")
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		logger.Error("Failed to start the service", "error", err)
	}

}
