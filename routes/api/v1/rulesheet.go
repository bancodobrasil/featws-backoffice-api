package v1

import (
	"github.com/bancodobrasil/featws-api/config"
	v1 "github.com/bancodobrasil/featws-api/controllers/v1"
	"github.com/bancodobrasil/featws-api/repository"
	"github.com/bancodobrasil/featws-api/services"
	"github.com/gin-gonic/gin"
)

func rulesheetsRouter(router *gin.RouterGroup) {

	cfg := config.GetConfig()

	repository := repository.GetRulesheets()

	gitlabService := services.NewGitlab(cfg)

	service := services.NewRulesheets(repository, gitlabService)

	controller := v1.NewRulesheets(service)

	router.POST("/", controller.CreateRulesheet())
	router.GET("/", controller.GetRulesheets())
	router.GET("/:id", controller.GetRulesheet())
	router.PUT("/:id", controller.UpdateRulesheet())
	router.DELETE("/:id", controller.DeleteRulesheet())
}
