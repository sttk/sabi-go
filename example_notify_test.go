package sabi_test

import (
	"fmt"
	"github.com/sttk-go/sabi"
	"time"
)

func ExampleAddAsyncErrHandler() {
	sabi.AddAsyncErrHandler(func(err sabi.Err, tm time.Time) {
		fmt.Println("Asynchronous error handling: " + err.Error())
	})
	sabi.FixErrCfgs()

	type FailToDoSomething struct{ Name string }

	sabi.ErrBy(FailToDoSomething{Name: "abc"})

	// Output:
	// Asynchronous error handling: {reason=FailToDoSomething, Name=abc}

	time.Sleep(100 * time.Millisecond)
	sabi.ClearErrHandlers()
}

func ExampleAddSyncErrHandler() {
	sabi.AddSyncErrHandler(func(err sabi.Err, tm time.Time) {
		fmt.Println("Synchronous error handling: " + err.Error())
	})
	sabi.FixErrCfgs()

	type FailToDoSomething struct{ Name string }

	sabi.ErrBy(FailToDoSomething{Name: "abc"})

	// Output:
	// Synchronous error handling: {reason=FailToDoSomething, Name=abc}

	sabi.ClearErrHandlers()
}

func ExampleFixErrCfgs() {
	sabi.AddSyncErrHandler(func(err sabi.Err, tm time.Time) {
		fmt.Println("This handler is registered")
	})

	sabi.FixErrCfgs()

	sabi.AddSyncErrHandler(func(err sabi.Err, tm time.Time) { // Bad example
		fmt.Println("This handler is not registered")
	})

	type FailToDoSomething struct{ Name string }

	sabi.ErrBy(FailToDoSomething{Name: "abc"})

	// Output:
	// This handler is registered

	sabi.ClearErrHandlers()
}
