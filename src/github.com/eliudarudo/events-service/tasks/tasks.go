package tasks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eliudarudo/event-service/dockerapi"
	"github.com/eliudarudo/event-service/env"
	"github.com/eliudarudo/event-service/interfaces"
	"github.com/eliudarudo/event-service/logs"
	"github.com/eliudarudo/event-service/util"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

var filename = "tasks/tasks.go"

func determineTask(requestBody map[string]interface{}) interfaces.TaskType {
	var task interfaces.TaskType
	var isString bool
	var isNumber bool

	for _, value := range requestBody {
		switch value.(type) {
		case int:
			isString = false
			isNumber = true

		case float64:
			isString = false
			isNumber = true

		case float32:
			isString = false
			isNumber = true

		case string:
			isString = true
			isNumber = false
		}
	}

	if isString {
		task = interfaces.STRING
	} else if isNumber {
		task = interfaces.NUMBER
	}

	return task
}

func determineSubtask(task interfaces.TaskType, requestBody map[string]interface{}) interfaces.SubTaskType {
	var subtask interfaces.SubTaskType

	var isAddition bool
	var isSubtration bool
	var isMultiplication bool
	var isDivision bool

	switch task {
	case interfaces.STRING:
		subtask = interfaces.ADD
	case interfaces.NUMBER:
		for key := range requestBody {
			isAddition = strings.Contains(key, "a")
			isSubtration = strings.Contains(key, "s")
			isMultiplication = strings.Contains(key, "m")
			isDivision = strings.Contains(key, "d")
		}

		if isAddition {
			subtask = interfaces.ADD
		} else if isSubtration {
			subtask = interfaces.ADD
		} else if isMultiplication {
			subtask = interfaces.MULTIPLY
		} else if isDivision {
			subtask = interfaces.DIVIDE
		}
	}

	return subtask
}

func taskDeterminer(requestBody map[string]interface{}, containerInfo interfaces.ContainerInfoStruct) (interfaces.TaskStruct, error) {
	task := determineTask(requestBody)
	subtask := determineSubtask(task, requestBody)

	requestID := uuid.New().String()

	chosenContainer := dockerapi.FetchEventContainer()

	marshalledRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "taskDeterminer", err.Error())
	}

	exportTask := interfaces.TaskStruct{
		task,
		subtask,
		containerInfo.ID,
		containerInfo.Service,
		requestID,
		string(marshalledRequestBody),
		chosenContainer.ID,
		chosenContainer.Service}

	return exportTask, nil

}

func waitForResult(requestID string) interfaces.ReceivedEventInterface {
	response := util.GetResponseFromBuffer(requestID)

	for {
		if len(response.ContainerID) > 0 {
			break
		}
		response = util.GetResponseFromBuffer(requestID)
	}

	return response
}

func sendTaskToEventsService(task interfaces.TaskStruct) {
	redisURI := fmt.Sprintf("%v:%v", env.RedisKeys.Host, env.RedisKeys.Port)

	publisher := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer publisher.Close()

	jsonifiedTask, _ := json.Marshal(task)
	jsonStringifiedTask := fmt.Sprintf("%#v", string(jsonifiedTask))
	
	err := publisher.Publish(env.EventService, jsonStringifiedTask).Err()
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "sendTaskToEventsService", err.Error())
	}
}

// TaskController takes request body and containerInfo and returns result
func TaskController(decodedRequestBody map[string]interface{}, containerInfo interfaces.ContainerInfoStruct) (interfaces.ResultStruct, error) {
	task, err := taskDeterminer(decodedRequestBody, containerInfo)
	if err != nil {
		return interfaces.ResultStruct{}, err
	}

	sendTaskToEventsService(task)

	response := waitForResult(task.RequestID)

	return interfaces.ResultStruct{"OK", response.ResponseBody}, nil
}
