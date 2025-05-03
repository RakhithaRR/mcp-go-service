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
	} else if mcpRequest.APIName == "" {
		logger.Warn("API name is not proided")
	} else if mcpRequest.Endpoint == "" {
		logger.Error("API endpoint is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API endpoint is required"})
		return
	} else if mcpRequest.Context == "" {
		logger.Error("API context is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API context is required"})
		return
	} else if mcpRequest.Version == "" {
		logger.Error("API version is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API version is required"})
		return
	} else if mcpRequest.Path == "" {
		logger.Error("Resource path is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource path is required"})
		return
	} else if mcpRequest.Verb == "" {
		logger.Error("HTTP verb is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "HTTP verb is required"})
		return
	} else if mcpRequest.Arguments == "" {
		logger.Error("Arguments are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arguments are required"})
		return
	}

	// Set logging context
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKey("toolName"), mcpRequest.ToolName)
	if mcpRequest.APIName != "" {
		ctx = context.WithValue(ctx, contextKey("apiName"), mcpRequest.APIName)
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
	err := router.Run("localhost:8080")
	if err != nil {
		logger.Error("Failed to start the service", "error", err)
	}

}
