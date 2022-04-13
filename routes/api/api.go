package api

import (
	v1 "github.com/bancodobrasil/featws-api/routes/api/v1"
	"github.com/gin-gonic/gin"
)

// Router define routes the API
func Router(router *gin.RouterGroup) {
	v1.Router(router.Group("/v1"))
}
