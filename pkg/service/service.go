package service

import (
	"github.com/gin-gonic/gin"
)

type Health struct {
	Status string `json:"status"`
}

func getHealth(c *gin.Context) {
	health := Health{Status: "OK"}
	c.IndentedJSON(200, health)
}

func GetRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/health", getHealth)
	return router
}
