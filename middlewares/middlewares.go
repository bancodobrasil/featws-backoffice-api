package middlewares

import "github.com/gin-gonic/gin"

type Middleware interface {
	Run()
}

func InitializeMiddlewares() {
	NewVerifyAuthTokenMiddleware()
}

// Helper function to abort the request with an error status code and message
func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}
