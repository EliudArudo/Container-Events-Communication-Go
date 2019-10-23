package initialise

import (
	"errors"

	"github.com/eliudarudo/event-service/dockerapi"
)

func initialiseDocker() error {
	myContainerInfo := dockerapi.GetMyContainerInfo()
	if len(myContainerInfo.ID) == 0 {
		return errors.New("Docker containers not initialised")
	}

	return nil
}
