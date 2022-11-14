package v1

import (
	goauthgin "github.com/bancodobrasil/goauth-gin"
	"github.com/gin-gonic/gin"
)

// Router define routes the API V1
func Router(router *gin.RouterGroup) {

	router.Use(goauthgin.Authenticate())
	rulesheetsRouter(router.Group("/rulesheets"))
	//rpcRouter(router.Group("/"))
}
