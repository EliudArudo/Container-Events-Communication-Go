package env

import (
	"os"
	"strconv"

	"github.com/eliudarudo/event-service/interfaces"
)

// EventService is what our event service is listening to
var EventService = "Event_Service"

// ConsumingService is what the consuming containers are listening to
var ConsumingService = "Consuming_Service"

// RedisKeys connects us to our redis instance
var RedisKeys = interfaces.RedisEnvInterface{Host: "localhost", Port: 6379}

// MongoKeys connects us to our mongodb instance
var MongoKeys = interfaces.MongoEnvInterface{Host: "localhost", Port: 27017, Database: "test"}

// InitialiseEnvironmentVariables initialises our env variables
func InitialiseEnvironmentVariables() {
	initialiseMongoEnv()
	initialiseRedisEnv()
	initialiseRedisChannelEnv()
}

func initialiseRedisEnv() {
	redisHostEnvVariableName := "REDIS_HOST"
	redisHost := os.Getenv(redisHostEnvVariableName)

	redisPortEnvVariableName := "REDIS_PORT"
	redisPort := os.Getenv(redisPortEnvVariableName)

	if len(redisHost) > 0 {
		convertedPort, _ := strconv.Atoi(redisPort)
		RedisKeys = interfaces.RedisEnvInterface{Host: redisHost, Port: convertedPort}
	}
}

func initialiseMongoEnv() {
	mongoHostEnvVariableName := "MONGOURI"
	mongoHost := os.Getenv(mongoHostEnvVariableName)

	mongoPortEnvVariableName := "MONGOPORT"
	mongoPort := os.Getenv(mongoPortEnvVariableName)

	mongoDatabaseEnvVariableName := "MONGODATABASE"
	mongoDatabase := os.Getenv(mongoDatabaseEnvVariableName)

	if len(mongoHost) > 0 {
		convertedPort, _ := strconv.Atoi(mongoPort)
		MongoKeys = interfaces.MongoEnvInterface{Host: mongoHost, Port: convertedPort, Database: mongoDatabase}
	}
}

func initialiseRedisChannelEnv() {
	eventServiceEnvVariableName := "EVENT_SERVICE_EVENT"
	eventServiceEnvVariableValue := os.Getenv(eventServiceEnvVariableName)

	if len(eventServiceEnvVariableValue) > 0 {
		EventService = eventServiceEnvVariableValue
	}

	consumingServiceEnvVariableName := "CONSUMING_SERVICE_EVENT"
	consumingServiceEnvVariableValue := os.Getenv(consumingServiceEnvVariableName)

	if len(consumingServiceEnvVariableValue) > 0 {
		ConsumingService = consumingServiceEnvVariableValue
	}
}
