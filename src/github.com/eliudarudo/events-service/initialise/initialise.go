package initialise

import (
	"github.com/eliudarudo/event-service/env"
	"github.com/eliudarudo/event-service/localredis"
	"github.com/eliudarudo/event-service/mongodb"
)

// Go fetches environment variables, prints our container info, tests mongodb connection and
// sets up our redis pubsub listeners
func Go() {
	env.FetchEnvironmentVariables()
	printMyContainerInfo()
	mongodb.TestMongoDBConnection()
	localredis.SetUpRedisPubSubListener()
}
