package initialise

import (
	"fmt"

	"github.com/eliudarudo/consuming-backend/dockerapi"
	"github.com/eliudarudo/consuming-backend/logs"
)

var filename = "initialise/dockerapi.go"

func printMyContainerInfo() {
	myContainerInfo := dockerapi.GetMyContainerInfo()
	if len(myContainerInfo.ID) == 0 {
		logs.StatusFileMessageLogging("FAILURE", filename, "printMyContainerInfo", "Docker containers not initialised")
		return
	}

	containerInfo := fmt.Sprintf("Docker is working, my container info is : \n %+v", myContainerInfo)

	logs.StatusFileMessageLogging("SUCCESS", filename, "printMyContainerInfo", containerInfo)
}
