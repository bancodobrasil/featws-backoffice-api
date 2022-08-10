package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/bancodobrasil/featws-api/routes/api"
	telemetry "github.com/bancodobrasil/gin-telemetry"
)

// ApiRoutes define all api routes
func ApiRoutes(router *gin.Engine) {
	// inject middleware
	router.Use(telemetry.Middleware("featws-api"))
	api.Router(router.Group("/api"))
}
