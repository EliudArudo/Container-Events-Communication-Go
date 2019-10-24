package logs

import "fmt"

// StatusFileMessageLogging creates a standard format for logging using a
// status - string "SUCCESS / "FAILURE"
// filename - where function lives
// functionName - from where the StatusFileMessageLogging is called
// message - we want logged
func StatusFileMessageLogging(status string, file string, functionName string, message string) {
	fmt.Printf("%v : %v : %v : %v \n", status, file, functionName, message)
}
