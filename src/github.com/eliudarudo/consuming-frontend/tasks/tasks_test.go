package tasks

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/eliudarudo/consuming-frontend/interfaces"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

/*

   Note - Test a function and
          put all tests that we expect under it

*/

func TestDetermineTask(t *testing.T) {
	numbersRequestBody := map[string]interface{}{"a1": 3, "a2": 4}
	stringRequestBody := map[string]interface{}{"a1": "a", "a2": "b"}
	invalidRequestBody := map[string]interface{}{}

	jsonifiedNumbersRequestBodyString, _ := json.Marshal(numbersRequestBody)
	jsonifiedStringRequestBodyString, _ := json.Marshal(stringRequestBody)
	jsonifiedInvalidRequestBodyString, _ := json.Marshal(invalidRequestBody)

	numbersTaskType := determineTask(numbersRequestBody)
	stringTaskType := determineTask(stringRequestBody)
	invalidTaskType := determineTask(invalidRequestBody)

	stringifiedNumbersType := fmt.Sprintf("%v", numbersTaskType)
	stringifiedStringType := fmt.Sprintf("%v", stringTaskType)
	stringifiedEmptyType := fmt.Sprintf("%v", invalidTaskType)

	stringifiedExpectedNumberType := fmt.Sprintf("%v", interfaces.NUMBER)
	stringifiedExpectedStringType := fmt.Sprintf("%v", interfaces.STRING)

	t.Logf("\tGiven requestBody = %v", string(jsonifiedNumbersRequestBodyString))

	t.Logf("\t\tTest: \tExpected interfaces.TaskType = %v\n", stringifiedExpectedNumberType)
	if numbersTaskType == interfaces.NUMBER {
		t.Logf("\t\t%v Got : %v", succeedIcon, stringifiedNumbersType)
	} else {
		t.Errorf("\t\t%v Got : %v", failIcon, stringifiedNumbersType)
	}

	t.Logf("\tGiven requestBody = %v", string(jsonifiedStringRequestBodyString))

	t.Logf("\t\tTest: \tExpected interfaces.TaskType = %v\n", stringifiedExpectedStringType)
	if stringTaskType == interfaces.STRING {
		t.Logf("\t\t%v Got : %v", succeedIcon, stringifiedStringType)
	} else {
		t.Errorf("\t\t%v Got : %v", failIcon, stringifiedStringType)
	}

	t.Logf("\tGiven invalid requestBody = %v", string(jsonifiedInvalidRequestBodyString))

	t.Log("\t\tTest: \tExpected nothing")
	if stringifiedEmptyType == "" {
		t.Logf("\t\t%v Got : nothing", succeedIcon)
	} else {
		t.Errorf("\t\t%v Got : %v", failIcon, stringifiedEmptyType)
	}
}

