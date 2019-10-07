package env

import (
	"os"
	"strconv"

	"github.com/eliudarudo/consuming-frontend/interfaces"
)

// EventService is what our event service is listening to
var EventService = "Event_Service"

// ConsumingService is what the consuming containers are listening to
var ConsumingService = "Consuming_Service"

// RedisKeys connects us to our redis instance
var RedisKeys = interfaces.RedisEnvInterface{Host: "localhost", Port: 6379}

// Port is our Gorilla mux http entry point
var Port = 4000

// InitialiseEnvironmentVariables initialises our env variables
func InitialiseEnvironmentVariables() {
	initialiseRedisEnv()
	initialisePortEnv()
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

func initialisePortEnv() {
	portEnvVariableName := "PORT"
	port := os.Getenv(portEnvVariableName)

	if len(port) > 0 {
		Port, _ = strconv.Atoi(port)
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
