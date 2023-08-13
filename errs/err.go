// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

// errs is the package for error processing for sabi framework.
// This package provides Err struct instead of Go standard error.
// Err takes an any struct which indicates a reason of an error of creating it.
// And Err has some functionalities to make it easier to hold, get, and notify
// error informations.
//
// # Creating an Err
//
// In contract that Go standard error requires to implement a struct with Error
// method to output a string of error content, Err only requires to implement a
// struct with no method.
//
//	type FailToDoSomething struct { Name string }
//	err := errs.New(FailToDoSomething{Name: name})
//
// And unlike Go standard error of which value is nil for no error, Err has an
// instance for no error by Ok function.
//
//	err := errs.Ok()
//
// In addition, Err provides methods to check whether an Err instance indicates
// an error or not.
//
//	if err.IsOk() { ... }     // err indicates no error.
//	if err.IsNotOk() { ... }  // err indicates an error.
//
// Also, Err provides the method to do a next function if it is no error.
//
//	func doSomething() errs.Err { ... }
//	func doNextThing() errs.Err { ... }
//	return doSomething().IfOk(doNextThing)
//
// # Distinction of error kinds
//
// To distinct error kinds, a type switch statement can be applied to a type of
// Err's reason.
//
//	switch err.Reason().(type) {
//	case nil:
//		...
//	case FailToDoSomething:
//		reason, ok := err.Reason().(FailToDoSomething)
//		...
//	default:
//		...
//	}
//
// # Error notifications
//
// This package support notification of error creations.
// Notification handlers can be registered with AddSyncHandler and
// AddAsyncHandler.
// When a Err is created with New function, handlers receives an Err instance
// and a ErrOcc instance.
// ErrOcc is a struct having informations when and where the error is occured.
//
//	errs.AddSyncHandler(func(err errs.Err, occ errs.ErrOcc) {
//	    logger.Printf("%s (%s:%d) %v\n",
//	        occ.Time().Format("2006-01-02T15:04:05Z"),
//	        occ.File(), occ.Line(), err)
//	})
//	errs.FixCfg()
//
//	errs.New(FailToDoSomething{Name: name})
package errs

import (
	"fmt"
	"reflect"
)

// Err is a struct type which represents an error with a reason.
// This instance can has an instance or pointer of any struct type of which
// indicates a reason by which this error is caused.
// A reason can has some fields that helps to know error situation where this
// error is caused.
type Err struct {
	reason any
	cause  error
}

var ok = Err{}

// Ok is the function which returns an Err instance of which reason is nil.
func Ok() Err {
	return ok
}

// New is the function which creates a new Err with a specified reason and an
// optional cause.
// A reason is a struct type of which name expresses what is a reason.
func New(reason any, cause ...error) Err {
	var e Err
	e.reason = reason

	if len(cause) > 0 {
		e.cause = cause[0]
	}

	notifyErr(e)

	return e
}

// IsOk is the method that checks whether this Err indicates there is no error.
func (e Err) IsOk() bool {
	return (e.reason == nil)
}

// IsNotOk is the method that checks whether this Err indicates there is an error.
func (e Err) IsNotOk() bool {
	return (e.reason != nil)
}

// IfOk is the method that executes an argument function if this Err indicates there is no error.
// If this Err indicates there is an error, this method just returns this Err it self.
func (e Err) IfOk(fn func() Err) Err {
	if e.IsOk() {
		return fn()
	}
	return e
}

// Reason is the method to get the reason of this error.
func (e Err) Reason() any {
	return e.reason
}

// ReasonName is the method to get the name of this reason's struct type.
func (e Err) ReasonName() string {
	if e.reason == nil {
		return ""
	}
	t := reflect.TypeOf(e.reason)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// ReasonPackage is the method to get the package path of this reason's struct
// type.
func (e Err) ReasonPackage() string {
	if e.reason == nil {
		return ""
	}
	t := reflect.TypeOf(e.reason)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath()
}

// Error is the method to get a string that expresses the content of this
// error.
func (e Err) Error() string {
	if e.reason == nil {
		return "{reason=nil}"
	}

	v := reflect.ValueOf(e.reason)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	s := "{reason=" + t.Name()

	n := v.NumField()
	for i := 0; i < n; i++ {
		k := t.Field(i).Name

		f := v.Field(i)
		if f.CanInterface() { // false if the field is nor public
			s += ", " + k + "=" + fmt.Sprintf("%v", f.Interface())
		}
	}

	if e.cause != nil {
		s += ", cause=" + e.cause.Error()
	}

	s += "}"
	return s
}

// Unwrap is the method to get an error which is wrapped in this error.
func (e Err) Unwrap() error {
	return e.cause
}

// Cause is the method to get the causal error of this Err.
func (e Err) Cause() error {
	return e.cause
}

// Get is the method to get a field value of the reason struct type by the
// specified name.
// If the specified named field is not found in the reason of this Err, this
// method finds a same named field in reasons of cause errors hierarchically.
func (e Err) Get(name string) any {
	if e.reason == nil {
		return nil
	}

	v := reflect.ValueOf(e.reason)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	f := v.FieldByName(name)
	if f.IsValid() && f.CanInterface() {
		return f.Interface()
	}

	if e.cause != nil {
		c, ok := e.cause.(Err)
		if ok {
			return c.Get(name)
		}
	}

	return nil
}

// Situation is the method to get a map which contains parameters that
// represents error situation.
func (e Err) Situation() map[string]any {
	var m map[string]any

	if e.reason == nil {
		return m
	}

	v := reflect.ValueOf(e.reason)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if e.cause != nil {
		c, ok := e.cause.(Err)
		if ok {
			m = c.Situation()
		}
	}

	if m == nil {
		m = make(map[string]any)
	}

	t := v.Type()
	n := v.NumField()

	for i := 0; i < n; i++ {
		k := t.Field(i).Name
		f := v.Field(i)
		if f.CanInterface() { // false if the field is nor public
			m[k] = f.Interface()
		}
	}

	return m
}