func TestDetermineSubtask(t *testing.T) {
	numbersAddRequestBody := map[string]interface{}{"a1": 3, "a2": 4}
	numbersSubtractRequestBody := map[string]interface{}{"s1": 3, "s2": 4}
	numbersMultiplyRequestBody := map[string]interface{}{"m1": 3, "m2": 4}
	numbersDivideRequestBody := map[string]interface{}{"d1": 3, "d2": 4}
	invalidRequestBody := map[string]interface{}{}

	jsonifiedNumbersAddRequestBodyString, _ := json.Marshal(numbersAddRequestBody)
	jsonifiedNumbersSubtractRequestBodyString, _ := json.Marshal(numbersSubtractRequestBody)
	jsonifiedNumbersMultiplyRequestBodyString, _ := json.Marshal(numbersMultiplyRequestBody)
	jsonifiedNumbersDivideRequestBodyString, _ := json.Marshal(numbersDivideRequestBody)
	jsonifiedInvalidRequestBodyString, _ := json.Marshal(invalidRequestBody)

	numbersAddTaskType := determineTask(numbersAddRequestBody)
	numbersSubtractTaskType := determineTask(numbersSubtractRequestBody)
	numbersMultiplyTaskType := determineTask(numbersMultiplyRequestBody)
	numbersDivideTaskType := determineTask(numbersDivideRequestBody)
	invalidTaskType := determineTask(invalidRequestBody)

	numbersAddSubTaskType := determineSubtask(numbersAddTaskType, numbersAddRequestBody)
	numbersSubtractSubTaskType := determineSubtask(numbersSubtractTaskType, numbersSubtractRequestBody)
	numbersMultiplySubTaskType := determineSubtask(numbersMultiplyTaskType, numbersMultiplyRequestBody)
	numbersDivideSubTaskType := determineSubtask(numbersDivideTaskType, numbersDivideRequestBody)
	invalidSubTaskType := determineSubtask(invalidTaskType, invalidRequestBody)

	stringifiedAddNumberType := fmt.Sprintf("%v", numbersAddTaskType)
	stringifiedSubtractNumberType := fmt.Sprintf("%v", numbersSubtractTaskType)
	stringifiedMultiplyNumberType := fmt.Sprintf("%v", numbersMultiplyTaskType)
	stringifiedDivideNumberType := fmt.Sprintf("%v", numbersDivideTaskType)
	stringifiedInvalidType := fmt.Sprintf("%v", invalidTaskType)

	stringifiedNumbersAddSubType := fmt.Sprintf("%v", numbersAddSubTaskType)
	stringifiedNumbersSubtractSubType := fmt.Sprintf("%v", numbersSubtractSubTaskType)
	stringifiedNumbersMultiplySubType := fmt.Sprintf("%v", numbersMultiplySubTaskType)
	stringifiedNumbersDivideSubType := fmt.Sprintf("%v", numbersDivideSubTaskType)
	stringifiedInvalidSubType := fmt.Sprintf("%v", invalidSubTaskType)

	stringifiedExpectedAddNumberType := fmt.Sprintf("%v", interfaces.NUMBER)
	stringifiedExpectedSubtractNumberType := fmt.Sprintf("%v", interfaces.NUMBER)
	stringifiedExpectedMultiplyNumberType := fmt.Sprintf("%v", interfaces.NUMBER)
	stringifiedExpectedDivideNumberType := fmt.Sprintf("%v", interfaces.NUMBER)

	stringifiedExpectedAddNumbersSubType := fmt.Sprintf("%v", interfaces.ADD)
	stringifiedExpectedSubtractNumbersSubType := fmt.Sprintf("%v", interfaces.SUBTRACT)
	stringifiedExpectedMultiplyNumbersSubType := fmt.Sprintf("%v", interfaces.MULTIPLY)
	stringifiedExpectedDivideNumbersSubType := fmt.Sprintf("%v", interfaces.DIVIDE)

	t.Logf("\tGiven requestBody = %v", string(jsonifiedNumbersAddRequestBodyString))

	t.Logf("\t\tTest: \tExpected interfaces.TaskType = %v and interfaces.SubTaskType = %v\n", stringifiedExpectedAddNumberType, stringifiedExpectedAddNumbersSubType)
	if numbersAddTaskType == interfaces.NUMBER && numbersAddSubTaskType == interfaces.ADD {
		t.Logf("\t\t%v Got : %v and %v", succeedIcon, stringifiedAddNumberType, stringifiedNumbersAddSubType)
	} else {
		t.Errorf("\t\t%v Got : '%v' and '%v'", failIcon, stringifiedAddNumberType, stringifiedNumbersAddSubType)
	}

	t.Logf("\tGiven requestBody = %v", string(jsonifiedNumbersSubtractRequestBodyString))

	t.Logf("\t\tTest: \tExpected interfaces.TaskType = %v and interfaces.SubTaskType = %v\n", stringifiedExpectedSubtractNumberType, stringifiedExpectedSubtractNumbersSubType)
	if numbersSubtractTaskType == interfaces.NUMBER && numbersSubtractSubTaskType == interfaces.SUBTRACT {
		t.Logf("\t\t%v Got : %v and %v", succeedIcon, stringifiedSubtractNumberType, stringifiedNumbersSubtractSubType)
	} else {
		t.Errorf("\t\t%v Got : '%v' and '%v'", failIcon, stringifiedSubtractNumberType, stringifiedNumbersSubtractSubType)
	}

	t.Logf("\tGiven requestBody = %v", string(jsonifiedNumbersMultiplyRequestBodyString))

	t.Logf("\t\tTest: \tExpected interfaces.TaskType = %v and interfaces.SubTaskType = %v\n", stringifiedExpectedMultiplyNumberType, stringifiedExpectedMultiplyNumbersSubType)
	if numbersMultiplyTaskType == interfaces.NUMBER && numbersMultiplySubTaskType == interfaces.MULTIPLY {
		t.Logf("\t\t%v Got : %v and %v", succeedIcon, stringifiedMultiplyNumberType, stringifiedNumbersMultiplySubType)
	} else {
		t.Errorf("\t\t%v Got : '%v' and '%v'", failIcon, stringifiedMultiplyNumberType, stringifiedNumbersMultiplySubType)
	}

	t.Logf("\tGiven requestBody = %v", string(jsonifiedNumbersDivideRequestBodyString))

	t.Logf("\t\tTest: \tExpected interfaces.TaskType = %v and interfaces.SubTaskType = %v\n", stringifiedExpectedDivideNumberType, stringifiedExpectedDivideNumbersSubType)
	if numbersDivideTaskType == interfaces.NUMBER && numbersDivideSubTaskType == interfaces.DIVIDE {
		t.Logf("\t\t%v Got : %v and %v", succeedIcon, stringifiedDivideNumberType, stringifiedNumbersDivideSubType)
	} else {
		t.Errorf("\t\t%v Got : '%v' and '%v'", failIcon, stringifiedDivideNumberType, stringifiedNumbersDivideSubType)
	}

	t.Logf("\tGiven requestBody = %v", string(jsonifiedInvalidRequestBodyString))

	t.Logf("\t\tTest: \tExpected interfaces.TaskType = '%v' and interfaces.SubTaskType = '%v'\n", "", "")
	if invalidTaskType == "" && invalidSubTaskType == "" {
		t.Logf("\t\t%v Got : '%v' and '%v'", succeedIcon, stringifiedInvalidType, stringifiedInvalidSubType)
	} else {
		t.Errorf("\t\t%v Got : '%v' and '%v'", failIcon, stringifiedInvalidType, stringifiedInvalidSubType)
	}
}

