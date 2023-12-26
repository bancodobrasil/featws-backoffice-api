package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/database"
	_ "github.com/bancodobrasil/featws-api/docs"
	"github.com/bancodobrasil/featws-api/routes"
	ginMonitor "github.com/bancodobrasil/gin-monitor"
	"github.com/bancodobrasil/goauth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

// This function sets up the logging configuration for a Go program.
func setupLog() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
}

// The following lines provide instructions for generating Swagger documentation.

// ------------------------------

// @title FeatWS API
// @version 1.0
// @Description Este projeto consiste em uma API cujo objetivo é fornecer operações para gerenciamento de repositórios e folhas de regra do sistema FeatWS. Através da API, é possível interagir entre a interface de usuário (UI) e o banco de dados, permitindo diversas interações, como as seguintes:
// @Description - [Post] Criação da Folha de Regra;
// @Description - [Get] Listar das Folhas de Regra;
// @Description - [Get] Obter folha de regra por ID;
// @Description - [Put] Atualizar uma folha de regra por ID;
// @Description - [Delete] Deletar uma folha de regra por ID.
// @Description
// @Description Antes de realizar as requisições no Swagger, é necessário autorizar o acesso clicando no botão **Authorize**, ao lado, e inserindo a senha correspondente. Após inserir o campo **value** e clicar no botão **Authorize**, o Swagger estará disponível para ser utilizado.
// @Description

// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9007
// @BasePath /api/v1

// @securityDefinitions.apikey Authentication Api Key
// @in header
// @name X-API-Key

// @securityDefinitions.apikey Authentication Bearer Token
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}

// end Swagger
// ------------------------------

// This is the main function of a Go program that sets up the configuration, database connection, and
// API routes using the Gin framework.
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

	// Connection with the DataBase containing the Rules and Rulesheets
	database.ConnectDB()

	isCmd := false

	// Perform the database migration
	if cfg.Migrate != "" {
		isCmd = true
		log.Debug("Migrating database...")
		m, err := migrate.New(
			"file://database/migrations/",
			"mysql://"+cfg.MysqlURI,
		)
		if err != nil {
			log.Fatalf("Migration failed with error: %v", err)
			os.Exit(1)
		}
		if strings.ToLower(cfg.Migrate) == "up" {
			if err := m.Up(); err != nil {
				if err.Error() == "no change" {
					log.Println("No change made by migration scripts")
				} else {
					log.Fatalf("Migration Up failed with error: %v", err)
					os.Exit(1)
				}
			}
		} else if strings.ToLower(cfg.Migrate) == "down" {
			if err := m.Down(); err != nil {
				log.Fatalf("Migration Down failed with error: %v", err)
				os.Exit(1)
			}
		} else if steps, err := strconv.Atoi(cfg.Migrate); err == nil {
			m.Steps(steps)
		}
	}

	// Successful migration
	if isCmd == true {
		log.Debug("Finished Successfully")
		os.Exit(0)
	}

	// This code block is setting up monitoring and logging for the Go API using the Gin framework.
	monitor, err := ginMonitor.New("v1.0.0", ginMonitor.DefaultErrorMessageKey, ginMonitor.DefaultBuckets)
	if err != nil {
		log.Panic(err)
	}
	gin.DefaultWriter = log.StandardLogger().WriterLevel(log.DebugLevel)
	gin.DefaultErrorWriter = log.StandardLogger().WriterLevel(log.ErrorLevel)

	router := gin.New()

	goauth.BootstrapMiddleware()

	router.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())
	router.Use(monitor.Prometheus())
	router.GET("metrics", gin.WrapH(promhttp.Handler()))

	// Setup Routers of health resources, swagger and home endpoint
	routes.SetupRoutes(router)
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = strings.Split(cfg.AllowOrigins, ",")
	configCors.AllowHeaders = append(configCors.AllowHeaders, "X-API-Key")
	router.Use(cors.New(configCors))

	// Setup API routers
	routes.APIRoutes(router)

	port := cfg.Port

	router.Run(":" + port)

	log.Infof("Listen on http://0.0.0.0:%s\n", port)
}
