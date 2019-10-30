package logic

import (
	"encoding/json"
	"fmt"

	"github.com/eliudarudo/event-service/databaseops"
	"github.com/eliudarudo/event-service/env"
	"github.com/eliudarudo/event-service/interfaces"
	"github.com/eliudarudo/event-service/logs"

	"github.com/go-redis/redis"
)

var filename = "logic/logic.go"

// EventDeterminer determines which type of event has been received through redis and channels it to respective handler functions
func EventDeterminer(sentEvent string, containerInfo interfaces.ContainerInfoStruct) {
	var debug1 string
	var event interfaces.ReceivedEventInterface

	json.Unmarshal([]byte(sentEvent), &event)

	// If event is still unmarshalled
	if len(event.ContainerID) == 0 {
		json.Unmarshal([]byte(sentEvent), &debug1)

		if err := json.Unmarshal([]byte(debug1), &event); err != nil {
			logs.StatusFileMessageLogging("FAILURE", filename, "EventDeterminer", err.Error())
		}
	}
	// else event is already marshalled

	eventIsOurs := event.ServiceContainerID == containerInfo.ID && event.ServiceContainerService == containerInfo.Service

	var taskType interfaces.EventTaskType

	if len(event.ResponseBody) > 0 {
		taskType = interfaces.RESPONSE
	} else {
		taskType = interfaces.TASK
	}

	if !eventIsOurs {
		return
	}

	switch taskType {
	case interfaces.TASK:
		recordAndAllocateTask(event)
	case interfaces.RESPONSE:
		modifyDatabaseAndSendBackResponse(event)
	}

}

func recordAndAllocateTask(task interfaces.ReceivedEventInterface) {

	initRecordInfo := databaseops.RecordNewTaskInDB(task)

	if len(initRecordInfo.ContainerID) > 0 && initRecordInfo.Existing {
		responseInfo := getParsedResponseInfo(task, *initRecordInfo)
		sendEventToContainer(*responseInfo)
		return
	}

	allocateTaskToConsumingContainer(*initRecordInfo)
}

func modifyDatabaseAndSendBackResponse(response interfaces.ReceivedEventInterface) {
	responseInfo, err := databaseops.CompleteExistingTaskRecordInDB(response)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "modifyDatabaseAndSendBackResponse", err.Error())
	}

	sendEventToContainer(*responseInfo)
}

func getParsedResponseInfo(task interfaces.ReceivedEventInterface, existingRecordInfo interfaces.InitialisedRecordInfoInterface) *interfaces.EventInterface {
	parsedResponseInfo := interfaces.EventInterface{
		RequestID:    task.RequestID,
		ContainerID:  task.ContainerID,
		Service:      task.Service,
		ResponseBody: existingRecordInfo.ResponseBody,
	}

	return &parsedResponseInfo
}

func allocateTaskToConsumingContainer(initRecordInfo interfaces.InitialisedRecordInfoInterface) {
	eventToSend := parseEventFromRecordInfo(initRecordInfo)

	sendEventToContainer(*eventToSend)
}

func sendEventToContainer(eventInfo interfaces.EventInterface) {

	redisURI := fmt.Sprintf("%v:%v", env.RedisKeys.Host, env.RedisKeys.Port)

	publisher := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer publisher.Close()

	exportedTask := interfaces.ReceivedEventInterface{
		RequestID:               eventInfo.RequestID,
		Task:                    eventInfo.Task,
		Subtask:                 eventInfo.Subtask,
		ContainerID:             eventInfo.ContainerID,
		Service:                 eventInfo.Service,
		RecordID:                eventInfo.RecordID,
		ServiceContainerID:      eventInfo.ServiceContainerID,
		ServiceContainerService: eventInfo.ServiceContainerService,
		ResponseBody:            eventInfo.ResponseBody,
		RequestBody:             eventInfo.RequestBody,
	}

	jsonifiedTask, _ := json.Marshal(exportedTask)
	jsonStringifiedTask := fmt.Sprintf("%#v", string(jsonifiedTask))

	err := publisher.Publish(env.ConsumingService, jsonStringifiedTask).Err()
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "sendTaskToEventsService", err.Error())
	}
}

func parseEventFromRecordInfo(initRecordInfo interfaces.InitialisedRecordInfoInterface) *interfaces.EventInterface {
	event := interfaces.EventInterface{
		ContainerID:             initRecordInfo.ChosenContainerID,
		Service:                 initRecordInfo.ChosenContainerService,
		RecordID:                initRecordInfo.RecordID,
		Task:                    initRecordInfo.Task,
		Subtask:                 initRecordInfo.Subtask,
		ServiceContainerID:      initRecordInfo.ServiceContainerID,
		ServiceContainerService: initRecordInfo.ServiceContainerService,
		RequestBody:             initRecordInfo.RequestBody,
	}

	return &event
}
