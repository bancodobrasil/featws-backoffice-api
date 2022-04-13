package database

import (
	"context"
	"time"

	"github.com/bancodobrasil/featws-api/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB ...
func ConnectDB() {

	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatalf("Não foi carregada configuracão!\n")
		return
	}

	cli, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, calcel := context.WithTimeout(context.Background(), 10*time.Second)
	defer calcel()
	err = cli.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = cli.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugln("Connected to MongoDB")
	client = cli
}

var client *mongo.Client

// GetCollection getting database collections
func GetCollection(collectionName string) *mongo.Collection {
	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatalf("Não foi carregada configuracão!\n")
		return nil
	}
	collection := client.Database(cfg.MongoDB).Collection(collectionName)
	return collection
}
