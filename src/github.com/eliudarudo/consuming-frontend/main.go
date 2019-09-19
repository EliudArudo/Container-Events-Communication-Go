package main

import (
	"fmt"

	"github.com/eliudarudo/consuming-frontend/dockerapi"
)

func forever() {
	for {
	}
}

func main() {
	myContainerInfo := dockerapi.GetMyContainerInfo()

	fmt.Print("My container info: %v", myContainerInfo)
	// block forever
	go forever()
	select {}
}
