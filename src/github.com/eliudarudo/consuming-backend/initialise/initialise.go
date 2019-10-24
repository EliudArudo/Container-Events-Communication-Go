package initialise

import (
	"github.com/eliudarudo/consuming-backend/env"
	"github.com/eliudarudo/consuming-backend/localredis"
)

// Go fetches environment variables, prints our container info and sets up our redis pubsub listeners
func Go() {
	env.FetchEnvironmentVariables()
	printMyContainerInfo()
	localredis.SetUpRedisPubSubListener()
}
