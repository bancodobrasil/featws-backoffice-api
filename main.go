package main

import (
	"os"
	"strings"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/database"
	_ "github.com/bancodobrasil/featws-api/docs"
	"github.com/bancodobrasil/featws-api/routes"
	ginMonitor "github.com/bancodobrasil/gin-monitor"
	telemetry "github.com/bancodobrasil/gin-telemetry"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func setupLog() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
}

// @title FeatWS API
// @version 1.0
// @description API Project to provide operations to manage FeatWS knowledge repositories rules
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9007
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}

// Run start the resolver server with resolverFunc
func main() {

	setupLog()

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Não foi possível carregar as configurações: %s\n", err)
		return
	}

	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatalf("Não foi carregada configuracão!\n")
		return
	}

	database.ConnectDB()

	monitor, err := ginMonitor.New("v1.0.0", ginMonitor.DefaultErrorMessageKey, ginMonitor.DefaultBuckets)
	if err != nil {
		log.Panic(err)
	}
	gin.DefaultWriter = log.StandardLogger().WriterLevel(log.DebugLevel)
	gin.DefaultErrorWriter = log.StandardLogger().WriterLevel(log.ErrorLevel)

	router := gin.New()

	router.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())
	router.Use(monitor.Prometheus())
	router.GET("metrics", gin.WrapH(promhttp.Handler()))
	router.Use(telemetry.Middleware("featws-api"))

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = strings.Split(cfg.AllowOrigins, ",")
	router.Use(cors.New(configCors))

	routes.SetupRoutes(router)

	port := cfg.Port

	router.Run(":" + port)

	log.Infof("Listen on http://0.0.0.0:%s\n", port)
}
