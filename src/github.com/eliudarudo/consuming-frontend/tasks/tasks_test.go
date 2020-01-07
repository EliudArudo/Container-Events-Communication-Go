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

	jsonifiedNumbersRequestBodyString, _ := json.Marshal(numbersRequestBody)
	jsonifiedStringRequestBodyString, _ := json.Marshal(stringRequestBody)

	numbersTaskType := determineTask(numbersRequestBody)
	stringTaskType := determineTask(stringRequestBody)

	stringifiedNumbersType := fmt.Sprintf("%v", numbersTaskType)
	stringifiedStringType := fmt.Sprintf("%v", stringTaskType)

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

}
