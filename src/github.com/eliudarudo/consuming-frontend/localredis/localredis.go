package localredis

import (
	"encoding/json"
	"fmt"

	"github.com/eliudarudo/consuming-frontend/dockerapi"
	"github.com/eliudarudo/consuming-frontend/env"
	"github.com/eliudarudo/consuming-frontend/interfaces"
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

	myContainerInfo := dockerapi.GetMyOfflineContainerInfo()
	event := interfaces.ReceivedEventInterface{}

	go func() {
		pubsub := redisClient.Subscribe(env.ConsumingService)
		for {
			message, err := pubsub.ReceiveMessage()
			if err != nil {
				logs.StatusFileMessageLogging("FAILURE", redisFilename, "initialiseRedis", err.Error())
			}

			event = *(parseAndReturnOurEvent(message.Payload, &myContainerInfo))

			if len(event.Service) == 0 {
				continue
			}

			go func() {
				logic.EventDeterminer(&event)
			}()
		}
	}()
}

func parseAndReturnOurEvent(sentEvent string, containerInfo *(interfaces.ContainerInfoStruct)) *interfaces.ReceivedEventInterface {
	var debug1 string
	var event interfaces.ReceivedEventInterface

	json.Unmarshal([]byte(sentEvent), &event)

	// If event is still unmarshalled
	if len(event.ContainerID) == 0 {
		json.Unmarshal([]byte(sentEvent), &debug1)

		if err := json.Unmarshal([]byte(debug1), &event); err != nil {
			logs.StatusFileMessageLogging("FAILURE", redisFilename, "EventDeterminer", err.Error())
		}
	}
	// else event is already marshalled

	eventIsOurs := event.ContainerID == containerInfo.ID && event.Service == containerInfo.Service

	if eventIsOurs {
		return &event
	}

	return &interfaces.ReceivedEventInterface{}
}
