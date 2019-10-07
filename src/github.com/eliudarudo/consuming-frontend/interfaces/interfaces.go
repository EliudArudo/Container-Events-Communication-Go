package interfaces

import "encoding/json"

// ContainerInfoStruct gets container ID and Service info
type ContainerInfoStruct struct {
	ID      string
	Service string
}

// ResultStruct is what TaskController returns after computation
type ResultStruct struct {
	Message string `json:"message"`
	Result  string `json:"result"`
}

// TaskType either 'Number' or 'String'
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

// EventTaskType either 'Task' or 'Response'
type EventTaskType int

const (
	// TASK has requestBody
	TASK EventTaskType = iota
	// RESPONSE has responseBody
	RESPONSE
)

// TaskStruct is the determined task to be sent to event service
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

// MarshalBinary -
func (e *TaskStruct) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalBinary -
func (e *TaskStruct) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	return nil
}

// ReceivedEventInterface is for events received
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

// RedisEnvInterface is an interface for our redis env variables
type RedisEnvInterface struct {
	Host string
	Port int
}