func TestGetTargetService(t *testing.T) {
	numberTaskKey := fmt.Sprintf("%v", interfaces.NUMBER)
	stringTaskKey := fmt.Sprintf("%v", interfaces.STRING)

	expectedNumberTaskTargetService := "event"
	expectedStringTaskTargetService := "event"

	targetNumberService, err := getTargetService(numberTaskKey)
	if err != nil {
		t.Errorf("\t%v Failed to open 'task-maps.json' with error: %v", failIcon, err.Error())
	}

	targetStringService, err := getTargetService(stringTaskKey)
	if err != nil {
		t.Errorf("\t%v Failed to open 'task-maps.json' with error: %v", failIcon, err.Error())
	}

	invalidTargetService, err := getTargetService("")
	if err != nil {
		t.Errorf("\t%v Failed to open 'task-maps.json' with error: %v", failIcon, err.Error())
	}

	t.Logf("\tGiven Task being = %v", numberTaskKey)

	t.Logf("\t\tTest: \tExpected target service = '%v'", expectedNumberTaskTargetService)
	if targetNumberService == expectedNumberTaskTargetService {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, targetNumberService)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, targetNumberService)
	}

	t.Logf("\tGiven Task being = %v", stringTaskKey)

	t.Logf("\t\tTest: \tExpected target service = '%v'", expectedStringTaskTargetService)
	if targetStringService == expectedStringTaskTargetService {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, targetStringService)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, targetStringService)
	}

	t.Log("\tGiven no valid Task")

	t.Log("\t\tTest: \tExpected target service = ''")
	if invalidTargetService == "" {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, invalidTargetService)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, invalidTargetService)
	}
}
