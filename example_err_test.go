package sabi_test

import (
	"errors"
	"fmt"
	"github.com/sttk-go/sabi"
)

func ExampleNewErr() {
	type /* error reason */ (
		FailToDoSomething           struct{}
		FailToDoSomethingWithParams struct {
			Param1 string
			Param2 int
		}
	)

	// (1) Creates an Err with no situation parameter.
	err := sabi.NewErr(FailToDoSomething{})
	fmt.Printf("(1) %v\n", err)

	// (2) Creates an Err with situation parameters.
	err = sabi.NewErr(FailToDoSomethingWithParams{
		Param1: "ABC",
		Param2: 123,
	})
	fmt.Printf("(2) %v\n", err)

	cause := errors.New("Causal error")

	// (3) Creates an Err with a causal error.
	err = sabi.NewErr(FailToDoSomething{}, cause)
	fmt.Printf("(3) %v\n", err)

	// (4) Creates an Err with situation parameters and a causal error.
	err = sabi.NewErr(FailToDoSomethingWithParams{
		Param1: "ABC",
		Param2: 123,
	}, cause)
	fmt.Printf("(4) %v\n", err)

	// Output:
	// (1) {reason=FailToDoSomething}
	// (2) {reason=FailToDoSomethingWithParams, Param1=ABC, Param2=123}
	// (3) {reason=FailToDoSomething, cause=Causal error}
	// (4) {reason=FailToDoSomethingWithParams, Param1=ABC, Param2=123, cause=Causal error}
}

func ExampleOk() {
	err := sabi.Ok()
	fmt.Printf("err = %v\n", err)
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	// Output:
	// err = {reason=nil}
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
	// {reason=FailToDoSomething, Param1=ABC, Param2=123, cause=Causal error}
}

func ExampleErr_Get() {
	type FailToDoSomething struct {
		Param1 string
		Param2 int
	}

	err := sabi.NewErr(FailToDoSomething{
		Param1: "ABC",
		Param2: 123,
	})
	fmt.Printf("Param1=%v\n", err.Get("Param1"))
	fmt.Printf("Param2=%v\n", err.Get("Param2"))
	fmt.Printf("Param3=%v\n", err.Get("Param3"))

	// Output:
	// Param1=ABC
	// Param2=123
	// Param3=<nil>
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

func ExampleErr_ReasonName() {
	type FailToDoSomething struct{}

	err := sabi.NewErr(FailToDoSomething{})
	fmt.Printf("%v\n", err.ReasonName())

	// Output:
	// FailToDoSomething
}

func ExampleErr_ReasonPackage() {
	type FailToDoSomething struct{}

	err := sabi.NewErr(FailToDoSomething{})
	fmt.Printf("%v\n", err.ReasonPackage())

	// Output:
	// github.com/sttk-go/sabi_test
}

func ExampleErr_Situation() {
	type FailToDoSomething struct {
		Param1 string
		Param2 int
	}

	err := sabi.NewErr(FailToDoSomething{
		Param1: "ABC",
		Param2: 123,
	})
	fmt.Printf("%v\n", err.Situation())

	// Output:
	// map[Param1:ABC Param2:123]
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
	fmt.Printf("errors.As(err, cause1) = %v\n", errors.Is(err, cause1))
	fmt.Printf("errors.As(err, cause2) = %v\n", errors.Is(err, cause2))

	// Output:
	// err.Unwrap() = Causal error 1
	// errors.Unwrap(err) = Causal error 1
	// errors.Is(err, cause1) = true
	// errors.Is(err, cause2) = false
	// errors.As(err, cause1) = true
	// errors.As(err, cause2) = false
}
