package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var fullContainerID string

type containerInfoStruct struct {
	id      string
	service string
}

func getMyContainerInfoFromContainerArray(containerArray []types.Container) containerInfoStruct {
	containerInfo := containerInfoStruct{}

	shortContainerID, _ := os.Hostname()

	if len(containerArray) > 0 {
		foundIndex := -1

		for index, container := range containerArray {
			if strings.Contains(container.ID, shortContainerID) {
				foundIndex = index
			}
		}

		if foundIndex != -1 {
			containerInfo.id = containerArray[foundIndex].ID

			containerService := containerArray[foundIndex].Labels["com.docker.swarm.service.name"]

			containerInfo.service = containerService
		}
	}

	return containerInfo
}

func getAndSetMyContainerID() {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	var containerArray []types.Container
	var containerInfo containerInfoStruct

	for {
		containerArray, _ = cli.ContainerList(context.Background(), types.ContainerListOptions{})
		containerInfo = getMyContainerInfoFromContainerArray(containerArray)

		if len(containerInfo.id) > 0 {
			break
		}
	}
	fmt.Printf("containerArray length is:%v\n", len(containerArray))

	fmt.Printf("My container info: %+v\n", containerInfo)

	fullContainerID = containerInfo.id

}

func forever() {
	for {
	}
}

func main() {
	getAndSetMyContainerID()

	fmt.Println("Service B running")
	// block forever
	go forever()
	select {}
}
