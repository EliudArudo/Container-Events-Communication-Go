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

// SetUpRedisPubSubListener checks if we're connected to redis and panics on failure
// It also set's up goroutines to listen to redis subscription messages
func SetUpRedisPubSubListener() {
	redisURI := fmt.Sprintf("%v:%v", env.RedisKeys.Host, env.RedisKeys.Port)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		panic("Cannot connect to Redis")
	}
	logs.StatusFileMessageLogging("SUCCESS", redisFilename, "initialiseRedis", "Redis successfully set up")

	pubsub := redisClient.Subscribe(env.EventService)
	for {
		message, err := pubsub.ReceiveMessage()
		go func() {
			if err != nil {
				logs.StatusFileMessageLogging("FAILURE", redisFilename, "initialiseRedis", err.Error())
			}

			myContainerInfo := dockerapi.GetMyOfflineContainerInfo()
			logic.EventDeterminer(message.Payload, *myContainerInfo)
		}()
	}

}
