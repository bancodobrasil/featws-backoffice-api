package auth

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
)

// VerifyAuthTokenMiddleware stores the OpenAM URL to be used for authentication
// and the signature key cache
type VerifyAuthTokenMiddleware struct {
	url               string
	ctx               context.Context
	signatureKeyCache *jwk.Cache
}

// NewVerifyAuthTokenMiddleware returns a new VerifyAuthTokenMiddleware instance
func NewVerifyAuthTokenMiddleware() *VerifyAuthTokenMiddleware {
	cfg := config.GetConfig()
	ctx := context.Background()

	verifyAuthTokenMiddleware := &VerifyAuthTokenMiddleware{
		url:               cfg.OpenAMURL,
		ctx:               ctx,
		signatureKeyCache: jwk.NewCache(ctx, jwk.WithRefreshWindow(1*time.Minute)),
	}

	verifyAuthTokenMiddleware.setup()

	return verifyAuthTokenMiddleware
}

// setup sets up the signature key cache
func (m *VerifyAuthTokenMiddleware) setup() {
	log.Println("Initializing VerifyAuthTokenMiddleware")
	m.signatureKeyCache.Register(m.url, jwk.WithMinRefreshInterval(5*time.Minute))
	_, err := m.signatureKeyCache.Refresh(m.ctx, m.url)
	if err != nil {
		log.Panicf("Failed to refresh OpenAM JWKS: %s\n", err)
	}
}

// Authenticate runs the authentication middleware
func (m *VerifyAuthTokenMiddleware) Authenticate(h *http.Header) (err error, statusCode int) {
	token, err, statusCode := m.extractTokenFromHeader(h)
	if err != nil {
		return err, statusCode
	}

	invalidJWTError := errors.New("Invalid JWT token")
	defaultStatusCode := 401

	msg, internalErr := jws.Parse([]byte(token))
	if internalErr != nil {
		return invalidJWTError, defaultStatusCode
	}

	key, err, statusCode := m.getSignatureKey()
	if err != nil {
		return err, statusCode
	}

	verified, internalErr := jws.Verify([]byte(token), jws.WithKey(jwa.RS256, key))
	if internalErr != nil {
		return invalidJWTError, defaultStatusCode
	}

	if !bytes.Equal(verified, msg.Payload()) {
		return invalidJWTError, defaultStatusCode
	}

	return nil, 0
}

func (m *VerifyAuthTokenMiddleware) extractTokenFromHeader(h *http.Header) (string, error, int) {
	authorizationHeader := h.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("Missing Authorization Header"), 401
	}
	splitHeader := strings.Split(authorizationHeader, "Bearer")
	if len(splitHeader) != 2 {
		return "", errors.New("Invalid Authorization Header"), 401
	}
	return strings.TrimSpace(splitHeader[1]), nil, 0
}
func (m *VerifyAuthTokenMiddleware) getSignatureKey() (jwk.Key, error, int) {
	keyset, err := m.signatureKeyCache.Get(m.ctx, m.url)
	errorMsg := "Failed to fetch OpenAM JWKS"
	if err != nil {
		log.Printf("%s: %s\n", errorMsg, err)
		return nil, errors.New(errorMsg), 502
	}
	key, exists := keyset.Key(0)
	if !exists {
		log.Printf("%s: %s\n", errorMsg, err)
		return nil, errors.New(errorMsg), 502
	}
	return key, nil, 0
}
