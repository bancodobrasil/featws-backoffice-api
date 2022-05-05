package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler ...
func HealthLiveHandler(c *gin.Context) {
	c.String(http.StatusOK, "Application is live!!!")
}
