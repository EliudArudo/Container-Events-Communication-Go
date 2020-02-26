package databaseops

import (
	"testing"

	"github.com/eliudarudo/event-service/interfaces"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

func TestGetTargetService(t *testing.T) {
	expectedNumberService := "backend"
	expectedStringService := "backend"

	numberTask := "NUMBER"
	stringTask := "STRING"

	gotNumberService, _ := getTargetService(numberTask)
	gotStringService, _ := getTargetService(stringTask)

	t.Logf("\tGiven a %v task", numberTask)

	t.Logf("\t\tTest: \tExpected targetService  = '%v'", expectedNumberService)

	if gotNumberService == expectedNumberService {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, gotNumberService)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, gotNumberService)
	}

	t.Logf("\tGiven a %v task", stringTask)

	t.Logf("\t\tTest: \tExpected targetService  = '%v'", expectedStringService)

	if gotStringService == expectedStringService {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, gotStringService)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, gotStringService)
	}

}

func TestGetParsedResponse(t *testing.T) {
	dummyFromRequestID := "dummyRequestID"
	dummyFromContainerID := "dummyContainerID"
	dummyFromContainerService := "dummyContainerService"
	dummyResponseBody := "ab"

	dummyOldTask := &interfaces.TaskModelInterface{
		FromRequestID:        dummyFromRequestID,
		FromContainerID:      dummyFromContainerID,
		FromContainerService: dummyFromContainerService,
	}

	dummyResponse := &interfaces.ReceivedEventInterface{
		ResponseBody: dummyResponseBody,
	}

	parsedResponse := getParsedResponse(dummyResponse, dummyOldTask)

	t.Logf("\tGiven a Task with FromRequestID: %v and Response with ResponseBody: %v", dummyOldTask.FromRequestID, dummyResponse.ResponseBody)

	t.Logf("\t\tTest: \tExpected parsedResponse.RequestID  = '%v' and parsedResponse.ResponseBody = '%v'", dummyOldTask.FromRequestID, dummyResponse.ResponseBody)

	if parsedResponse.RequestID == dummyOldTask.FromRequestID && parsedResponse.ResponseBody == dummyResponse.ResponseBody {
		t.Logf("\t\t%v Got : '%v' and '%v'", succeedIcon, parsedResponse.RequestID, parsedResponse.ResponseBody)
	} else {
		t.Errorf("\t\t%v Got : '%v' and '%v'", failIcon, parsedResponse.RequestID, parsedResponse.ResponseBody)
	}
}
