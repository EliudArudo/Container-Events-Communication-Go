package logic

import (
	"encoding/json"

	"github.com/eliudarudo/consuming-frontend/interfaces"
	"github.com/eliudarudo/consuming-frontend/logs"
	"github.com/eliudarudo/consuming-frontend/util"
)

var filename = "logic/logic.go"

// EventDeterminer determines which type of event
// has been received through redis
func EventDeterminer(sentEvent string, containerInfo interfaces.ContainerInfoStruct) {
	var debug1 string
	var event interfaces.ReceivedEventInterface

	json.Unmarshal([]byte(sentEvent), &debug1)

	if err := json.Unmarshal([]byte(debug1), &event); err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "EventDeterminer", err.Error())
	}

	eventIsOurs := event.ContainerID == containerInfo.ID && event.Service == containerInfo.Service

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
		// Empty function here - frontend does not receive task requests
		// for now
	case interfaces.RESPONSE:
		util.PushResponseToBuffers(event)
	}
}
