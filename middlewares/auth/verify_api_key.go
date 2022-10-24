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
func (m *VerifyAPIKeyMiddleware) Authenticate(h *http.Header) (err error, statusCode int) {
	key, err, statusCode := m.extractKeyFromHeader(h)
	if err != nil {
		return err, statusCode
	}

	if key != m.key {
		return errors.New("Unauthorized"), 401
	}

	return nil, 0
}

func (m *VerifyAPIKeyMiddleware) extractKeyFromHeader(h *http.Header) (key string, err error, statusCode int) {
	authorizationHeader := h.Get("X-API-Key")
	if authorizationHeader == "" {
		return "", errors.New("Missing X-API-Key Header"), 401
	}
	return authorizationHeader, nil, 0
}
