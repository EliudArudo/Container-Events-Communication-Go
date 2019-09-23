package initialise

// Go initialises http,mongodb and redis
func Go() {

	// initialiseDocker()
	initialiseRedis()
	initialiseMongoDB()
	initialiseRouter()
}
