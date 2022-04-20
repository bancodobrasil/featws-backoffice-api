package main

import (
	"os"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/featws-api/routes"
	ginMonitor "github.com/bancodobrasil/gin-monitor"
	telemetry "github.com/bancodobrasil/gin-telemetry"
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
	routes.SetupRoutes(router)
	router.Use(monitor.Prometheus())
	router.GET("metrics", gin.WrapH(promhttp.Handler()))
	router.Use(telemetry.Middleware("featws-api"))

	port := cfg.Port

	router.Run(":" + port)

	log.Infof("Listen on http://0.0.0.0:%s\n", port)
}
