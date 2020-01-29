package dockerapi

import (
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"

	"github.com/eliudarudo/consuming-backend/interfaces"
)

// Struct for a DockerAPI object
type Struct struct {
	methodsToCall map[string]bool
}

// ExpectToCall creates a mocking expectation
func (dockerMock *Struct) ExpectToCall(methodName string) {
	if dockerMock.methodsToCall == nil {
		dockerMock.methodsToCall = make(map[string]bool)
	}

	dockerMock.methodsToCall[methodName] = false
}

// Verify checks which methods are called
func (dockerMock *Struct) Verify(t *testing.T) {
	for methodName, called := range dockerMock.methodsToCall {
		if !called {
			t.Errorf("Expected to call '%s', but it was not called", methodName)
		}
	}
}

// Restore clears all method calls
func (dockerMock *Struct) Restore() {
	dockerMock.methodsToCall = nil
}

// SetMyContainerInfo allows for tests to set myContainerInfo
func SetMyContainerInfo(containerInfo *interfaces.ReceivedEventInterface) {
	myContainerInfo = interfaces.ContainerInfoStruct{
		ID:      containerInfo.ContainerID,
		Service: containerInfo.Service,
	}
}

var myContainerInfo interfaces.ContainerInfoStruct

// GetMyContainerInfo gets all docker containers and stores this container's info in the global
// myContainerInfo variable
func GetMyContainerInfo() interfaces.ContainerInfoStruct {
	for {
		initialise()

		if len(myContainerInfo.ID) > 0 {
			break
		}
	}

	return myContainerInfo
}

// GetMyOfflineContainerInfo get's container info from global myContainerInfo variable
// If it does not exist, it reinitialises the container info fetch and returns it
func GetMyOfflineContainerInfo() interfaces.ContainerInfoStruct {
	for {
		if len(myContainerInfo.ID) > 0 {
			break
		}
		GetMyContainerInfo()
	}

	return myContainerInfo
}

func initialise() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	var containerArray []types.Container
	var containerInfo interfaces.ContainerInfoStruct

	for {
		containerArray, _ = cli.ContainerList(context.Background(), types.ContainerListOptions{})
		containerInfo = getMyContainerInfoFromContainerArray(containerArray)

		if len(containerInfo.ID) > 0 {
			break
		}
	}

	myContainerInfo = containerInfo
}

func getMyContainerInfoFromContainerArray(containerArray []types.Container) interfaces.ContainerInfoStruct {
	containerInfo := interfaces.ContainerInfoStruct{}

	shortContainerID, _ := os.Hostname()

	if len(containerArray) > 0 {
		foundIndex := -1

		for index, container := range containerArray {
			if strings.Contains(container.ID, shortContainerID) {
				foundIndex = index
			}
		}

		if foundIndex != -1 {
			containerInfo.ID = containerArray[foundIndex].ID

			containerService := containerArray[foundIndex].Labels["com.docker.swarm.service.name"]

			containerInfo.Service = containerService
		}
	}

	return containerInfo
}
