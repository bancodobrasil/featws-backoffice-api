package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func HomeHandler() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.String(http.StatusOK, "FeatWS Resolver Works!!!")
// 	}

// }

// HomeHandler ...
func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "FeatWS API Works!!!")
}
