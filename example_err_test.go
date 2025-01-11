package sabi_test

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sttk/sabi"
)

func ExampleNewErr() {
	type /* error reasons */ (
		FailToDoSomething  struct{}
		FailToDoWithParams struct {
			Param1 string
			Param2 int
		}
	)

	// (1) Creates an Err with no parameter.
	err := sabi.NewErr(FailToDoSomething{})
	fmt.Printf("(1) %v\n", err)

	// (2) Creates an Err with parameters.
	err = sabi.NewErr(FailToDoWithParams{
		Param1: "ABC",
		Param2: 123,
	})
	fmt.Printf("(2) %v\n", err)

	cause := errors.New("Causal error")

	// (3) Creates an Err with a causal error.
	err = sabi.NewErr(FailToDoSomething{}, cause)
	fmt.Printf("(3) %v\n", err)

	// (4) Creates an Err with parameters and a causal error.
	err = sabi.NewErr(FailToDoWithParams{
		Param1: "ABC",
		Param2: 123,
	}, cause)
	fmt.Printf("(4) %v\n", err)
	// Output:
	// (1) github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.FailToDoSomething }
	// (2) github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.FailToDoWithParams { Param1: ABC, Param2: 123 } }
	// (3) github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.FailToDoSomething, cause = Causal error }
	// (4) github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.FailToDoWithParams { Param1: ABC, Param2: 123 }, cause = Causal error }
}

func ExampleOk() {
	err := sabi.Ok()
	fmt.Printf("err = %v\n", err)
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())
	// Output:
	// err = github.com/sttk/sabi.Err { reason = nil }
	// err.IsOk() = true
}

func ExampleErr_Cause() {
	type FailToDoSomething struct{}

	cause := errors.New("Causal error")

	err := sabi.NewErr(FailToDoSomething{}, cause)
	fmt.Printf("%v\n", err.Cause())
	// Output:
	// Causal error
}

func ExampleErr_Error() {
	type FailToDoSomething struct {
		Param1 string
		Param2 int
	}

	cause := errors.New("Causal error")

	err := sabi.NewErr(FailToDoSomething{
		Param1: "ABC",
		Param2: 123,
	}, cause)
	fmt.Printf("%v\n", err.Error())
	// Output:
	// github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.FailToDoSomething { Param1: ABC, Param2: 123 }, cause = Causal error }
}

func ExampleErr_IsOk() {
	err := sabi.Ok()
	fmt.Printf("%v\n", err.IsOk())

	type FailToDoSomething struct{}
	err = sabi.NewErr(FailToDoSomething{})
	fmt.Printf("%v\n", err.IsOk())
	// Output:
	// true
	// false
}

func ExampleErr_IsNotOk() {
	err := sabi.Ok()
	fmt.Printf("%v\n", err.IsNotOk())

	type FailToDoSomething struct{}
	err = sabi.NewErr(FailToDoSomething{})
	fmt.Printf("%v\n", err.IsNotOk())
	// Output:
	// false
	// true
}

func ExampleErr_Reason() {
	type FailToDoSomething struct {
		Param1 string
	}

	err := sabi.NewErr(FailToDoSomething{Param1: "value1"})
	switch err.Reason().(type) {
	case FailToDoSomething:
		fmt.Println("The reason of the error is: FailToDoSomething")
		reason := err.Reason().(FailToDoSomething)
		fmt.Printf("The value of reason.Param1 is: %v\n", reason.Param1)
	}

	err = sabi.NewErr(&FailToDoSomething{Param1: "value2"})
	switch err.Reason().(type) {
	case *FailToDoSomething:
		fmt.Println("The reason of the error is: *FailToDoSomething")
		reason := err.Reason().(*FailToDoSomething)
		fmt.Printf("The value of reason.Param1 is: %v\n", reason.Param1)
	}
	// Output:
	// The reason of the error is: FailToDoSomething
	// The value of reason.Param1 is: value1
	// The reason of the error is: *FailToDoSomething
	// The value of reason.Param1 is: value2
}

func ExampleErr_Unwrap() {
	type FailToDoSomething struct{}

	cause1 := errors.New("Causal error 1")
	cause2 := errors.New("Causal error 2")

	err := sabi.NewErr(FailToDoSomething{}, cause1)

	fmt.Printf("err.Unwrap() = %v\n", err.Unwrap())
	fmt.Printf("errors.Unwrap(err) = %v\n", errors.Unwrap(err))
	fmt.Printf("errors.Is(err, cause1) = %v\n", errors.Is(err, cause1))
	fmt.Printf("errors.Is(err, cause2) = %v\n", errors.Is(err, cause2))
	// Output:
	// err.Unwrap() = Causal error 1
	// errors.Unwrap(err) = Causal error 1
	// errors.Is(err, cause1) = true
	// errors.Is(err, cause2) = false
}

func ExampleErr_IfOkThen() {
	type FailToDoSomething struct{}

	err := sabi.Ok()
	err.IfOkThen(func() sabi.Err {
		fmt.Println("execute if non error.")
		return sabi.Ok()
	})

	err = sabi.NewErr(FailToDoSomething{})
	err.IfOkThen(func() sabi.Err {
		fmt.Println("not execute if some error.")
		return sabi.Ok()
	})
	// Output:
	// execute if non error.
}

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
	// This handler is registered at example_err_test.go:222

	sabi.ClearErrHandlers()
}
