package initialise

import (
	"github.com/eliudarudo/consuming-frontend/env"
	"github.com/eliudarudo/consuming-frontend/localredis"
)

// Go fetches environment variables, prints our container info, tests mongodb connection and
// sets up our redis pubsub listeners
func Go() {
	env.FetchEnvironmentVariables()
	printMyContainerInfo()
	localredis.SetUpRedisPubSubListener()
	setUpRouter()
}
