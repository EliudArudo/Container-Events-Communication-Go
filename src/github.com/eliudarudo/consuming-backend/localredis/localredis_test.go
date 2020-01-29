package localredis

import (
	"encoding/json"
	"testing"

	"github.com/eliudarudo/consuming-backend/dockerapi"
	"github.com/eliudarudo/consuming-backend/interfaces"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

func TestParseAndReturnOurEvent(t *testing.T) {
	dummyContainerID := "dummyContainerID"
	dummyContainerService := "dummyContainerService"

	dummySentEvent := interfaces.ReceivedEventInterface{
		ContainerID: dummyContainerID,
		Service:     dummyContainerService,
	}

	dummyContainerInfo := interfaces.ContainerInfoStruct{
		ID:      dummyContainerID,
		Service: dummyContainerService,
	}

	dockerapi.SetMyContainerInfo(&dummySentEvent)

	dummySentEventString, _ := json.Marshal(dummySentEvent)

	receivedEvent := parseAndReturnOurEvent(string(dummySentEventString), &dummyContainerInfo)

	t.Log("\tGiven a our event is sent through redis")

	t.Logf("\t\tTest: \tExpected receivedEvent.ContainerID = '%v'", dummyContainerID)

	if receivedEvent.ContainerID == dummyContainerInfo.ID {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, receivedEvent.ContainerID)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, receivedEvent.ContainerID)
	}
}
