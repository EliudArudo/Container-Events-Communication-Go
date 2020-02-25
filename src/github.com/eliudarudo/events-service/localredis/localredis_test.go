package localredis

import (
	"encoding/json"
	"testing"

	"github.com/eliudarudo/event-service/dockerapi"
	"github.com/eliudarudo/event-service/interfaces"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

func TestParseAndReturnOurEvent(t *testing.T) {
	dummyContainerID := "dummyContainerID"
	dummyContainerService := "dummyContainerService"

	dummySentEvent := interfaces.ReceivedEventInterface{
		ServiceContainerID:      dummyContainerID,
		ServiceContainerService: dummyContainerService,
	}

	dummyContainerInfo := interfaces.ContainerInfoStruct{
		ID:      dummyContainerID,
		Service: dummyContainerService,
	}

	dockerapi.SetMyContainerInfo(&dummySentEvent)

	dummySentEventString, _ := json.Marshal(dummySentEvent)

	receivedEvent := parseAndReturnOurEvent(string(dummySentEventString), &dummyContainerInfo)

	t.Log("\tGiven a our event is sent through redis")

	t.Logf("\t\tTest: \tExpected receivedEvent.ServiceContainerID = '%v'", dummyContainerID)

	if receivedEvent.ServiceContainerID == dummyContainerInfo.ID {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, receivedEvent.ServiceContainerID)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, receivedEvent.ServiceContainerID)
	}
}
