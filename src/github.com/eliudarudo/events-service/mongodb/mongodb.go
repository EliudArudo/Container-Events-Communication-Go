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

var monoMongoClient *(mongo.Client)

// TestMongoDBConnection pings mongodb and checks for a working connection
// If there isn't, it'll crash the container till it gets a connection
func TestMongoDBConnection() {

	mongoURI := fmt.Sprintf("%v://mongodb:%v/%v", env.MongoKeys.Host, env.MongoKeys.Port, env.MongoKeys.Database)

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic("Could not create a new MongoClient")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		panic("Could not connect to MongoDB")
	}
	defer cancel()

	logs.StatusFileMessageLogging("SUCCESS", filename, "initialiseMongoDB", "Successfully connected to MongoDB")
}

// GetClient establishes a connection to mongodb and returns a working Client
func GetClient() *(mongo.Client) {

	if monoMongoClient == nil {
		mongoURI := fmt.Sprintf("%v://mongodb:%v/%v", env.MongoKeys.Host, env.MongoKeys.Port, env.MongoKeys.Database)
		optionsObject := options.Client().ApplyURI(mongoURI)
		mongoClient, err := mongo.NewClient(optionsObject)
		// Just as an option here
		// optionsObject.SetMaxPoolSize(10)
		if err != nil {
			logs.StatusFileMessageLogging("FAILURE", filename, "GetClient", "Failed to create a new MongoClient")
		}

		err = mongoClient.Connect(context.TODO())
		if err != nil {
			logs.StatusFileMessageLogging("FAILURE", filename, "GetClient", "Failed to connect to the new MongoClient")
		}

		monoMongoClient = mongoClient

		return mongoClient
	}

	return monoMongoClient

}
