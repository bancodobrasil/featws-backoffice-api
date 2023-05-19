package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler function returns a JSON response with a message indicating that the FeatWS API is working.
func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "FeatWS API Works!!!",
	})
}
