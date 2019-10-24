package interfaces

import "encoding/json"

// ContainerInfoStruct maps a docker container ID and Service
type ContainerInfoStruct struct {
	ID      string
	Service string
}

// TaskType determines which service is selected
type TaskType string

const (
	// NUMBER is everything except a1 and a2 strings
	NUMBER TaskType = "NUMBER"
	// STRING is a1 and a2 strings
	STRING TaskType = "STRING"
)

// SubTaskType specifies a subtask to task
type SubTaskType string

const (
	// ADD using a1 and a2
	ADD SubTaskType = "ADD"
	// MULTIPLY using m1 and m2
	MULTIPLY = "MULTIPLY"
	// SUBTRACT using s1 and s2
	SUBTRACT = "SUBTRACT"
	// DIVIDE using d1 and d2
	DIVIDE = "DIVIDE"
)

// EventTaskType determines next action on event being received from consuming containers
type EventTaskType int

const (
	// TASK determined by presense of requestBody field in ReceivedEventInterface object
	TASK EventTaskType = iota
	// RESPONSE determined by presense of responseBody field in ReceivedEventInterface object
	RESPONSE
)

// TaskStruct is the most basic form of a Task for internal processing
type TaskStruct struct {
	Task                    TaskType    `json:"task"`
	Subtask                 SubTaskType `json:"subtask"`
	ContainerID             string      `json:"containerId"`
	Service                 string      `json:"service"`
	RequestID               string      `json:"requestId"`
	RequestBody             string      `json:"requestBody"`
	ServiceContainerID      string      `json:"serviceContainerId"`
	ServiceContainerService string      `json:"serviceContainerService"`
}

// MarshalBinary marshals []byte
func (e *TaskStruct) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalBinary unmarshals a []byte
func (e *TaskStruct) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	return nil
}

// ReceivedEventInterface is structure of object we've just received from redis pubsub
type ReceivedEventInterface struct {
	RequestID               string      `json:"requestId"`
	ContainerID             string      `json:"containerId"`
	Service                 string      `json:"service"`
	ResponseBody            string      `json:"responseBody"`
	RecordID                string      `json:"recordId"`
	Task                    TaskType    `json:"task"`
	Subtask                 SubTaskType `json:"subtask"`
	RequestBody             string      `json:"requestBody"`
	ServiceContainerID      string      `json:"serviceContainerId"`
	ServiceContainerService string      `json:"serviceContainerService"`
}

// RedisEnvInterface defines Host and Port fields for redis keys
type RedisEnvInterface struct {
	Host string
	Port int
}
