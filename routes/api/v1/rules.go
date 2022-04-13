package v1

import (
	v1 "github.com/bancodobrasil/featws-api/controllers/v1"
	"github.com/gin-gonic/gin"
)

func rulesRouter(router *gin.RouterGroup) {
	router.POST("/", v1.CreateRule())
	router.GET("/", v1.GetRules())
	router.GET("/:id", v1.GetRule())
	router.PUT("/:id", v1.UpdateRule())
	router.DELETE("/:id", v1.DeleteRule())
}
