package databaseops

import (
	"context"
	"fmt"
	"time"

	"github.com/eliudarudo/event-service/dockerapi"
	"github.com/eliudarudo/event-service/env"
	"github.com/eliudarudo/event-service/interfaces"
	"github.com/eliudarudo/event-service/logs"
	"github.com/eliudarudo/event-service/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var filename = "databaseops/databaseops.go"

func getExistingRequestDocumentID(request string) string {

	mongoClient := mongodb.GetClient()

	collection := mongoClient.Database(env.MongoKeys.Database).Collection("requests")

	foundRequest := interfaces.RequestModelInterface{}

	// filter := bson.D{{"request", request}}
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

	return finalID
}

func getExistingTask(task interfaces.ReceivedEventInterface) interfaces.TaskModelInterface {
	existingTask := interfaces.TaskModelInterface{}

	requestBodyID := getExistingRequestDocumentID(task.RequestBody)

	if len(requestBodyID) > 0 {
		mongoClient := mongodb.GetClient()
		collection := mongoClient.Database(env.MongoKeys.Database).Collection("tasks")

		// filter := bson.D{{"fromContainerService", task.Service}, {"task", task.Task}, {"subtask", task.Subtask}, {"requestBodyId", task.RequestID}}
		filter := bson.M{
			"fromContainerService": bson.M{
				"$eq": task.Service,
			},
			"task": bson.M{
				"$eq": task.Task,
			},
			"subtask": bson.M{
				"$eq": task.Subtask,
			},
			"requestBodyId": bson.M{
				"$eq": requestBodyID,
			},
		}

		err := collection.FindOne(context.TODO(), filter).Decode(&existingTask)
		if err != nil {
			logs.StatusFileMessageLogging("FAILURE", filename, "getExistingTask", err.Error())
		}

		fmt.Printf("\n \n---------> existingTask  : %+v \n", existingTask)
	}

	return existingTask
}

func getExistingParsedTask(mongoDBTask interfaces.TaskModelInterface) interfaces.InitialisedRecordInfoInterface {
	toResponseID := mongoDBTask.ToResponseBodyID
	parsedTask := interfaces.InitialisedRecordInfoInterface{}

	if len(toResponseID) == 0 {
		return parsedTask
	}

	response := interfaces.ResponseModelInterface{}

	mongoClient := mongodb.GetClient()
	collection := mongoClient.Database(env.MongoKeys.Database).Collection("responses")

	docID, err := primitive.ObjectIDFromHex(toResponseID)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getExistingParsedTask", err.Error())
	}

	// filter := bson.D{{"_id", docID}}
	filter := bson.M{
		"_id": bson.M{
			"$eq": docID,
		},
	}
	err = collection.FindOne(context.TODO(), filter).Decode(&response)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getExistingParsedTask", err.Error())
	}

	parsedTask = interfaces.InitialisedRecordInfoInterface{
		ContainerID:             mongoDBTask.FromContainerID,
		ContainerService:        mongoDBTask.FromContainerService,
		RecordID:                mongoDBTask.ID.Hex(),
		Task:                    mongoDBTask.Task,
		Subtask:                 mongoDBTask.Subtask,
		ServiceContainerID:      mongoDBTask.ServiceContainerID,
		ServiceContainerService: mongoDBTask.ServiceContainerService,
		ResponseBody:            response.Response,
	}

	return parsedTask
}

func getNewParsedTask(mongoDBTask interfaces.TaskModelInterface, selectedContainerInfo interfaces.ContainerInfoStruct) interfaces.InitialisedRecordInfoInterface {
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

	mongoClient := mongodb.GetClient()
	collection := mongoClient.Database(env.MongoKeys.Database).Collection("responses")

	// filter := bson.D{{"id", mongoDBTask.ToResponseBodyID}}
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

	return parsedTask
}

func saveNewRequestAndGetID(requestBody string) string {

	mongoClient := mongodb.GetClient()

	collection := mongoClient.Database(env.MongoKeys.Database).Collection("requests")

	newRequest := interfaces.RequestModelInterface{Request: requestBody}

	newRequest.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.TODO(), newRequest)

	var finalID string

	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "saveNewRequestAndGetID", err.Error())
		return finalID
	}

	// newRequest.ID.Hex() is the way to access the string value
	// How to find a document using its id
	// https://kb.objectrocket.com/mongo-db/how-to-find-a-mongodb-document-by-its-bson-objectid-using-golang-452

	return newRequest.ID.Hex()

}

