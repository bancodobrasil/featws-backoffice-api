package routes

import (
	"github.com/bancodobrasil/featws-api/routes/api"
	"github.com/gin-gonic/gin"
)

// SetupRoutes define all routes
func SetupRoutes(router *gin.Engine) {

	homeRouter(router.Group("/"))
	api.Router(router.Group("/api"))
}
