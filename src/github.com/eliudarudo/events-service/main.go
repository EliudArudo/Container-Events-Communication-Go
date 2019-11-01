package main

import (
	"runtime"

	"github.com/eliudarudo/event-service/initialise"
)

func init() {

	// Allocate one logical processor for the scheduler to use.
	runtime.GOMAXPROCS(1)
}

func main() {
	initialise.Go()
}
