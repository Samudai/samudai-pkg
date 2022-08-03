package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Samudai/samudai-pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitMongo() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_URL")).
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	logger.LogMessage("info", "mongo db connected")
}

func GetMongo() *mongo.Client {
	return client
}

func CloseMongo() {
	client.Disconnect(context.TODO())
}
