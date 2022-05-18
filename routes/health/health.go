package health

import (
	"github.com/bancodobrasil/featws-api/controllers"
	"github.com/gin-gonic/gin"
)

// Router ...
func Router(router *gin.RouterGroup) {

	healthController := controllers.NewHealthController()
	router.GET("/live", healthController.HealthLiveHandler())
	router.GET("/ready", healthController.HealthReadyHandler())
}
