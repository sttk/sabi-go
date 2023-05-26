// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"fmt"
	"reflect"
)

// Err is a struct which represents an error with a reason.
type Err struct {
	reason any
	cause  error
}

var ok = Err{}

// Ok is a function which returns an Err of which reason is nil.
func Ok() Err {
	return ok
}

// NewErr is a function which creates a new Err with a specified reason and
// an optional cause.
// A reason is a struct of which name expresses what is a reason.
func NewErr(reason any, cause ...error) Err {
	var err Err
	err.reason = reason

	if len(cause) > 0 {
		err.cause = cause[0]
	}

	notifyErr(err)

	return err
}

// IsOk method checks whether this Err indicates no error.
func (err Err) IsOk() bool {
	return (err.reason == nil)
}

// IsNotOk method checks whether this Err indicates an error.
func (err Err) IsNotOk() bool {
	return (err.reason != nil)
}

// Reason method returns an err reaason struct.
func (err Err) Reason() any {
	return err.reason
}

// ReasonName method returns a name of a reason struct type.
func (err Err) ReasonName() string {
	if err.reason == nil {
		return ""
	}
	t := reflect.TypeOf(err.reason)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// ReasonPackage method returns a package path of a reason struct type.
func (err Err) ReasonPackage() string {
	if err.reason == nil {
		return ""
	}
	t := reflect.TypeOf(err.reason)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath()
}

// Cause method returns a causal error of this Err.
func (err Err) Cause() error {
	return err.cause
}

// Error method returns a string which expresses this error.
func (err Err) Error() string {
	if err.reason == nil {
		return "{reason=nil}"
	}

	v := reflect.ValueOf(err.reason)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	s := "{reason=" + t.Name()

	n := v.NumField()
	for i := 0; i < n; i++ {
		k := t.Field(i).Name

		f := v.Field(i)
		if f.CanInterface() {
			s += ", " + k + "=" + fmt.Sprintf("%v", f.Interface())
		}
	}

	if err.cause != nil {
		s += ", cause=" + err.cause.Error()
	}

	s += "}"
	return s
}

// Unwrap method returns an error which is wrapped in this error.
func (err Err) Unwrap() error {
	return err.cause
}

// Get method returns a parameter value of a specified name, which is one of
// fields of the reason struct.
// If the specified named field is not found in this Err and this cause is
// also Err struct, this method digs hierarchically to find the field.
func (err Err) Get(name string) any {
	if err.reason == nil {
		return nil
	}

	v := reflect.ValueOf(err.reason)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	f := v.FieldByName(name)
	if f.IsValid() && f.CanInterface() {
		return f.Interface()
	}

	if err.cause != nil {
		t := reflect.TypeOf(err.cause)
		_, ok := t.MethodByName("Reason")
		if ok {
			_, ok := t.MethodByName("Get")
			if ok {
				return err.cause.(Err).Get(name)
			}
		}
	}

	return nil
}

// Situation method returns a map containing the field names and values of this
// reason struct and of this cause if it is also Err struct.
func (err Err) Situation() map[string]any {
	var m map[string]any

	if err.reason == nil {
		return m
	}

	v := reflect.ValueOf(err.reason)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if err.cause != nil {
		t := reflect.TypeOf(err.cause)
		_, ok := t.MethodByName("Reason")
		if ok {
			_, ok := t.MethodByName("Situation")
			if ok {
				m = err.cause.(Err).Situation()
			}
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
		if f.CanInterface() { // false if field is not public
			m[k] = f.Interface()
		}
	}

	return m
}

// IfOk method executes an argument function if this Err indicates non error.
// If this Err indicates some error, this method just returns this Err.
func (err Err) IfOk(fn func() Err) Err {
	if err.IsOk() {
		return fn()
	}
	return err
}
