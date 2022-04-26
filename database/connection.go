package database

import (
	"context"
	"time"

	"github.com/bancodobrasil/featws-api/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// GetConn ...
func GetConn() *gorm.DB {
	return db
}

// GetModel ...
func GetModel(value interface{}) *gorm.DB {
	return db.Model(value)
}

// ConnectDB ...
func ConnectDB() {

	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatalf("Não foi carregada configuracão!\n")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbConn, err := gorm.Open(mysql.Open(cfg.MysqlURI+"?parseTime=true"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = dbConn.ConnPool.QueryContext(ctx, "SELECT 1")
	if err != nil {
		log.Fatal(err)
	}

	log.Debugln("Connected to Mysql...")
	db = dbConn
}

// GetCollection getting database collections
// func GetCollection(collectionName string) *mongo.Collection {
// 	cfg := config.GetConfig()
// 	if cfg == nil {
// 		log.Fatalf("Não foi carregada configuracão!\n")
// 		return nil
// 	}
// 	collection := client.Database(cfg.MongoDB).Collection(collectionName)
// 	return collection
// }
