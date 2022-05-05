package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthLiveHandler ...
func HealthLiveHandler(c *gin.Context) {
	c.String(http.StatusOK, "Application is live!!!")
}
