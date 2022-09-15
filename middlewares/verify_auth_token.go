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
	OpenAMURL         string
	ctx               context.Context
	signatureKeyCache *jwk.Cache
}

var verifyAuthTokenMiddleware *VerifyAuthTokenMiddleware

func NewVerifyAuthTokenMiddleware() {
	cfg := config.GetConfig()
	ctx := context.Background()

	verifyAuthTokenMiddleware = &VerifyAuthTokenMiddleware{
		OpenAMURL:         cfg.OpenAMURL,
		ctx:               ctx,
		signatureKeyCache: jwk.NewCache(ctx, jwk.WithRefreshWindow(5*time.Minute)),
	}

	verifyAuthTokenMiddleware.signatureKeyCache.Register(cfg.OpenAMURL, jwk.WithMinRefreshInterval(15*time.Minute))
	_, err := verifyAuthTokenMiddleware.signatureKeyCache.Refresh(ctx, cfg.OpenAMURL)
	if err != nil {
		log.Panicf("Failed to refresh OpenAM JWKS: %s\n", err)
	}
}

func RunVerifyAuthTokenMiddleware() gin.HandlerFunc {
	return verifyAuthTokenMiddleware.Run()
}

func (m *VerifyAuthTokenMiddleware) Run() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)

		msg, err := jws.Parse([]byte(token))
		errorMsg := "Invalid Auth JWT"
		if err != nil {
			respondWithError(c, 401, errorMsg)
		}

		key := getSignatureKey(c)
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

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func extractTokenFromHeader(c *gin.Context) string {
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
func getSignatureKey(c *gin.Context) jwk.Key {
	keyset, err := verifyAuthTokenMiddleware.signatureKeyCache.Get(verifyAuthTokenMiddleware.ctx, verifyAuthTokenMiddleware.OpenAMURL)
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
