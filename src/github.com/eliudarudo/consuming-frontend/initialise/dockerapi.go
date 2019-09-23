package initialise

import (
	"errors"

	"github.com/eliudarudo/consuming-frontend/dockerapi"
	"github.com/eliudarudo/consuming-frontend/interfaces"
)

func initialiseDocker() error {
	var myContainerInfo interfaces.ContainerInfoStruct = dockerapi.GetMyContainerInfo()
	if len(myContainerInfo.ID) == 0 {
		return errors.New("Docker containers not initialised")
	}

	return nil
}
