package health

import (
	"github.com/bancodobrasil/featws-api/controllers"
	"github.com/gin-gonic/gin"
)

// Router ...
func Router(router *gin.RouterGroup) {
	router.GET("/live", controllers.HealthLiveHandler())
}
