package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler ...
func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "FeatWS API Works!!!",
	})
}
