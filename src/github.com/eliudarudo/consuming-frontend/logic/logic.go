package logic

import (
	"github.com/eliudarudo/consuming-frontend/interfaces"
	"github.com/eliudarudo/consuming-frontend/util"
)

var filename = "logic/logic.go"

// EventDeterminer determines which type of event
// has been received through redis
/*
   Test
    - util.PushResponseToBuffers called if it's a TASK
*/
func EventDeterminer(event *(interfaces.ReceivedEventInterface)) {
	var taskType interfaces.EventTaskType

	if len((*event).ResponseBody) > 0 {
		taskType = interfaces.RESPONSE
	} else {
		taskType = interfaces.TASK
	}

	switch taskType {
	case interfaces.TASK:
		// Empty function here - frontend does not receive task requests
		// for now
	case interfaces.RESPONSE:
		util.PushResponseToBuffers(event)
	}
}
