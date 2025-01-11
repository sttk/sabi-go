// Copyright (C) 2023-2025 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"
)

type errEntry struct {
	name string
	err  Err
	next *errEntry
}

// AsyncGroup manages multiple asynchronous operations and allows for waiting until all operations
// have completed.
// It also handles error collection to store errors that encounter during executions.
type AsyncGroup struct {
	wg      sync.WaitGroup
	errHead *errEntry
	errLast *errEntry
	mutex   sync.Mutex
	name    string
}

// Adds a new asynchronous operation to this instance.
// The operation is provided as a function that returns an Err.
// If the function encounters an error, it is recorded internally.
func (ag *AsyncGroup) Add(fn func() Err) {
	ag.wg.Add(1)
	go func(name string) {
		defer ag.wg.Done()
		err := fn()
		if err.IsNotOk() {
			ag.mutex.Lock()
			defer ag.mutex.Unlock()
			ag.addErr(name, err)
		}
	}(ag.name)
}

func (ag *AsyncGroup) join() map[string]Err {
	ag.wg.Wait()

	m := make(map[string]Err)
	for ent := ag.errHead; ent != nil; ent = ent.next {
		m[ent.name] = ent.err
	}
	return m
}

func (ag *AsyncGroup) addErr(name string, err Err) {
	ent := &errEntry{name: name, err: err}
	if ag.errLast == nil {
		ag.errHead = ent
		ag.errLast = ent
	} else {
		ag.errLast.next = ent
		ag.errLast = ent
	}
}
