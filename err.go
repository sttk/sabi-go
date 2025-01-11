// Copyright (C) 2022-2025 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// Err is the struct that represents an error used commonly across the
// Sabi Framework.
//
// It encapsulates the reason for the error, which can be any data type.
// Typically, the reason is an instance of a struct, which makes it easy
// to uniquely identify the error kind and location in the source code.
// In addition, since a struct can store additional informations as their
// fields, it is possible to provide more detailed information about the
// error.
//
// The reason for the error can be distinguished with a switch statement
// and type assertion, so it is easy to handle the error in a type-safe
// manner.
//
// This struct also contains an optional cause error, which is the error
// that caused the current error. This is useful for chaining errors.
//
// This struct implements the error interface, so it can be used as an
// error object in Go programs.
// And since this struct implements the Unwrap method, it can be used as
// a wrapper error object in Go programs.
type Err struct {
	reason any
	cause  error
}

// Ok returns an instance of Err with no reason, indicating no error.
// This is the default "no error" state.
func Ok() Err {
	return Err{}
}

// NewErr creates a new Err instance with the provided reason.
// Optionally, a cause can also be supplied, which represents a lower-level
// error.
func NewErr(reason any, cause ...error) Err {
	var e Err
	e.reason = reason

	if len(cause) > 0 {
		e.cause = cause[0]
	}

	notifyErr(e)

	return e
}

// Reason returns the reason for the error, which can be any type.
// This helps in analyzing why the error occurred.
func (e Err) Reason() any {
	return e.reason
}

// Error returns a string representation of the Err instance.
// It formats the error, including the package path, reason, and cause.
func (e Err) Error() string {
	var sb strings.Builder

	t := reflect.TypeOf(e)
	s := t.PkgPath()
	if len(s) > 0 {
		sb.WriteString(s)
		sb.WriteByte('.')
	}
	sb.WriteString(t.Name())

	sb.WriteString(" { reason = ")

	if e.reason == nil {
		sb.WriteString("nil")
	} else {
		v := reflect.ValueOf(e.reason)

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			if v.CanInterface() {
				sb.WriteString(fmt.Sprintf("%v", v.Interface()))
			}
		} else {
			t := v.Type()

			s := t.PkgPath()
			if len(s) > 0 {
				sb.WriteString(s)
				sb.WriteByte('.')
			}
			sb.WriteString(t.Name())

			n := v.NumField()

			if n > 0 {
				sb.WriteString(" { ")

				for i := 0; i < n; i++ {
					if i > 0 {
						sb.WriteString(", ")
					}

					k := t.Field(i).Name

					f := v.Field(i)
					if f.CanInterface() { // false if the field is not public
						sb.WriteString(k)
						sb.WriteString(": ")
						sb.WriteString(fmt.Sprintf("%v", f.Interface()))
					}
				}

				sb.WriteString(" }")
			}
		}
	}

	if e.cause != nil {
		sb.WriteString(", cause = ")
		sb.WriteString(e.cause.Error())
	}

	sb.WriteString(" }")

	return sb.String()
}

// Unwrap returns the underlying cause of the error, allowing it to be chained.
// This helps in accessing the root cause when errors are wrapped.
func (e Err) Unwrap() error {
	return e.cause
}

// Cause returns the cause of the error.
// This is similar to Unwrap but provides a direct access method.
func (e Err) Cause() error {
	return e.cause
}

// IsOk returns true if the Err instance has no reason, indicating no error.
// This is used to check if the operation was successful.
func (e Err) IsOk() bool {
	return (e.reason == nil)
}

// IsNotOk returns true if the Err instance has a reason, indicating an error
// occurred.
// This is the inverse of IsOk, used to determine if an error is present.
func (e Err) IsNotOk() bool {
	return (e.reason != nil)
}

// IfOkThen executes the provided function if no error is present (IsOk).
// This is useful for chaining operations that only proceed if no error has
// occurred.
func (e Err) IfOkThen(fn func() Err) Err {
	if e.IsOk() {
		return fn()
	}
	return e
}

// ErrOccasion represents the details of an error occurrence,
// including the time it occurred, and the source file and line number.
type ErrOccasion struct {
	time time.Time
	file string
	line int
}

// Time returns the time when the error occurred.
func (e ErrOccasion) Time() time.Time {
	return e.time
}

// File returns the name of the source file where the error occurred.
func (e ErrOccasion) File() string {
	return e.file
}

// Line returns the line number in the source file where the error occurred.
func (e ErrOccasion) Line() int {
	return e.line
}

type errHandlerListItem struct {
	handler func(Err, ErrOccasion)
	next    *errHandlerListItem
}

type errHandlerList struct {
	head *errHandlerListItem
	last *errHandlerListItem
}

var (
	syncErrHandlers    = errHandlerList{nil, nil}
	asyncErrHandlers   = errHandlerList{nil, nil}
	isErrHandlersFixed = false
)

// AddSyncErrHandler adds a new synchronous error handler to the global handler list.
// It will not add the handler if the handlers have been fixed using FixErrHandlers.
func AddSyncErrHandler(handler func(Err, ErrOccasion)) {
	if isErrHandlersFixed {
		return
	}

	last := syncErrHandlers.last
	syncErrHandlers.last = &errHandlerListItem{handler, nil}

	if last != nil {
		last.next = syncErrHandlers.last
	}

	if syncErrHandlers.head == nil {
		syncErrHandlers.head = syncErrHandlers.last
	}
}

// AddAsyncErrHandler adds a new asynchronous error handler to the global handler list.
// It will not add the handler if the handlers have been fixed using FixErrHandlers.
func AddAsyncErrHandler(handler func(Err, ErrOccasion)) {
	if isErrHandlersFixed {
		return
	}

	last := asyncErrHandlers.last
	asyncErrHandlers.last = &errHandlerListItem{handler, nil}

	if last != nil {
		last.next = asyncErrHandlers.last
	}

	if asyncErrHandlers.head == nil {
		asyncErrHandlers.head = asyncErrHandlers.last
	}
}

// FixErrHandlers prevents further modification of the error handler lists.
// Before this is called, no Err is notified to the handlers.
// After this is called, no new handlers can be added, and Err(s) is notified to the
// handlers.
func FixErrHandlers() {
	isErrHandlersFixed = true
}

func notifyErr(err Err) {
	if !isErrHandlersFixed {
		return
	}

	if syncErrHandlers.head == nil && asyncErrHandlers.head == nil {
		return
	}

	var occ ErrOccasion
	occ.time = time.Now().UTC()

	_, file, line, ok := runtime.Caller(2)
	if ok {
		occ.file = filepath.Base(file)
		occ.line = line
	}

	for item := syncErrHandlers.head; item != nil; item = item.next {
		item.handler(err, occ)
	}

	for item := asyncErrHandlers.head; item != nil; item = item.next {
		go item.handler(err, occ)
	}
}
