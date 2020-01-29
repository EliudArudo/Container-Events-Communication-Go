package dockerapi

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/eliudarudo/consuming-backend/interfaces"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

func TestGetMyOfflineContainerInfo(t *testing.T) {
	dummyContainerID := "dummyContainerId"
	myContainerInfo.ID = dummyContainerID

	foundContainerInfo := GetMyOfflineContainerInfo()

	t.Log("\tGiven already existing myContainerInfo")

	t.Logf("\t\tTest: \tExpected myContainerInfo.ID = '%v'", dummyContainerID)

	if foundContainerInfo.ID == dummyContainerID {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, foundContainerInfo.ID)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, foundContainerInfo.ID)
	}
}

func TestGetParsedContainers(t *testing.T) {

	dummyContainerID := "dummyDockerContainerID"
	dummyContainerService := "dummyDockerService"

	dummyDockerContainers := []types.Container{
		{
			ID: dummyContainerID,
			Labels: map[string]string{
				"com.docker.swarm.service.name": dummyContainerService,
			},
		},
	}

	parsedContainers := []interfaces.ContainerInfoStruct{
		{
			ID:      dummyContainerID,
			Service: dummyContainerService,
		},
	}

	t.Log("\tGiven raw docker containers")

	t.Log("\t\tTest: \tExpected dummyDockerContainers to be parsed to parsedContainers")

	if parsedContainers[0].ID == dummyDockerContainers[0].ID {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, parsedContainers[0].ID)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, parsedContainers[0].ID)
	}

}

func TestGetMyContainerInfoFromContainerArray(t *testing.T) {

	shortContainerID, _ := os.Hostname()
	dummyContainerService := "dummyDockerService"

	dummyDockerContainers := []types.Container{
		{
			ID: shortContainerID,
			Labels: map[string]string{
				"com.docker.swarm.service.name": dummyContainerService,
			},
		},
	}

	foundContainerInfo := getMyContainerInfoFromContainerArray(dummyDockerContainers)

	foundContainerInfoString := fmt.Sprintf("%#v", foundContainerInfo)

	t.Log("\tGiven raw docker containers")

	t.Logf("\t\tTest: \tExpected containerInfo to be : \n%v", foundContainerInfoString)

	if foundContainerInfo.Service == dummyDockerContainers[0].Labels["com.docker.swarm.service.name"] {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, foundContainerInfo.Service)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, foundContainerInfo.Service)
	}
}
