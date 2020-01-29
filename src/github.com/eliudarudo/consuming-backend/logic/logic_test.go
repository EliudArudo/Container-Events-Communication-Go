package logic

import (
	"testing"

	"github.com/eliudarudo/consuming-backend/interfaces"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

var dummyRequestID = "dummyRequestID"

var dummyEvent = &interfaces.ReceivedEventInterface{}

func TestPerformLogic(t *testing.T) {
	stringsAdd := "{ \"a1\": \"a\", \"a2\": \"b\"}"
	expectedStringsAdd := "ab"

	dummyEvent.Task = "STRING"
	dummyEvent.Subtask = "ADD"

	dummyEvent.RequestBody = stringsAdd

	result := performLogic(dummyEvent)

	t.Logf("\tGiven dummyEvent.RequestBody = '%v', dummyEvent.Task = '%v' and dummyEvent.Subtask = '%v'", dummyEvent.RequestBody, dummyEvent.Task, dummyEvent.Subtask)

	t.Log("\t\tTest: \tExpected result to be 'ab'")

	if *result == expectedStringsAdd {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, *result)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, *result)
	}

	numbersAdd := "{ \"a1\": 3, \"a2\": 5 }"
	expectedNumbersAdd := "8"

	dummyEvent.Task = "NUMBER"
	dummyEvent.Subtask = "ADD"

	dummyEvent.RequestBody = numbersAdd

	result = performLogic(dummyEvent)

	t.Logf("\tGiven dummyEvent.RequestBody = '%v', dummyEvent.Task = '%v' and dummyEvent.Subtask = '%v'", dummyEvent.RequestBody, dummyEvent.Task, dummyEvent.Subtask)

	t.Logf("\t\tTest: \tExpected result to be '%v'", expectedNumbersAdd)

	if *result == expectedNumbersAdd {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, *result)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, *result)
	}

	numbersSubtract := "{ \"s1\": 3, \"s2\": 5 }"
	expectedNumbersSubtract := "-2"

	dummyEvent.Task = "NUMBER"
	dummyEvent.Subtask = "SUBTRACT"

	dummyEvent.RequestBody = numbersSubtract

	result = performLogic(dummyEvent)

	t.Logf("\tGiven dummyEvent.RequestBody = '%v', dummyEvent.Task = '%v' and dummyEvent.Subtask = '%v'", dummyEvent.RequestBody, dummyEvent.Task, dummyEvent.Subtask)

	t.Logf("\t\tTest: \tExpected result to be '%v'", expectedNumbersSubtract)

	if *result == expectedNumbersSubtract {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, *result)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, *result)
	}

	numbersMultiply := "{ \"s1\": 3, \"s2\": 5 }"
	expectedNumbersMultiply := "15"

	dummyEvent.Task = "NUMBER"
	dummyEvent.Subtask = "MULTIPLY"

	dummyEvent.RequestBody = numbersMultiply

	result = performLogic(dummyEvent)

	t.Logf("\tGiven dummyEvent.RequestBody = '%v', dummyEvent.Task = '%v' and dummyEvent.Subtask = '%v'", dummyEvent.RequestBody, dummyEvent.Task, dummyEvent.Subtask)

	t.Logf("\t\tTest: \tExpected result to be '%v'", expectedNumbersMultiply)

	if *result == expectedNumbersMultiply {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, *result)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, *result)
	}

	numbersDivide := "{ \"d1\": 3, \"d2\": 5 }"
	expectedNumbersDivide := "0.6"

	dummyEvent.Task = "NUMBER"
	dummyEvent.Subtask = "DIVIDE"

	dummyEvent.RequestBody = numbersDivide

	result = performLogic(dummyEvent)

	t.Logf("\tGiven dummyEvent.RequestBody = '%v', dummyEvent.Task = '%v' and dummyEvent.Subtask = '%v'", dummyEvent.RequestBody, dummyEvent.Task, dummyEvent.Subtask)

	t.Logf("\t\tTest: \tExpected result to be '%v'", expectedNumbersDivide)

	if *result == expectedNumbersDivide {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, *result)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, *result)
	}

}
