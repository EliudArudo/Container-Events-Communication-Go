package initialise

import (
	"github.com/eliudarudo/event-service/env"
	"github.com/eliudarudo/event-service/localredis"
	"github.com/eliudarudo/event-service/mongodb"
)

// Go initialises http,mongodb and redis
func Go() {
	env.InitialiseEnvironmentVariables()
	initialiseDocker()
	mongodb.InitialiseMongoDB()
	localredis.InitialiseRedis()
}
