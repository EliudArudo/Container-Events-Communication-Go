package logs

import "fmt"

// StatusFileMessageLogging simply logs an error or success case
func StatusFileMessageLogging(status string, file string, functionName string, message string) {
	fmt.Printf("%v : %v : %v : %v \n", status, file, functionName, message)
}
