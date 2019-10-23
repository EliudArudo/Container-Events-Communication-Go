package localredis

import (
	"fmt"

	"github.com/eliudarudo/event-service/dockerapi"
	"github.com/eliudarudo/event-service/env"
	"github.com/eliudarudo/event-service/logic"
	"github.com/eliudarudo/event-service/logs"
	"github.com/go-redis/redis"
)

var redisFilename = "redis/localredis.go"

// RedisClient is what we export outside
var RedisClient *redis.Client

// InitialiseRedis is a function that fills up client with right info
func InitialiseRedis() {
	redisURI := fmt.Sprintf("%v:%v", env.RedisKeys.Host, env.RedisKeys.Port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// Error may occur here
	// defer client.Close()

	_, err := RedisClient.Ping().Result()
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", redisFilename, "initialiseRedis", err.Error())
	}
	logs.StatusFileMessageLogging("SUCCESS", redisFilename, "initialiseRedis", "Redis successfully set up")

	// defer pubsub.Close()

	pubsub := RedisClient.Subscribe(env.EventService)
	for {
		message, err := pubsub.ReceiveMessage()
		go func() {
			if err != nil {
				logs.StatusFileMessageLogging("FAILURE", redisFilename, "initialiseRedis", err.Error())
				panic(err)
			}

			myContainerInfo := dockerapi.GetMyOfflineContainerInfo()
			logic.EventDeterminer(message.Payload, myContainerInfo)
		}()
	}

}
