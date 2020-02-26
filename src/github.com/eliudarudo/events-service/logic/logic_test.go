package logic

import (
	"testing"

	"github.com/eliudarudo/event-service/interfaces"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

func TestGetParsedResponseInfo(t *testing.T) {
	dummyRequestID := "dummyRequestID"
	dummyContainerID := "dummyContainerID"
	dummyContainerService := "dummyContainerService"
	dummyResponseBody := "cd"

	dummyTask := &interfaces.ReceivedEventInterface{
		RequestID:   dummyRequestID,
		ContainerID: dummyContainerID,
		Service:     dummyContainerService,
	}

	dummyExistingRecordInfo := &interfaces.InitialisedRecordInfoInterface{
		ResponseBody: dummyResponseBody,
	}

	parsedResponseInfo := getParsedResponseInfo(dummyTask, dummyExistingRecordInfo)

	t.Logf("\tGiven a Task with RequestID: %v and existingRecordInfo with ResponseBody: %v", dummyTask.RequestID, dummyExistingRecordInfo.ResponseBody)

	t.Logf("\t\tTest: \tExpected parsedResponse.RequestID  = '%v' and parsedResponse.ResponseBody = '%v'", dummyTask.RequestID, dummyExistingRecordInfo.ResponseBody)

	if parsedResponseInfo.RequestID == dummyTask.RequestID && parsedResponseInfo.ResponseBody == dummyExistingRecordInfo.ResponseBody {
		t.Logf("\t\t%v Got : '%v' and '%v'", succeedIcon, parsedResponseInfo.RequestID, parsedResponseInfo.ResponseBody)
	} else {
		t.Errorf("\t\t%v Got : '%v' and '%v'", failIcon, parsedResponseInfo.RequestID, parsedResponseInfo.ResponseBody)
	}
}

func TestParseEventFromRecordInfo(t *testing.T) {
	dummyContainerID := "dummyContainerID"
	dummyContainerService := "dummyContainerService"

	dummyInitRecordInfo := interfaces.InitialisedRecordInfoInterface{
		ChosenContainerID:      dummyContainerID,
		ChosenContainerService: dummyContainerService,
	}

	event := parseEventFromRecordInfo(dummyInitRecordInfo)

	t.Logf("\tGiven a Initialised Record Info with ChosenContainerID: %v", dummyInitRecordInfo.ChosenContainerID)

	t.Logf("\t\tTest: \tExpected event.ContainerID  = '%v'", dummyInitRecordInfo.ChosenContainerID)

	if event.ContainerID == dummyInitRecordInfo.ChosenContainerID {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, event.ContainerID)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, event.ContainerID)
	}
}
