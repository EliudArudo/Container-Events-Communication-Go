package dockerapi

import (
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"

	"github.com/eliudarudo/event-service/interfaces"
)

var myContainerInfo interfaces.ContainerInfoStruct

// GetMyContainerInfo gets all docker containers and stores this container's info in the global
// myContainerInfo variable
func GetMyContainerInfo() *interfaces.ContainerInfoStruct {
	for {
		initialise()

		if len(myContainerInfo.ID) > 0 {
			break
		}
	}

	return &myContainerInfo
}

// GetMyOfflineContainerInfo get's container info from global myContainerInfo variable
// If it does not exist, it reinitialises the container info fetch and returns it
func GetMyOfflineContainerInfo() *interfaces.ContainerInfoStruct {
	for {
		if len(myContainerInfo.ID) > 0 {
			break
		}
		GetMyContainerInfo()
	}

	return &myContainerInfo
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
	containerInfo := &interfaces.ContainerInfoStruct{}

	for {
		containerArray, _ = cli.ContainerList(context.Background(), types.ContainerListOptions{})
		containerInfo = getMyContainerInfoFromContainerArray(containerArray)

		if len(containerInfo.ID) > 0 {
			break
		}
	}

	myContainerInfo = *containerInfo
}

func getMyContainerInfoFromContainerArray(containerArray []types.Container) *interfaces.ContainerInfoStruct {
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

	return &containerInfo
}

// FetchConsumingContainer returns a randomly selected container from target service

func FetchConsumingContainer(containerServiceKeyword string) *interfaces.ContainerInfoStruct {
	freshContainers := getFreshContainers()

	var selectedContainers []interfaces.ContainerInfoStruct

	for _, container := range freshContainers {
		lowerCaseContainerService := strings.ToLower(container.Service)
		containerBelongsToSelectedService := strings.Contains(lowerCaseContainerService, containerServiceKeyword)

		if containerBelongsToSelectedService {
			selectedContainers = append(selectedContainers, container)
		}
	}

	rand.Seed(time.Now().Unix())
	randomIndex := rand.Int() % len(selectedContainers)

	var randomlySelectedContainer interfaces.ContainerInfoStruct = selectedContainers[randomIndex]

	selectedContainer := interfaces.ContainerInfoStruct{ID: randomlySelectedContainer.ID, Service: randomlySelectedContainer.Service}

	for {
		if len(selectedContainer.ID) > 0 {
			break
		} else {
			FetchConsumingContainer(containerServiceKeyword)
		}
	}

	return &selectedContainer
}
