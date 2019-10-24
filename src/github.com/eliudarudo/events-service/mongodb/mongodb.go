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
	mongoURI := fmt.Sprintf("%v://mongodb:%v/%v", env.MongoKeys.Host, env.MongoKeys.Port, env.MongoKeys.Database)

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "GetClient", "Failed to create a new MongoClient")
	}

	err = mongoClient.Connect(context.TODO())
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "GetClient", "Failed to connect to the new MongoClient")
	}

	return mongoClient
}
