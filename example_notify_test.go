package sabi_test

import (
	"fmt"
	"github.com/sttk-go/sabi"
	"strconv"
	"time"
)

func ExampleAddAsyncErrHandler() {
	sabi.AddAsyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {
		fmt.Println("Asynchronous error handling: " + err.Error())
	})
	sabi.FixErrCfgs()

	type FailToDoSomething struct{ Name string }

	sabi.NewErr(FailToDoSomething{Name: "abc"})

	// Output:
	// Asynchronous error handling: {reason=FailToDoSomething, Name=abc}

	time.Sleep(100 * time.Millisecond)
	sabi.ClearErrHandlers()
}

func ExampleAddSyncErrHandler() {
	sabi.AddSyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {
		fmt.Println("Synchronous error handling: " + err.Error())
	})
	sabi.FixErrCfgs()

	type FailToDoSomething struct{ Name string }

	sabi.NewErr(FailToDoSomething{Name: "abc"})

	// Output:
	// Synchronous error handling: {reason=FailToDoSomething, Name=abc}

	sabi.ClearErrHandlers()
}

func ExampleFixErrCfgs() {
	sabi.AddSyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {
		fmt.Println("This handler is registered at " + occ.File() + ":" +
			strconv.Itoa(occ.Line()))
	})

	sabi.FixErrCfgs()

	sabi.AddSyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {
		fmt.Println("This handler is not registered")
	})

	type FailToDoSomething struct{ Name string }

	sabi.NewErr(FailToDoSomething{Name: "abc"})

	// Output:
	// This handler is registered at example_notify_test.go:57

	sabi.ClearErrHandlers()
}
