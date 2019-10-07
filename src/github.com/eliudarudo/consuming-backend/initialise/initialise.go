package initialise

import (
	"github.com/eliudarudo/consuming-frontend/env"
	"github.com/eliudarudo/consuming-frontend/localredis"
)

// Go initialises http,mongodb and redis
func Go() {
	env.InitialiseEnvironmentVariables()
	// initialiseDocker()
	localredis.InitialiseRedis()
}
