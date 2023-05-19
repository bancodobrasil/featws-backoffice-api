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

// HealthController the health endpoints controller
type HealthController struct {
	health healthcheck.Handler
}

// NewHealthController returns a new instance of the HealthController struct with a newHandler.
func NewHealthController() *HealthController {
	return &HealthController{
		health: newHandler(),
	}
}

var health = healthcheck.NewHandler()

// newHandler creates a new healthcheck handler and adds liveness and readiness checks for goroutine
// count, Gitlab URL, and database connection.
func newHandler() healthcheck.Handler {
	cfg := config.GetConfig()
	health.AddLivenessCheck("goroutine-threshold", goroutine.Count(100))

	if cfg.GitlabURL != "" {
		rawGitlabURL := cfg.GitlabURL
		gitlab, _ := url.Parse(rawGitlabURL)

		if gitlab.Scheme == "" {
			log.Fatal("gitlab must have a scheme: http:// or https://")
		}

		if gitlab.Host == "" {
			log.Fatal("gitlab must have a host: example.com")
		}
		gitlabURL := gitlab.Scheme + "://" + gitlab.Host
		health.AddReadinessCheck("gitlab", Get(gitlabURL, 1*time.Second))
	}

	database, err := database.GetConn().DB()
	if err != nil {
		log.Fatal(err)
	}

	health.AddReadinessCheck("database", db.Ping(database, 1*time.Second))
	return health
}

// Get returns a check function that performs an HTTP GET request to a specified URL with a
// specified timeout and returns an error if the response status code is not 200.
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

// HealthLiveHandler is returning a Gin middleware function that wraps the `LiveEndpoint` method of the
// `health` handler of the `HealthController` struct. The `LiveEndpoint` method is used to check if the
// application is alive and responding to requests. The `gin.WrapH` function is used to convert the
// `http.HandlerFunc` returned by `LiveEndpoint` into a `gin.HandlerFunc` that can be used as
// middleware in a Gin router.
func (c *HealthController) HealthLiveHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.LiveEndpoint))
}

// HealthReadyHandler is returning a Gin middleware function that wraps the `ReadyEndpoint` method of the
// `health` handler of the `HealthController` struct. The `ReadyEndpoint` method is used to check if
// the application is ready to receive requests. The `gin.WrapH` function is used to convert the
// `http.HandlerFunc` returned by `ReadyEndpoint` into a `gin.HandlerFunc` that can be used as
// middleware in a Gin router.
func (c *HealthController) HealthReadyHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.ReadyEndpoint))
}
