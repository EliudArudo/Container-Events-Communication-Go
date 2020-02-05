package databaseops

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/eliudarudo/event-service/dockerapi"
	"github.com/eliudarudo/event-service/env"
	"github.com/eliudarudo/event-service/interfaces"
	"github.com/eliudarudo/event-service/logs"
	"github.com/eliudarudo/event-service/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var filename = "databaseops/databaseops.go"

func getDatabaseCollection(collectionName string) *(mongo.Collection) {
	mongoClient := mongodb.GetClient()

	collection := mongoClient.Database(env.MongoKeys.Database).Collection(collectionName)

	return collection
}

func getExistingRequestDocumentID(request string) *string {
	collection := getDatabaseCollection("requests")

	foundRequest := interfaces.RequestModelInterface{}

	filter := bson.M{
		"request": bson.M{
			"$eq": request,
		},
	}
	err := collection.FindOne(context.TODO(), filter).Decode(&foundRequest)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getExistingRequestDocumentID", err.Error())
	}

	var finalID string

	if len(foundRequest.Request) > 0 {
		finalID = foundRequest.ID.Hex()
	}

	return &finalID
}

func getExistingTask(task *(interfaces.ReceivedEventInterface)) *interfaces.TaskModelInterface {
	existingTask := interfaces.TaskModelInterface{}

	requestBodyID := *(getExistingRequestDocumentID((*task).RequestBody))

	if len(requestBodyID) > 0 {
		collection := getDatabaseCollection("tasks")

		filter := bson.M{
			"fromContainerService": bson.M{
				"$eq": (*task).Service,
			},
			"task": bson.M{
				"$eq": (*task).Task,
			},
			"subtask": bson.M{
				"$eq": (*task).Subtask,
			},
			"requestBodyId": bson.M{
				"$eq": requestBodyID,
			},
		}

		err := collection.FindOne(context.TODO(), filter).Decode(&existingTask)
		if err != nil {
			logs.StatusFileMessageLogging("FAILURE", filename, "getExistingTask", err.Error())
		}
	}

	return &existingTask
}

func getExistingParsedTask(mongoDBTask *(interfaces.TaskModelInterface)) *interfaces.InitialisedRecordInfoInterface {
	toResponseID := (*mongoDBTask).ToResponseBodyID
	parsedTask := &interfaces.InitialisedRecordInfoInterface{}

	if len(toResponseID) == 0 {
		return parsedTask
	}

	response := interfaces.ResponseModelInterface{}

	collection := getDatabaseCollection("responses")

	docID, err := primitive.ObjectIDFromHex(toResponseID)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getExistingParsedTask", err.Error())
	}

	filter := bson.M{
		"_id": bson.M{
			"$eq": docID,
		},
	}
	err = collection.FindOne(context.TODO(), filter).Decode(&response)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getExistingParsedTask", err.Error())
	}

	parsedTask = &interfaces.InitialisedRecordInfoInterface{
		ContainerID:             (*mongoDBTask).FromContainerID,
		ContainerService:        (*mongoDBTask).FromContainerService,
		RecordID:                (*mongoDBTask).ID.Hex(),
		Task:                    (*mongoDBTask).Task,
		Subtask:                 (*mongoDBTask).Subtask,
		ServiceContainerID:      (*mongoDBTask).ServiceContainerID,
		ServiceContainerService: (*mongoDBTask).ServiceContainerService,
		ResponseBody:            response.Response,
	}

	return parsedTask
}

func getNewParsedTask(mongoDBTask interfaces.TaskModelInterface, selectedContainerInfo interfaces.ContainerInfoStruct) *interfaces.InitialisedRecordInfoInterface {
	parsedTask := interfaces.InitialisedRecordInfoInterface{
		ContainerID:             mongoDBTask.FromContainerID,
		ContainerService:        mongoDBTask.FromContainerService,
		RecordID:                mongoDBTask.ID.Hex(),
		Task:                    mongoDBTask.Task,
		Subtask:                 mongoDBTask.Subtask,
		ServiceContainerID:      mongoDBTask.ServiceContainerID,
		ServiceContainerService: mongoDBTask.ServiceContainerService,
		ChosenContainerID:       selectedContainerInfo.ID,
		ChosenContainerService:  selectedContainerInfo.Service,
	}

	response := interfaces.ResponseModelInterface{}

	collection := getDatabaseCollection("responses")

	filter := bson.M{
		"_id": bson.M{
			"$eq": mongoDBTask.ToResponseBodyID,
		},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&response)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getNewParsedTask", err.Error())
	} else {
		parsedTask.ResponseBody = response.Response
	}

	return &parsedTask
}

func saveNewRequestAndGetID(requestBody string) string {

	collection := getDatabaseCollection("requests")

	newRequest := interfaces.RequestModelInterface{Request: requestBody}

	newRequest.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.TODO(), newRequest)

	var finalID string

	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "saveNewRequestAndGetID", err.Error())
		return finalID
	}

	return newRequest.ID.Hex()

}

func getTargetService(key string) (string, error) {
	jsonFile, err := os.Open("../tasks/task-maps.json")
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	return result[key], nil
}

