package sabi_test

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sttk/sabi"
)

func ExampleAddAsyncErrHandler() {
	sabi.AddAsyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {
		fmt.Println("Asynchronous error handling: " + err.Error())
	})
	sabi.FixErrHandlers()

	type FailToDoSomething struct{ Name string }

	sabi.NewErr(FailToDoSomething{Name: "abc"})
	// Output:
	// Asynchronous error handling: github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.FailToDoSomething { Name: abc } }

	time.Sleep(100 * time.Millisecond)
	sabi.ClearErrHandlers()
}

func ExampleAddSyncErrHandler() {
	sabi.AddSyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {
		fmt.Println("Synchronous error handling: " + err.Error())
	})
	sabi.FixErrHandlers()

	type FailToDoSomething struct{ Name string }

	sabi.NewErr(FailToDoSomething{Name: "abc"})
	// Output:
	// Synchronous error handling: github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.FailToDoSomething { Name: abc } }

	sabi.ClearErrHandlers()
}

func ExampleFixErrHandlers() {
	sabi.AddSyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {
		fmt.Println("This handler is registered at " + occ.File() + ":" + strconv.Itoa(occ.Line()))
	})

	type FailToDoSomething struct{ Name string }

	sabi.NewErr(FailToDoSomething{Name: "abc"})

	sabi.FixErrHandlers()

	sabi.AddSyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {
		fmt.Println("This handler is not registered")
	})

	sabi.NewErr(FailToDoSomething{Name: "def"})
	// Output:
	// This handler is registered at example_err-notify_test.go:57

	sabi.ClearErrHandlers()
}
