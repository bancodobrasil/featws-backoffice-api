package routes

import (
	"github.com/bancodobrasil/featws-api/routes/api"
	telemetry "github.com/bancodobrasil/gin-telemetry"
	"github.com/gin-gonic/gin"
)

// APIRoutes define all api routes
func APIRoutes(router *gin.Engine) {
	// inject middleware
	router.Use(telemetry.Middleware("featws-api"))
	api.Router(router.Group("/api"))
}