func recordNewInitialisedTaskWithRequestID(funcTask *(interfaces.ReceivedEventInterface), requestBodyID string) *interfaces.InitialisedRecordInfoInterface {
	parsedTask := &interfaces.InitialisedRecordInfoInterface{}

	strigifiedTask := fmt.Sprintf("%v", (*funcTask).Task)
	targetService, err := getTargetService(strigifiedTask)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "recordNewInitialisedTaskWithRequestID", err.Error())
		return parsedTask
	}

	selectedContainer := dockerapi.FetchConsumingContainer(targetService)
	myContainerInfo := dockerapi.GetMyOfflineContainerInfo()

	newTask := interfaces.TaskModelInterface{
		FromRequestID:           (*funcTask).RequestID,
		FromContainerID:         (*funcTask).ContainerID,
		FromContainerService:    (*funcTask).Service,
		FromReceivedTime:        time.Now(),
		Task:                    (*funcTask).Task,
		Subtask:                 (*funcTask).Subtask,
		RequestBodyID:           requestBodyID,
		ToContainerID:           selectedContainer.ID,
		ToContainerService:      selectedContainer.Service,
		ServiceContainerID:      myContainerInfo.ID,
		ServiceContainerService: myContainerInfo.Service,
		ToReceivedTime:          time.Now(),
		ToResponseBodyID:        "Go Event Service: Nothing received yet",
		FromSentTime:            time.Now(),
	}

	collection := getDatabaseCollection("tasks")

	newTask.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.TODO(), newTask)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "recordNewInitialisedTaskWithRequestID", err.Error())
	}

	parsedTask = getNewParsedTask(newTask, *selectedContainer)
	return parsedTask
}

func recordNewTaskAndRequest(task *(interfaces.ReceivedEventInterface)) *interfaces.InitialisedRecordInfoInterface {
	requestBodyID := saveNewRequestAndGetID((*task).RequestBody)

	initialisedInfo := recordNewInitialisedTaskWithRequestID(task, requestBodyID)

	initialisedInfo.RequestBody = task.RequestBody

	return initialisedInfo
}

func getParsedResponse(funcResponse *(interfaces.ReceivedEventInterface), oldTask *(interfaces.TaskModelInterface)) *interfaces.EventInterface {
	response := interfaces.EventInterface{
		RequestID:    (*oldTask).FromRequestID,
		ContainerID:  (*oldTask).FromContainerID,
		Service:      (*oldTask).FromContainerService,
		ResponseBody: (*funcResponse).ResponseBody,
	}

	return &response
}

func saveNewResponseAndGetID(funcResponse *(interfaces.ReceivedEventInterface)) string {

	collection := getDatabaseCollection("responses")

	newResponse := interfaces.ResponseModelInterface{Response: (*funcResponse).ResponseBody}
	newResponse.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.TODO(), newResponse)

	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "saveNewResponseAndGetID", err.Error())
	}

	return newResponse.ID.Hex()
}

func completeRecordInDB(funcResponse *(interfaces.ReceivedEventInterface), receivedTime time.Time, responseBodyID string) {
	fromSentTime := time.Now()
	collection := getDatabaseCollection("tasks")

	docID, err := primitive.ObjectIDFromHex((*funcResponse).RecordID)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getExistingParsedTask", err.Error())
	}

	filter := bson.M{
		"_id": bson.M{
			"$eq": docID,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"toReceivedTime":   receivedTime,
			"toResponseBodyId": responseBodyID,
			"fromSentTime":     fromSentTime,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getExistingTask", err.Error())
	}
}

// RecordNewTaskInDB checks if there's an existing task and if not, records a new task and request
func RecordNewTaskInDB(task *(interfaces.ReceivedEventInterface)) *interfaces.InitialisedRecordInfoInterface {
	existingTask := getExistingTask(task)

	if len(existingTask.FromRequestID) > 0 {
		parsedTask := getExistingParsedTask(existingTask)
		parsedTask.Existing = true

		return parsedTask
	}

	initRecordInfo := recordNewTaskAndRequest(task)

	return initRecordInfo
}

// CompleteExistingTaskRecordInDB gets response and returns existing response if it exists, otherwise it stores a new response
// and returns parsed object
func CompleteExistingTaskRecordInDB(funcResponse *(interfaces.ReceivedEventInterface)) (*interfaces.EventInterface, error) {

	preexistingResponse := interfaces.ResponseModelInterface{}

	channel1 := make(chan bool)
	channel2 := make(chan *(interfaces.TaskModelInterface))

	go func(funcResponse *(interfaces.ReceivedEventInterface)) {
		collection := getDatabaseCollection("responses")
		filter := bson.M{
			"response": bson.M{
				"$eq": (*funcResponse).ResponseBody,
			},
		}

		err := collection.FindOne(context.TODO(), filter).Decode(&preexistingResponse)
		if err != nil {
			logs.StatusFileMessageLogging("FAILURE", filename, "CompleteExistingTaskRecordInDB", err.Error())
		}

		if len(preexistingResponse.Response) == 0 {
			toReceivedTime := time.Now()
			toResponseBodyID := saveNewResponseAndGetID(funcResponse)
			completeRecordInDB(funcResponse, toReceivedTime, toResponseBodyID)
		}

		channel1 <- true

	}(funcResponse)

	go func(funcResponse *(interfaces.ReceivedEventInterface)) {
		var task interfaces.TaskModelInterface

		collection := getDatabaseCollection("tasks")

		docID, err := primitive.ObjectIDFromHex((*funcResponse).RecordID)
		if err != nil {
			logs.StatusFileMessageLogging("FAILURE", filename, "CompleteExistingTaskRecordInDB", err.Error())
		}

		filter := bson.M{
			"_id": bson.M{
				"$eq": docID,
			},
		}

		err = collection.FindOne(context.TODO(), filter).Decode(&task)
		if err != nil {
			logs.StatusFileMessageLogging("FAILURE", filename, "CompleteExistingTaskRecordInDB", err.Error())
		}

		channel2 <- &task

	}(funcResponse)

	<-channel1
	newTask := <-channel2
	/* ---> funcResponse and task */
	response := getParsedResponse(funcResponse, newTask)
	return response, nil
}
