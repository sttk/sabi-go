// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// ErrOccasion is a struct which contains time and posision in a source file
// when an Err occured.
type ErrOccasion struct {
	time time.Time
	file string
	line int
}

// Time is a method which returns time when this Err occured.
func (e ErrOccasion) Time() time.Time {
	return e.time
}

// File is a method which returns the file name where this Err occured.
func (e ErrOccasion) File() string {
	return e.file
}

// Line is a method which returns the line number where this Err occured.
func (e ErrOccasion) Line() int {
	return e.line
}

type handlerListElem struct {
	handler func(Err, ErrOccasion)
	next    *handlerListElem
}

type handlerList struct {
	head *handlerListElem
	last *handlerListElem
}

var (
	syncErrHandlers  = handlerList{nil, nil}
	asyncErrHandlers = handlerList{nil, nil}
	isErrCfgsFixed   = false
	errCfgMutex      = sync.Mutex{}
)

// Adds an Err creation event handler which is executed synchronously.
// Handlers added with this method are executed in the order of addition.
func AddSyncErrHandler(handler func(Err, ErrOccasion)) {
	errCfgMutex.Lock()
	defer errCfgMutex.Unlock()

	if isErrCfgsFixed {
		return
	}

	last := syncErrHandlers.last
	syncErrHandlers.last = &handlerListElem{handler, nil}

	if last != nil {
		last.next = syncErrHandlers.last
	}

	if syncErrHandlers.head == nil {
		syncErrHandlers.head = syncErrHandlers.last
	}
}

// Adds a Err creation event handlers which is executed asynchronously.
func AddAsyncErrHandler(handler func(Err, ErrOccasion)) {
	errCfgMutex.Lock()
	defer errCfgMutex.Unlock()

	if isErrCfgsFixed {
		return
	}

	last := asyncErrHandlers.last
	asyncErrHandlers.last = &handlerListElem{handler, nil}

	if last != nil {
		last.next = asyncErrHandlers.last
	}

	if asyncErrHandlers.head == nil {
		asyncErrHandlers.head = asyncErrHandlers.last
	}
}

// Fixes configuration for Err creation event handlers.
// After calling this function, handlers cannot be registered any more and the
// notification becomes effective.
func FixErrCfgs() {
	isErrCfgsFixed = true
}

func notifyErr(err Err) {
	if !isErrCfgsFixed {
		return
	}

	if syncErrHandlers.head == nil && asyncErrHandlers.head == nil {
		return
	}

	var occ ErrOccasion
	occ.time = time.Now()

	_, file, line, ok := runtime.Caller(2)
	if ok {
		occ.file = filepath.Base(file)
		occ.line = line
	}

	for el := syncErrHandlers.head; el != nil; el = el.next {
		el.handler(err, occ)
	}

	if asyncErrHandlers.head != nil {
		for el := asyncErrHandlers.head; el != nil; el = el.next {
			go el.handler(err, occ)
		}
	}
}
