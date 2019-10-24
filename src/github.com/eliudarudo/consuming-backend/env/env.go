package env

import (
	"os"
	"strconv"

	"github.com/eliudarudo/consuming-backend/interfaces"
)

// EventService is the default channel listened to from redis pubsub
var EventService = "Event_Service"

// ConsumingService is the default channel for publishing redis messages
var ConsumingService = "Consuming_Service"

// RedisKeys are the default redis database keys
var RedisKeys = interfaces.RedisEnvInterface{Host: "localhost", Port: 6379}

// FetchEnvironmentVariables checks if we have environment variables set and defaults to default values above
func FetchEnvironmentVariables() {
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
