package v1

import (
	"strings"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/middlewares"
	"github.com/gin-gonic/gin"
)

// Router define routes the API V1
func Router(router *gin.RouterGroup) {
	cfg := config.GetConfig()

	switch strings.ToLower(cfg.AuthMode) {
	case "openam":
		router.Use(middlewares.VerifyAuthToken())
	}

	rulesheetsRouter(router.Group("/rulesheets"))
	//rpcRouter(router.Group("/"))
}
