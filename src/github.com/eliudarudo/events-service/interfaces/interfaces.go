package interfaces

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContainerInfoStruct gets container ID and Service info
type ContainerInfoStruct struct {
	ID      string
	Service string
}

// TODO - remove all unneeded code for each microservice
// // ResultStruct is what TaskController returns after computation
// type ResultStruct struct {
// 	Message string `json:"message"`
// 	Result  string `json:"result"`
// }

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

// MarshalBinary -
func (e *ReceivedEventInterface) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalBinary -
func (e *ReceivedEventInterface) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	return nil
}

// RedisEnvInterface is an interface for our redis env variables
type RedisEnvInterface struct {
	Host string
	Port int
}

// MongoEnvInterface is an interface for our mongodb env variables
type MongoEnvInterface struct {
	Host     string
	Port     int
	Database string
}

// TODO - add 'omitempty' to all the other service interfaces

// InitialisedRecordInfoInterface -
type InitialisedRecordInfoInterface struct {
	ContainerID             string      `json:"containerId"`
	ContainerService        string      `json:"containerService"`
	RecordID                string      `json:"recordId"`
	Task                    TaskType    `json:"task"`
	Subtask                 SubTaskType `json:"subtask"`
	RequestBody             string      `json:"requestBody"`
	ServiceContainerID      string      `json:"serviceContainerId"`
	ServiceContainerService string      `json:"serviceContainerService"`

	Existing               bool   `json:"existing"`
	ResponseBody           string `json:"responseBody"`
	ChosenContainerID      string `json:"chosenContainerId"`
	ChosenContainerService string `json:"chosenContainerService"`
}

// EventInterface -
type EventInterface struct {
	RequestID   string `json:"requestId"`
	ContainerID string `json:"containerId"`
	Service     string `json:"service"`

	// Received user>'this_container'>'event'>'service'>'event'>'this_container'
	ResponseBody string `json:"responseBody"`

	// Received 'event'>'this_container'
	RecordID                string      `json:"recordId"`
	Task                    TaskType    `json:"task"`
	Subtask                 SubTaskType `json:"subtask"`
	ServiceContainerID      string      `json:"serviceContainerId"`
	ServiceContainerService string      `json:"serviceContainerService"`
	RequestBody             string      `json:"requestBody"`
}

// Mongo Models

// RequestModelInterface -
type RequestModelInterface struct {
	ID      primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	Request string             `json:"request" bson:"request"`
}

// TaskModelInterface -
type TaskModelInterface struct {
	ID                      primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	FromRequestID           string             `json:"fromRequestId" bson:"fromRequestId"`
	FromContainerID         string             `json:"fromContainerId" bson:"fromContainerId"`
	FromContainerService    string             `json:"fromContainerService" bson:"fromContainerService"`
	FromReceivedTime        time.Time          `json:"fromReceivedTime" bson:"fromReceivedTime"`
	Task                    TaskType           `json:"task" bson:"task"`
	Subtask                 SubTaskType        `json:"subtask" bson:"subtask"`
	RequestBodyID           string             `json:"requestBodyId" bson:"requestBodyId"`
	ToContainerID           string             `json:"toContainerId" bson:"toContainerId"`
	ToContainerService      string             `json:"toContainerService" bson:"toContainerService"`
	ServiceContainerID      string             `json:"serviceContainerId" bson:"serviceContainerId"`
	ServiceContainerService string             `json:"serviceContainerService" bson:"serviceContainerService"`

	ToReceivedTime   time.Time `json:"toReceivedTime" bson:"toReceivedTime"`
	ToResponseBodyID string    `json:"toResponseBodyId" bson:"toResponseBodyId"`
	FromSentTime     time.Time `json:"fromSentTime" bson:"fromSentTime"`
}

// ResponseModelInterface -
type ResponseModelInterface struct {
	ID       primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	Response string             `json:"response" bson:"response"`
}
