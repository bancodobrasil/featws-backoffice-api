package v1

import (
	v1 "github.com/bancodobrasil/featws-api/controllers/v1"
	"github.com/gin-gonic/gin"
)

func rulesheetsRouter(router *gin.RouterGroup) {
	router.POST("/", v1.CreateRulesheet())
	router.GET("/", v1.GetRulesheets())
	router.GET("/:id", v1.GetRulesheet())
	router.PUT("/:id", v1.UpdateRulesheet())
	router.DELETE("/:id", v1.DeleteRulesheet())
}
