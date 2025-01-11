// Copyright (C) 2022-2025 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"path/filepath"
	"runtime"
	"time"
)

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
