// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"
	"time"
)

type handlerListElem struct {
	handler func(Err, time.Time)
	next    *handlerListElem
}

type handlerList struct {
	head *handlerListElem
	last *handlerListElem
}

var (
	syncErrHandlers  = handlerList{nil, nil}
	asyncErrHandlers = handlerList{nil, nil}
	isErrCfgSealed   = false
	errCfgMutex      = sync.Mutex{}
)

// Adds an Err creation event handler which is executed synchronously.
// Handlers added with this method are executed in the order of addition.
func AddSyncErrHandler(handler func(Err, time.Time)) {
	errCfgMutex.Lock()
	defer errCfgMutex.Unlock()

	if isErrCfgSealed {
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
func AddAsyncErrHandler(handler func(Err, time.Time)) {
	errCfgMutex.Lock()
	defer errCfgMutex.Unlock()

	if isErrCfgSealed {
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

// Seals configuration for Err creation event handlers.
// After calling this function, handlers cannot be registered any more and the
// notification becomes effective.
func SealErrCfgs() {
	isErrCfgSealed = true
}

func notifyErr(err Err) {
	if !isErrCfgSealed {
		return
	}

	if syncErrHandlers.head == nil && asyncErrHandlers.head == nil {
		return
	}

	now := time.Now()

	for el := syncErrHandlers.head; el != nil; el = el.next {
		el.handler(err, now)
	}

	if asyncErrHandlers.head != nil {
		go func() {
			for el := asyncErrHandlers.head; el != nil; el = el.next {
				go el.handler(err, now)
			}
		}()
	}
}
