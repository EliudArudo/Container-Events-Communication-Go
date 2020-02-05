package databaseops

import (
	"testing"
)

const succeedIcon = "\u2713"
const failIcon = "\u2717"

func TestGetTargetService(t *testing.T) {
	expectedNumberService := "backend"
	expectedStringService := "backend"

	numberTask := "NUMBER"
	stringTask := "STRING"

	gotNumberService, _ := getTargetService(numberTask)
	gotStringService, _ := getTargetService(stringTask)

	t.Logf("\tGiven a %v task", numberTask)

	t.Logf("\t\tTest: \tExpected targetService  = '%v'", expectedNumberService)

	if gotNumberService == expectedNumberService {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, gotNumberService)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, gotNumberService)
	}

	t.Logf("\tGiven a %v task", stringTask)

	t.Logf("\t\tTest: \tExpected targetService  = '%v'", expectedStringService)

	if gotStringService == expectedStringService {
		t.Logf("\t\t%v Got : '%v'", succeedIcon, gotStringService)
	} else {
		t.Errorf("\t\t%v Got : '%v'", failIcon, gotStringService)
	}

}

func TestGetParsedResponse(t *testing.T) {

}
