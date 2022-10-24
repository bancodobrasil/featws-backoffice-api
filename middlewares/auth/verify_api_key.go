package auth

import (
	"errors"
	"net/http"

	"github.com/bancodobrasil/featws-api/config"
)

// VerifyAPIKeyMiddleware stores the API key to be used for authentication
type VerifyAPIKeyMiddleware struct {
	key string
}

// NewVerifyAPIKeyMiddleware returns a new VerifyAPIKeyMiddleware instance
func NewVerifyAPIKeyMiddleware() *VerifyAPIKeyMiddleware {
	cfg := config.GetConfig()

	return &VerifyAPIKeyMiddleware{
		key: cfg.AuthAPIKey,
	}

}

// Authenticate runs the authentication middleware
func (m *VerifyAPIKeyMiddleware) Authenticate(h *http.Header) (statusCode int, err error) {
	key, statusCode, err := m.extractKeyFromHeader(h)
	if err != nil {
		return statusCode, err
	}

	if key != m.key {
		return 401, errors.New("Unauthorized")
	}

	return 0, nil
}

func (m *VerifyAPIKeyMiddleware) extractKeyFromHeader(h *http.Header) (key string, statusCode int, err error) {
	authorizationHeader := h.Get("X-API-Key")
	if authorizationHeader == "" {
		return "", 401, errors.New("Missing X-API-Key Header")
	}
	return authorizationHeader, 0, nil
}
