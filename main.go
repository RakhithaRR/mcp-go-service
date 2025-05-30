package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	mcp "mcp-server/pkg/mcp"
	"mcp-server/pkg/service"
)

var logger = service.GetLogger()

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
	} else if mcpRequest.Arguments == "" {
		logger.Error("Arguments are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arguments are required"})
		return
	} else if mcpRequest.Schema == "" {
		logger.Warn("Input schema is not provided")
	}
	if mcpRequest.IsProxy {
		if mcpRequest.API.APIName == "" {
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
		}
	} else {
		if mcpRequest.Backend.Endpoint == "" {
			logger.Error("Backend endpoint is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Backend endpoint is required"})
			return
		} else if mcpRequest.Backend.Target == "" {
			logger.Error("Backend target is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Backend target is required"})
			return
		} else if mcpRequest.Backend.Verb == "" {
			logger.Error("Backend verb is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Backend verb is required"})
			return
		}
	}

	// Set logging context
	ctx := context.Background()

	ctx = context.WithValue(ctx, service.ToolNameKey, mcpRequest.ToolName)
	if mcpRequest.API.APIName != "" {
		ctx = context.WithValue(ctx, service.ApiNameKey, mcpRequest.API.APIName)
	}

	if mcpRequest.IsProxy && mcpRequest.API.Auth == "" {
		logger.WarnContext(ctx, "Authentication is not provided for the underlying API. Assuming no authentication is required.")
	}

	// Call the underlying API
	logger.Info("Calling underlying API", "tool_name", mcpRequest.ToolName)
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
	cfg, err := service.InitConfig()
	if err != nil {
		logger.Error("Failed to get configurations", "error", err)
		return
	}
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info(fmt.Sprintf("Starting server on %s...", address))
	if cfg.Server.Secure {
		err = router.RunTLS(address, cfg.Server.CertPath, cfg.Server.KeyPath)
	} else {
		logger.Warn("Starting server in insecure mode.")
		err = router.Run(address)
	}
	if err != nil {
		logger.Error("Failed to start the service", "error", err)
		return
	}

}
