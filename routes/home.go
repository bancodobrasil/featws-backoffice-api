package routes

import (
	"github.com/bancodobrasil/featws-api/controllers"
	"github.com/gin-gonic/gin"
)

func homeRouter(router *gin.RouterGroup) {
	router.GET("/", controllers.HomeHandler)
}
