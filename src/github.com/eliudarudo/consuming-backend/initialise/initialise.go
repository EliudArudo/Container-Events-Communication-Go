package initialise

import (
	"github.com/eliudarudo/consuming-backend/env"
	"github.com/eliudarudo/consuming-backend/localredis"
)

// Go initialises http,mongodb and redis
func Go() {
	env.InitialiseEnvironmentVariables()
	// initialiseDocker()
	localredis.InitialiseRedis()
}
