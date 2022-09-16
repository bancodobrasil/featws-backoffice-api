package middlewares

import (
	"bytes"
	"context"
	"log"
	"strings"
	"time"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
)

type VerifyAuthTokenMiddleware struct {
	url               string
	ctx               context.Context
	signatureKeyCache *jwk.Cache
}

var verifyAuthTokenMiddleware *VerifyAuthTokenMiddleware

// Middleware function to verify the JWT token
func VerifyAuthToken() gin.HandlerFunc {
	return verifyAuthTokenMiddleware.Run()
}

func NewVerifyAuthTokenMiddleware() {
	cfg := config.GetConfig()
	ctx := context.Background()

	verifyAuthTokenMiddleware = &VerifyAuthTokenMiddleware{
		url:               cfg.OpenAMURL,
		ctx:               ctx,
		signatureKeyCache: jwk.NewCache(ctx, jwk.WithRefreshWindow(1*time.Minute)),
	}

	verifyAuthTokenMiddleware.setup()
}

func (m *VerifyAuthTokenMiddleware) setup() {
	log.Println("Initializing VerifyAuthTokenMiddleware")
	m.signatureKeyCache.Register(m.url, jwk.WithMinRefreshInterval(5*time.Minute))
	_, err := m.signatureKeyCache.Refresh(m.ctx, m.url)
	if err != nil {
		log.Panicf("Failed to refresh OpenAM JWKS: %s\n", err)
	}
}

func (m *VerifyAuthTokenMiddleware) Run() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractTokenFromHeader(c)

		msg, err := jws.Parse([]byte(token))
		errorMsg := "Invalid Auth JWT"
		if err != nil {
			respondWithError(c, 401, errorMsg)
		}

		key := m.getSignatureKey(c)
		verified, err := jws.Verify([]byte(token), jws.WithKey(jwa.RS256, key))
		if err != nil {
			respondWithError(c, 401, errorMsg)
		}

		if !bytes.Equal(verified, msg.Payload()) {
			respondWithError(c, 401, errorMsg)
		}

		c.Next()
	}
}

func (m *VerifyAuthTokenMiddleware) extractTokenFromHeader(c *gin.Context) string {
	authorizationHeader := c.Request.Header.Get("Authorization")
	if authorizationHeader == "" {
		respondWithError(c, 401, "Missing Authorization Header")
	}
	splitHeader := strings.Split(authorizationHeader, "Auth JWT")
	if len(splitHeader) != 2 {
		respondWithError(c, 401, "Invalid Authorization Header")
	}
	return strings.TrimSpace(splitHeader[1])
}
func (m *VerifyAuthTokenMiddleware) getSignatureKey(c *gin.Context) jwk.Key {
	keyset, err := m.signatureKeyCache.Get(m.ctx, m.url)
	errorMsg := "Failed to fetch OpenAM JWKS"
	if err != nil {
		log.Printf("%s: %s\n", errorMsg, err)
		respondWithError(c, 502, errorMsg)
	}
	key, exists := keyset.Key(0)
	if !exists {
		log.Printf("%s: %s\n", errorMsg, err)
		respondWithError(c, 502, errorMsg)
	}
	return key
}
