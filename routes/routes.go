package routes

import (
	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/docs"
	"github.com/bancodobrasil/featws-api/routes/api"
	"github.com/bancodobrasil/featws-api/routes/health"
	telemetry "github.com/bancodobrasil/gin-telemetry"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// SetupRoutes define all routes
func SetupRoutes(router *gin.Engine) {

	cfg := config.GetConfig()

	docs.SwaggerInfo.Host = cfg.ExternalHost

	homeRouter(router.Group("/"))
	health.Router(router.Group("/health"))
	// setup swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// APIRoutes define all api routes
func APIRoutes(router *gin.Engine) {
	// inject middleware
	router.Use(telemetry.Middleware("featws-api"))
	api.Router(router.Group("/api"))
}
