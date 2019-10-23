package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eliudarudo/event-service/env"
	"github.com/eliudarudo/event-service/logs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var filename = "mongodb/mongodb.go"

// MongoClient is what we export
var MongoClient *(mongo.Client)

// InitialiseMongoDB initialised the client
func InitialiseMongoDB() {

	mongoURI := fmt.Sprintf("%v://mongodb:%v/%v", env.MongoKeys.Host, env.MongoKeys.Port, env.MongoKeys.Database)

	MongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "initialiseMongoDB", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	err = MongoClient.Connect(ctx)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "initialiseMongoDB", err.Error())
	}
	defer cancel()

	logs.StatusFileMessageLogging("SUCCESS", filename, "initialiseMongoDB", "Successfully connected to MongoDB")
}

// GetClient -
func GetClient() *(mongo.Client) {
	mongoURI := fmt.Sprintf("%v://mongodb:%v/%v", env.MongoKeys.Host, env.MongoKeys.Port, env.MongoKeys.Database)

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "initialiseMongoDB", err.Error())
	}

	err = mongoClient.Connect(context.TODO())
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "initialiseMongoDB", err.Error())
	}

	return mongoClient

}
