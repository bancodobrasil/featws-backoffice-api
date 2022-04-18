package main

import (
	"os"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/featws-api/routes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

	router := gin.New()

	routes.SetupRoutes(router)

	port := cfg.Port

	router.Run(":" + port)

	log.Infof("Listen on http://0.0.0.0:%s\n", port)
}
