// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
)

// Err is a structure type which represents an error with a reason.
type Err struct {
	reason any
	file   string
	line   int
	cause  error
}

// NoError is an error reason which indicates no error.
type NoError struct{}

var ok = Err{reason: NoError{}}

// Ok is a function which returns an Err value of which reason is NoError.
func Ok() Err {
	return ok
}

// NewErr is a function which creates a new Err value with a reason and a cause.
// A reason is a structure type of which name expresses what is a reason.
func NewErr(reason any, cause ...error) Err {
	var err Err
	err.reason = reason

	if len(cause) > 0 {
		err.cause = cause[0]
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		err.file = filepath.Base(file)
		err.line = line
	}

	notifyErr(err)

	return err
}

// IsOk method determines whether an Err indicates no error.
func (err Err) IsOk() bool {
	switch err.reason.(type) {
	case NoError, *NoError:
		return true
	default:
		return false
	}
}

// Reason method returns an error reason structure.
func (err Err) Reason() any {
	return err.reason
}

// ReasonName method returns a name of a reason structure type.
func (err Err) ReasonName() string {
	t := reflect.TypeOf(err.reason)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// ReasonPackage method returns a package path of a reason structure type.
func (err Err) ReasonPackage() string {
	t := reflect.TypeOf(err.reason)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath()
}

// FileName method returns a source file name where an Err was caused.
func (err Err) FileName() string {
	return err.file
}

// LineNumber method returns a line number in a source file where an Err was
// caused.
func (err Err) LineNumber() int {
	return err.line
}

// Cause method returns a causal error of an Err.
func (err Err) Cause() error {
	return err.cause
}

// Error method returns a string which expresses this error.
func (err Err) Error() string {
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
			s += fmt.Sprintf(", %s=%v", k, f.Interface())
		}
	}

	if err.cause != nil {
		s += ", cause=" + err.cause.Error()
	}

	s += "}"
	return s
}

// Unwrap method returns an error which is wrapped by this error.
func (err Err) Unwrap() error {
	return err.cause
}

// Get method returns a parameter value of a specified name, which is one of
// parameters which represents situation when an Err was caused.
// If a parameter is not found in an Err and its .cause is also an Err, this
// method digs hierarchically to find a parameter which has same name.
func (err Err) Get(name string) any {
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

// Situation method returns a map containing parameters which represent a
// situation when an Err was caused.
// If a .cause is an Err, a returned map includes parameters of .cause
// hierarchically.
func (err Err) Situation() map[string]any {
	v := reflect.ValueOf(err.reason)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var m map[string]any

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
