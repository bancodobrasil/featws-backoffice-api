package middlewares

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
)

var cfg = config.GetConfig()
var ctx = context.Background()
var signatureKeyCache = jwk.NewCache(ctx)

func InitializeSignatureKeyCache() {
	signatureKeyCache.Register(cfg.OpenAMURL, jwk.WithMinRefreshInterval(15*time.Minute))
	_, err := signatureKeyCache.Refresh(ctx, cfg.OpenAMURL)
	if err != nil {
		fmt.Printf("Failed to refresh OpenAM JWKS: %s\n", err)
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
	keyset, err := signatureKeyCache.Get(ctx, cfg.OpenAMURL)
	errorMsg := "Failed to fetch OpenAM JWKS"
	if err != nil {
		fmt.Printf("%s: %s\n", errorMsg, err)
		respondWithError(c, 502, errorMsg)
	}
	key, exists := keyset.Key(0)
	if !exists {
		fmt.Printf("%s: %s\n", errorMsg, err)
		respondWithError(c, 502, errorMsg)
	}
	return key
}
func VerifyAuthTokenMiddleware() gin.HandlerFunc {
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
