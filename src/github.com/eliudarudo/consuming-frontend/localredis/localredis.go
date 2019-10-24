package localredis

import (
	"fmt"

	"github.com/eliudarudo/consuming-frontend/dockerapi"
	"github.com/eliudarudo/consuming-frontend/env"
	"github.com/eliudarudo/consuming-frontend/logic"
	"github.com/eliudarudo/consuming-frontend/logs"
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

	go func() {
		pubsub := redisClient.Subscribe(env.ConsumingService)
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
	}()
}
