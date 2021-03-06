package dockerapi

import (
	"errors"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"

	"github.com/eliudarudo/consuming-frontend/interfaces"
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

// FetchEventContainer should return fetched event container
func (dockerMock *Struct) FetchEventContainer(targetService string) interfaces.ContainerInfoStruct {
	return FetchEventContainer(targetService)
}

var myContainerInfo interfaces.ContainerInfoStruct

// SetMyContainerInfo allows for tests to set myContainerInfo
func SetMyContainerInfo(containerInfo *interfaces.ReceivedEventInterface) {
	myContainerInfo = interfaces.ContainerInfoStruct{
		ID:      containerInfo.ContainerID,
		Service: containerInfo.Service,
	}
}

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

func getParsedContainers(containerArray []types.Container) ([]interfaces.ContainerInfoStruct, error) {
	if len(containerArray) == 0 {
		return nil, errors.New("No containers to parse")
	}

	parsedContainers := make([]interfaces.ContainerInfoStruct, len(containerArray))

	for index, container := range containerArray {
		parsedContainers[index].ID = container.ID

		containerService := container.Labels["com.docker.swarm.service.name"]

		parsedContainers[index].Service = containerService
	}

	return parsedContainers, nil
}

func getFreshContainers() []interfaces.ContainerInfoStruct {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	var containerArray []types.Container
	var parsedContainers []interfaces.ContainerInfoStruct

	for {
		containerArray, _ = cli.ContainerList(context.Background(), types.ContainerListOptions{})
		parsedContainers, err = getParsedContainers(containerArray)
		if err == nil && len(parsedContainers) > 0 {
			break
		}
	}

	return parsedContainers
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

// TODO - make sure we can select event service here

// FetchEventContainer returns a randomly selected container from target service
func FetchEventContainer(serviceKeyword string) interfaces.ContainerInfoStruct {
	freshContainers := getFreshContainers()

	var selectedContainers []interfaces.ContainerInfoStruct

	for _, container := range freshContainers {
		lowerCaseContainerService := strings.ToLower(container.Service)
		serviceKeyword = strings.ToLower(serviceKeyword)
		containerBelongsToSelectedService := strings.Contains(lowerCaseContainerService, serviceKeyword)

		if containerBelongsToSelectedService {
			selectedContainers = append(selectedContainers, container)
		}
	}

	for {
		if len(selectedContainers) > 0 {
			break
		}
		FetchEventContainer(serviceKeyword)
	}

	rand.Seed(time.Now().Unix())
	randomIndex := rand.Int() % len(selectedContainers)

	randomlySelectedContainer := selectedContainers[randomIndex]

	selectedContainer := interfaces.ContainerInfoStruct{ID: randomlySelectedContainer.ID, Service: randomlySelectedContainer.Service}

	return selectedContainer
}
