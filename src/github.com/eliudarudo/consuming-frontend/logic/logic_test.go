package logic

import (
	"testing"

	"github.com/eliudarudo/consuming-frontend/interfaces"
	"github.com/eliudarudo/consuming-frontend/util"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

func TestEventDeterminer(t *testing.T) {
	dummyRequestID := "dummyRequestID"
	dummyResponseBody := "dummyResponseBody"

	dummyEvent := &interfaces.ReceivedEventInterface{
		RequestID:    dummyRequestID,
		ResponseBody: dummyResponseBody,
	}

	EventDeterminer(dummyEvent)

	t.Log("\tGiven our event comes in")

	t.Log("\t\tTest: \tExpected event to be pushed to util.ResponseBuffer")

	if util.ResponseBuffer[0] == dummyRequestID {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, util.ResponseBuffer[0])
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, util.ResponseBuffer[0])
	}

	util.ResponseBuffer = []string{}
}
