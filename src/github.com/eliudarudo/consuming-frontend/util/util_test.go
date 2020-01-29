package util

import (
	"testing"

	"github.com/eliudarudo/consuming-frontend/interfaces"
)

var dummyRequestID = "dummyRequestID"
var dummyResponseBody = "dummyResponseBody"

var dummyEvent = &interfaces.ReceivedEventInterface{
	RequestID:    dummyRequestID,
	ResponseBody: dummyResponseBody,
}

const succeedIcon = "\u2713"
const failIcon = "\u2717"

func TestPushResponseToBuffers(t *testing.T) {
	PushResponseToBuffers(dummyEvent)

	t.Log("\tGiven event is pushed in")

	t.Logf("\t\tTest: \tExpected first event RequestID to be '%v'", dummyEvent.RequestID)

	if ResponseBuffer[0] == dummyEvent.RequestID {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, ResponseBuffer[0])
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, ResponseBuffer[0])
	}

	clearResponseFromBuffers(dummyEvent.RequestID)
}

func TestEventPushedToResponsesSlice(t *testing.T) {
	PushResponseToBuffers(dummyEvent)

	t.Log("\tGiven event is pushed in")

	t.Logf("\t\tTest: \tExpected first response's RequestID to be '%v'", dummyEvent.RequestID)

	if responses[0].RequestID == dummyEvent.RequestID {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, responses[0].RequestID)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, responses[0].RequestID)
	}

	clearResponseFromBuffers(dummyEvent.RequestID)
}

func TestClearResponseFromBuffers(t *testing.T) {
	PushResponseToBuffers(dummyEvent)

	t.Log("\tGiven an existing event is pushed and cleared from ResponseBuffer")

	t.Log("\t\tTest: \tExpected ResponseBuffer.length to be 0")

	clearResponseFromBuffers(dummyEvent.RequestID)

	if len(ResponseBuffer) == 0 {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, len(ResponseBuffer))
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, len(ResponseBuffer))
	}
}

func TestGetResponseFromBuffer(t *testing.T) {
	PushResponseToBuffers(dummyEvent)

	response := GetResponseFromBuffer(dummyEvent.RequestID)

	t.Log("\tGiven an existing event in the buffers")

	t.Logf("\t\tTest: \tExpected response to have a ResponseBody of '%v'", dummyEvent.ResponseBody)

	if response.ResponseBody == dummyEvent.ResponseBody {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, response.ResponseBody)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, response.ResponseBody)
	}

	clearResponseFromBuffers(dummyEvent.RequestID)
}
