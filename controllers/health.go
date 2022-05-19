package controllers

import (
	"net/http"
	"time"

	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/healthcheck"
	"github.com/bancodobrasil/healthcheck/checks/db"
	"github.com/bancodobrasil/healthcheck/checks/goroutine"
	healthHttp "github.com/bancodobrasil/healthcheck/checks/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HealthController ...
type HealthController struct {
	health healthcheck.Handler
}

// NewHealthController ...
func NewHealthController() *HealthController {
	return &HealthController{
		health: newHandler(),
	}
}

func newHandler() healthcheck.Handler {
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", goroutine.Count(100))
	database, err := database.GetConn().DB()
	if err != nil {
		log.Fatal(err)
	}

	health.AddReadinessCheck("database", db.Ping(database, 1*time.Second))
	health.AddReadinessCheck("gitlab", healthHttp.Get("https://about.gitlab.com", 1*time.Second))
	return health
}

// HealthLiveHandler ...
func (c *HealthController) HealthLiveHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.LiveEndpoint))
}

// HealthReadyHandler ...
func (c *HealthController) HealthReadyHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.ReadyEndpoint))
}
