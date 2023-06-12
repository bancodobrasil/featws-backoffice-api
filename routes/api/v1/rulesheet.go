package v1

import (
	"github.com/bancodobrasil/featws-api/config"
	v1 "github.com/bancodobrasil/featws-api/controllers/v1"
	"github.com/bancodobrasil/featws-api/repository"
	"github.com/bancodobrasil/featws-api/services"
	"github.com/gin-gonic/gin"
)

// rulesheetsRouter sets up the routing for CRUD operations on rulesheets using Gin framework
func rulesheetsRouter(router *gin.RouterGroup) {

	cfg := config.GetConfig()

	// Repository enables the code to access the data repository and perform CRUD operations on the rulesheets.
	repository := repository.GetRulesheets()

	// The project connects to GitLab to access the rulesheets and perform CRUD operations
	gitlabService := services.NewGitlab(cfg)

	// The service variable is creating a new instance of the Rulesheets service from the services package,
	// it takes the parameters: repository and gitlabService. These parameters allow the service to access
	// the data repository and the GitLab service
	service := services.NewRulesheets(repository, gitlabService)

	// The controller is creating a new instance of the "Rulesheets" controller from the "v1"
	// package and passing an instance of the service as a parameter. This allows the controller
	// to have access to the business logic and functionalities provided by the service.
	controller := v1.NewRulesheets(service)

	// These are the API endpoints
	router.POST("/", controller.CreateRulesheet())
	router.GET("/", controller.GetRulesheets())
	router.GET("/:id", controller.GetRulesheet())
	router.PUT("/:id", controller.UpdateRulesheet())
	router.DELETE("/:id", controller.DeleteRulesheet())
}
