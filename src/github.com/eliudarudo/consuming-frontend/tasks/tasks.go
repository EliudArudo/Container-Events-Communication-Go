package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/eliudarudo/consuming-frontend/dockerapi"
	"github.com/eliudarudo/consuming-frontend/env"
	"github.com/eliudarudo/consuming-frontend/interfaces"
	"github.com/eliudarudo/consuming-frontend/logs"
	"github.com/eliudarudo/consuming-frontend/util"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

var filename = "tasks/tasks.go"

var waitingTimeForResponseMS = 50

/*
   Test
*/
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

/*
   Test
*/
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
			subtask = interfaces.SUBTRACT
		} else if isMultiplication {
			subtask = interfaces.MULTIPLY
		} else if isDivision {
			subtask = interfaces.DIVIDE
		}
	}

	return subtask
}

/*
   Test
*/
func getTargetService(key string) (string, error) {
	jsonFile, err := os.Open("tasks/task-maps.json")
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	return result[key], nil
}

/*
   Test
*/
func taskDeterminer(requestBody map[string]interface{}, containerInfo interfaces.ContainerInfoStruct) (interfaces.TaskStruct, error) {
	task := determineTask(requestBody)
	subtask := determineSubtask(task, requestBody)

	requestID := uuid.New().String()

	strigifiedTask := fmt.Sprintf("%v", task)
	targetService, err := getTargetService(strigifiedTask)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "taskDeterminer", err.Error())
	}
	chosenContainer := dockerapi.FetchEventContainer(targetService)

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

/*
   Test
*/
func waitForResult(requestID string, expiresAt int64) *(interfaces.ReceivedEventInterface) {
	response := util.GetResponseFromBuffer(requestID)

	for {
		// Our speed breaker
		time.Sleep(time.Millisecond * time.Duration(waitingTimeForResponseMS))

		timeNow := int64(time.Now().Unix())

		if timeNow >= expiresAt || len((*response).ContainerID) > 0 {
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

// TaskController takes request body from http requests, sends it through redis pubsub and waits for response on the same channel
/*
   Test
*/
func TaskController(decodedRequestBody map[string]interface{}, containerInfo interfaces.ContainerInfoStruct) (interfaces.ResultStruct, error) {
	task, err := taskDeterminer(decodedRequestBody, containerInfo)
	if err != nil {
		return interfaces.ResultStruct{}, err
	}

	sendTaskToEventsService(task)

	channel := make(chan *(interfaces.ReceivedEventInterface))

	go func(task *(interfaces.TaskStruct)) {
		expiresAt := int64(time.Now().Add(5 * time.Second).Unix())
		channel <- waitForResult((*task).RequestID, expiresAt)
	}(&task)

	response := <-channel

	result := interfaces.ResultStruct{"FAIL", "0"}

	if len((*response).ResponseBody) > 0 {
		result = interfaces.ResultStruct{"OK", (*response).ResponseBody}
	}

	return result, nil
}
