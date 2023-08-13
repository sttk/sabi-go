package errs_test

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sttk/sabi/errs"
)

func ExampleAddAsyncHandler() {
	errs.AddAsyncHandler(func(err errs.Err, occ errs.ErrOcc) {
		fmt.Println("Asynchronous error handling: " + err.Error())
	})
	errs.FixCfg()

	type FailToDoSomething struct{ Name string }

	errs.New(FailToDoSomething{Name: "abc"})
	// Output:
	// Asynchronous error handling: {reason=FailToDoSomething, Name=abc}

	time.Sleep(100 * time.Millisecond)
	errs.ClearErrHandlers()
}

func ExampleAddSyncHandler() {
	errs.AddSyncHandler(func(err errs.Err, occ errs.ErrOcc) {
		fmt.Println("Synchronous error handling: " + err.Error())
	})
	errs.FixCfg()

	type FailToDoSomething struct{ Name string }

	errs.New(FailToDoSomething{Name: "abc"})
	// Output:
	// Synchronous error handling: {reason=FailToDoSomething, Name=abc}

	errs.ClearErrHandlers()
}

func ExampleFixCfg() {
	errs.AddSyncHandler(func(err errs.Err, occ errs.ErrOcc) {
		fmt.Println("This handler is registered at " + occ.File() + ":" +
			strconv.Itoa(occ.Line()))
	})

	errs.FixCfg()

	errs.AddSyncHandler(func(err errs.Err, occ errs.ErrOcc) {
		fmt.Println("This handler is not registered")
	})

	type FailToDoSomething struct{ Name string }

	errs.New(FailToDoSomething{Name: "abc"})
	// Output:
	// This handler is registered at example_notify_test.go:56

	errs.ClearErrHandlers()
}
