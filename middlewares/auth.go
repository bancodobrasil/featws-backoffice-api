package middlewares

import (
	"net/http"
	"strings"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/middlewares/auth"
	"github.com/gin-gonic/gin"
)

// AuthenticationMiddleware is the interface that wraps the AuthenticateFunc method
// and is used to authenticate the request
type AuthenticationMiddleware interface {
	Authenticate(h *http.Header) (err error, statusCode int)
}

// AuthMiddleware stores the authentication middlewares to be run
type AuthMiddleware struct {
	mode        []string
	middlewares []AuthenticationMiddleware
}

// NewAuthMiddleware returns a new AuthMiddleware instance
func NewAuthMiddleware() *AuthMiddleware {
	cfg := config.GetConfig()

	mode := strings.Split(strings.ToLower(cfg.AuthMode), ",")

	if len(mode) == 0 || mode[0] == "" || mode[0] == "none" {
		return nil
	}

	middlewares := []AuthenticationMiddleware{}

	for _, m := range mode {
		switch m {
		case "api_key":
			middlewares = append(middlewares, auth.NewVerifyAPIKeyMiddleware())
		case "openam":
			middlewares = append(middlewares, auth.NewVerifyAuthTokenMiddleware())
		}
	}

	return &AuthMiddleware{
		mode:        mode,
		middlewares: middlewares,
	}
}

// Run executes all the authentication middlewares in the order they were added.
// If any of the middlewares does not return an error, the request proceeds to the next handler.
// If the last middleware returns an error, the request is aborted.
func (m *AuthMiddleware) Run() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var statusCode int

		for _, middleware := range m.middlewares {
			err, statusCode = middleware.Authenticate(&c.Request.Header)
			if err == nil {
				c.Next()
				return
			}
		}

		if err != nil {
			respondWithError(c, &MiddlewareError{
				Code:    statusCode,
				Message: err.Error(),
			})
		}
	}
}
