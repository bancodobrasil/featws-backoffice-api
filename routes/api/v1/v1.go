package v1

import (
	"github.com/bancodobrasil/featws-api/middlewares"
	"github.com/gin-gonic/gin"
)

// Router define routes the API V1
func Router(router *gin.RouterGroup) {
	if middlewares.Authentication != nil {
		router.Use(middlewares.Authentication.Run())
	}

	rulesheetsRouter(router.Group("/rulesheets"))
	//rpcRouter(router.Group("/"))
}
