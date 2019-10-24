package interfaces

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// MarshalBinary marshals []byte
func (e *ReceivedEventInterface) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalBinary unmarshals a []byte
func (e *ReceivedEventInterface) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	return nil
}

// RedisEnvInterface defines Host and Port fields for redis keys
type RedisEnvInterface struct {
	Host string
	Port int
}

// MongoEnvInterface defines Host, Port and Database fields for redis keys
type MongoEnvInterface struct {
	Host     string
	Port     int
	Database string
}

// InitialisedRecordInfoInterface defines a Task record that has just been received from redis pubsub
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

// EventInterface is a base event interface for all our event-based objects
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

// Mongo DB Models

// RequestModelInterface is a model for our mongodb 'requests' collection
type RequestModelInterface struct {
	ID      primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	Request string             `json:"request" bson:"request"`
}

// TaskModelInterface is a model for our mongodb 'tasks' collection
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

// ResponseModelInterface is a model for our mongodb 'responses' collection
type ResponseModelInterface struct {
	ID       primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	Response string             `json:"response" bson:"response"`
}
