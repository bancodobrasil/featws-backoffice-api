package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/healthcheck"
	"github.com/bancodobrasil/healthcheck/checks/db"
	"github.com/bancodobrasil/healthcheck/checks/goroutine"
	"github.com/gin-gonic/gin"
	"github.com/gsdenys/healthcheck/checks"
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
	cfg := config.GetConfig()
	rawGitlabURL := cfg.GitlabURL
	gitlab, _ := url.Parse(rawGitlabURL)
	gitlabUrl := gitlab.Scheme + "://" + gitlab.Host
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", goroutine.Count(100))
	database, err := database.GetConn().DB()
	if err != nil {
		log.Fatal(err)
	}

	health.AddReadinessCheck("database", db.Ping(database, 1*time.Second))
	health.AddReadinessCheck("gitlab", Get(gitlabUrl, 1*time.Second))
	return health
}

// Get was the function that allow follow the url
func Get(url string, timeout time.Duration) checks.Check {
	client := http.Client{
		Timeout: timeout,
	}

	return func() error {
		resp, err := client.Get(url)

		if err != nil {
			return err
		}

		resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("returned status %d", resp.StatusCode)
		}

		return nil
	}
}

// HealthLiveHandler ...
func (c *HealthController) HealthLiveHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.LiveEndpoint))
}

// HealthReadyHandler ...
func (c *HealthController) HealthReadyHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.ReadyEndpoint))
}
