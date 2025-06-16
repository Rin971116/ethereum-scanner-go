package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloWorldHandler handles the hello world request
func HelloWorldHandler(c *gin.Context) {
	// Respond with a simple message
	c.String(http.StatusOK, "Hello, World!")
}
