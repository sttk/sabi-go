// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package errs

import (
	"path/filepath"
	"runtime"
	"time"
)

// ErrOcc is the struct type that contains time and position in a source
// file when an Err occured.
type ErrOcc struct {
	time time.Time
	file string
	line int
}

// Time is the method to get time when this Err occured.
func (e ErrOcc) Time() time.Time {
	return e.time
}

// Line is the method to get the line number where this Err occured.
func (e ErrOcc) Line() int {
	return e.line
}

// File is the method to get the file name where this Err occured.
func (e ErrOcc) File() string {
	return e.file
}

type handlerListEntry struct {
	handler func(Err, ErrOcc)
	next    *handlerListEntry
}

type handlerList struct {
	head *handlerListEntry
	last *handlerListEntry
}

var (
	syncErrHandlers  = handlerList{nil, nil}
	asyncErrHandlers = handlerList{nil, nil}
	isErrCfgFixed    = false
)

// AddSyncHandler is the function that adds an Err creation event handler.
// Handlers added with this method are executed synchronously in the order of
// addition.
func AddSyncHandler(handler func(Err, ErrOcc)) {
	if isErrCfgFixed {
		return
	}

	last := syncErrHandlers.last
	syncErrHandlers.last = &handlerListEntry{handler, nil}

	if last != nil {
		last.next = syncErrHandlers.last
	}

	if syncErrHandlers.head == nil {
		syncErrHandlers.head = syncErrHandlers.last
	}
}

// AddAsyncHandler is the function that adds an Err creation event handler.
// Handlers added with this method are executed asynchronously.
func AddAsyncHandler(handler func(Err, ErrOcc)) {
	if isErrCfgFixed {
		return
	}

	last := asyncErrHandlers.last
	asyncErrHandlers.last = &handlerListEntry{handler, nil}

	if last != nil {
		last.next = asyncErrHandlers.last
	}

	if asyncErrHandlers.head == nil {
		asyncErrHandlers.head = asyncErrHandlers.last
	}
}

// FixCfg is the function to fix the configuration of error processing.
// After calling this function, handlers cannot be added any more and the
// notification becomes effective.
func FixCfg() {
	isErrCfgFixed = true
}

func notifyErr(err Err) {
	if !isErrCfgFixed {
		return
	}

	if syncErrHandlers.head == nil && asyncErrHandlers.head == nil {
		return
	}

	var occ ErrOcc
	occ.time = time.Now()

	_, file, line, ok := runtime.Caller(2)
	if ok {
		occ.file = filepath.Base(file)
		occ.line = line
	}

	for el := syncErrHandlers.head; el != nil; el = el.next {
		el.handler(err, occ)
	}

	for el := asyncErrHandlers.head; el != nil; el = el.next {
		go el.handler(err, occ)
	}
}