func recordNewInitialisedTaskWithRequestID(funcTask interfaces.ReceivedEventInterface, requestBodyId string) interfaces.InitialisedRecordInfoInterface {
	// TODO - Find a way to read json and put it in docker api
	selectedContainer := dockerapi.FetchConsumingContainer("backend")
	myContainerInfo := dockerapi.GetMyOfflineContainerInfo()

	newTask := interfaces.TaskModelInterface{
		FromRequestID:           funcTask.RequestID,
		FromContainerID:         funcTask.ContainerID,
		FromContainerService:    funcTask.Service,
		FromReceivedTime:        time.Now(),
		Task:                    funcTask.Task,
		Subtask:                 funcTask.Subtask,
		RequestBodyID:           requestBodyId,
		ToContainerID:           selectedContainer.ID,
		ToContainerService:      selectedContainer.Service,
		ServiceContainerID:      myContainerInfo.ID,
		ServiceContainerService: myContainerInfo.Service,
	}

	mongoClient := mongodb.GetClient()
	collection := mongoClient.Database(env.MongoKeys.Database).Collection("tasks")

	newTask.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.TODO(), newTask)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "recordNewInitialisedTaskWithRequestID", err.Error())
	}

	parsedTask := getNewParsedTask(newTask, selectedContainer)
	return parsedTask
}

func recordNewTaskAndRequest(task interfaces.ReceivedEventInterface) interfaces.InitialisedRecordInfoInterface {
	requestBodyID := saveNewRequestAndGetID(task.RequestBody)

	initialisedInfo := recordNewInitialisedTaskWithRequestID(task, requestBodyID)

	initialisedInfo.RequestBody = task.RequestBody

	return initialisedInfo
}

func getParsedResponse(funcResponse interfaces.ReceivedEventInterface, oldTask interfaces.TaskModelInterface) interfaces.EventInterface {
	response := interfaces.EventInterface{
		RequestID:    oldTask.FromRequestID,
		ContainerID:  oldTask.FromContainerID,
		Service:      oldTask.FromContainerService,
		ResponseBody: funcResponse.ResponseBody,
	}

	return response
}

func saveNewResponseAndGetID(funcResponse interfaces.ReceivedEventInterface) string {

	mongoClient := mongodb.GetClient()
	collection := mongoClient.Database(env.MongoKeys.Database).Collection("responses")

	newResponse := interfaces.ResponseModelInterface{Response: funcResponse.ResponseBody}
	newResponse.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.TODO(), newResponse)

	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "saveNewResponseAndGetID", err.Error())
	}

	return newResponse.ID.Hex()
}

func completeRecordInDB(funcResponse interfaces.ReceivedEventInterface, receivedTime time.Time, responseBodyID string) {
	fromSentTime := time.Now()

	mongoClient := mongodb.GetClient()
	collection := mongoClient.Database(env.MongoKeys.Database).Collection("tasks")

	docID, err := primitive.ObjectIDFromHex(funcResponse.RecordID)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "getExistingParsedTask", err.Error())
	}

	// filter := bson.D{{"_id", docID}}
	filter := bson.M{
		"_id": bson.M{
			"$eq": docID,
		},
	}

	// update := bson.D{
	// 	{"toReceivedTime", receivedTime},
	// 	{"toResponseBodyId", responseBodyID},
	// 	{"fromSentTime", fromSentTime}}

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

// RecordNewTaskInDB -
func RecordNewTaskInDB(task interfaces.ReceivedEventInterface) interfaces.InitialisedRecordInfoInterface {

	fmt.Printf("\n \n--------->  task  : %+v \n", task)
	existingTask := getExistingTask(task)

	if len(existingTask.FromRequestID) > 0 {
		parsedTask := getExistingParsedTask(existingTask)
		parsedTask.Existing = true

		return parsedTask
	}

	initRecordInfo := recordNewTaskAndRequest(task)

	return initRecordInfo
}

// CompleteExistingTaskRecordInDB -
func CompleteExistingTaskRecordInDB(funcResponse interfaces.ReceivedEventInterface) (interfaces.EventInterface, error) {

	preexistingResponse := interfaces.ResponseModelInterface{}

	mongoClient := mongodb.GetClient()
	collection := mongoClient.Database(env.MongoKeys.Database).Collection("responses")

	// filter := bson.D{{"response", funcResponse.ResponseBody}}
	filter := bson.M{
		"response": bson.M{
			"$eq": funcResponse.ResponseBody,
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

	task := interfaces.TaskModelInterface{}

	collection = mongoClient.Database(env.MongoKeys.Database).Collection("tasks")

	docID, err := primitive.ObjectIDFromHex(funcResponse.RecordID)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "CompleteExistingTaskRecordInDB", err.Error())
	}

	// filter = bson.D{{"_id", docID}}
	filter = bson.M{
		"_id": bson.M{
			"$eq": docID,
		},
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "CompleteExistingTaskRecordInDB", err.Error())
	}

	response := getParsedResponse(funcResponse, task)
	return response, nil
}
