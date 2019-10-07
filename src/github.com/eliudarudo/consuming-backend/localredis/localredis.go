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

	pong, err := RedisClient.Ping().Result()
	fmt.Println(pong, err)

	// defer pubsub.Close()

	pubsub := RedisClient.Subscribe(env.ConsumingService)
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
