package middlewares

import "github.com/gin-gonic/gin"

// Middleware is the interface that wraps the Run method
// which is used to run the middleware
type Middleware interface {
	Run() gin.HandlerFunc
}

// Authentication Middleware
var Authentication Middleware

// InitializeMiddlewares initializes the middlewares
func InitializeMiddlewares() {
	Authentication = NewAuthMiddleware()
}

// MiddlewareError is the error type returned by the middlewares
type MiddlewareError struct {
	Code    int
	Message string
}

// Error implements the error interface
func (e *MiddlewareError) Error() string {
	return e.Message
}

// Helper function to abort the request with an error status code and message
func respondWithError(c *gin.Context, e *MiddlewareError) {
	c.AbortWithStatusJSON(e.Code, gin.H{"error": e.Message})
}
