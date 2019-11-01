package logic

import (
	"encoding/json"
	"fmt"

	"github.com/eliudarudo/consuming-backend/env"
	"github.com/go-redis/redis"

	"github.com/eliudarudo/consuming-backend/interfaces"
	"github.com/eliudarudo/consuming-backend/logs"
)

var filename = "logic/logic.go"

// EventDeterminer determines which type of event
// has been received through redis
func EventDeterminer(event *(interfaces.ReceivedEventInterface)) {
	var taskType interfaces.EventTaskType

	if len((*event).ResponseBody) > 0 {
		taskType = interfaces.RESPONSE
	} else {
		taskType = interfaces.TASK
	}

	switch taskType {
	case interfaces.TASK:
		// Already has fromContainerID and fromContainerService, to it sends back directly to the
		// event service it got this from
		performTaskAndRespond(event)
	case interfaces.RESPONSE:
		// Empty function here - backend does not receive responses
		// for now
	}

}

func performTaskAndRespond(task *(interfaces.ReceivedEventInterface)) {
	results := performLogic(task)
	sendTaskToEventsService(task, results)
}

func sendTaskToEventsService(task *(interfaces.ReceivedEventInterface), results *string) {
	redisURI := fmt.Sprintf("%v:%v", env.RedisKeys.Host, env.RedisKeys.Port)

	publisher := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer publisher.Close()

	exportedTask := interfaces.ReceivedEventInterface{
		Task:                    (*task).Task,
		Subtask:                 (*task).Subtask,
		ContainerID:             (*task).ContainerID,
		Service:                 (*task).Service,
		RecordID:                (*task).RecordID,
		ServiceContainerID:      (*task).ServiceContainerID,
		ServiceContainerService: (*task).ServiceContainerService,
		ResponseBody:            *results,
	}

	jsonifiedTask, _ := json.Marshal(exportedTask)
	jsonStringifiedTask := fmt.Sprintf("%#v", string(jsonifiedTask))

	err := publisher.Publish(env.EventService, jsonStringifiedTask).Err()
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "sendTaskToEventsService", err.Error())
	}
}

func getObjectKeys(_map map[string]interface{}) []string {
	keys := make([]string, 0, len(_map))
	for k := range _map {
		keys = append(keys, k)
	}

	return keys
}

func performLogic(task *(interfaces.ReceivedEventInterface)) *string {

	var decodedRequestBody map[string]interface{}
	err := json.Unmarshal([]byte((*task).RequestBody), &decodedRequestBody)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "performLogic", err.Error())
	}

	keyArray := getObjectKeys(decodedRequestBody)

	item1 := decodedRequestBody[keyArray[0]]
	item2 := decodedRequestBody[keyArray[1]]

	var result interface{}

	if (*task).Task == interfaces.STRING &&
		(*task).Subtask == interfaces.ADD {
		str1, _ := item1.(string)
		str2, _ := item2.(string)

		result = devAddStrings(str1, str2)
	} else {
		num1, _ := item1.(float64)
		num2, _ := item2.(float64)

		if (*task).Subtask == interfaces.ADD {
			result = devAddNumber(num1, num2)
		} else if (*task).Subtask == interfaces.SUBTRACT {
			result = devSubtractNumber(num1, num2)
		} else if (*task).Subtask == interfaces.MULTIPLY {
			result = devMultiplyNumber(num1, num2)
		} else if (*task).Subtask == interfaces.DIVIDE {
			result = devDivideNumber(num1, num2)
		}
	}

	stringifiedResult := fmt.Sprintf("%v", result)

	return &stringifiedResult
}

func devAddStrings(string1 string, string2 string) string {
	concatString := string1 + string2
	return concatString
}

func devAddNumber(number1 float64, number2 float64) float64 {
	addedNumber := number1 + number2
	return addedNumber
}

func devSubtractNumber(number1 float64, number2 float64) float64 {
	subtractedNumber := number1 - number2
	return subtractedNumber
}

func devMultiplyNumber(number1 float64, number2 float64) float64 {
	multipliedNumber := number1 * number2
	return multipliedNumber
}

func devDivideNumber(number1 float64, number2 float64) float64 {
	dividedNumber := number1 / number2
	return dividedNumber
}
