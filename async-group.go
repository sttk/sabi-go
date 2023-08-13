// Copyright (C) 2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"

	"github.com/sttk/sabi/errs"
)

// AsyncGroup is the interface to execute added functions asynchronously.
// The method: Add is to add target functions.
// This interface is used as an argument of DaxSrc#Setup, DaxConn#Commit, and DaxConn#Rollback.
type AsyncGroup interface {
	Add(fn func() errs.Err)
}

type errEntry[N comparable] struct {
	name N
	err  errs.Err
	next *errEntry[N]
}

type asyncGroupAsync[N comparable] struct {
	wg      sync.WaitGroup
	errHead *errEntry[N]
	errLast *errEntry[N]
	mutex   sync.Mutex
	name    N
}

func (ag *asyncGroupAsync[N]) Add(fn func() errs.Err) {
	ag.wg.Add(1)
	go func(name N) {
		defer ag.wg.Done()
		err := fn()
		if err.IsNotOk() {
			ag.mutex.Lock()
			defer ag.mutex.Unlock()
			ag.addErr(name, err)
		}
	}(ag.name)
}

func (ag *asyncGroupAsync[N]) wait() {
	ag.wg.Wait()
}

func (ag *asyncGroupAsync[N]) addErr(name N, err errs.Err) {
	ent := &errEntry[N]{name: name, err: err}
	if ag.errLast == nil {
		ag.errHead = ent
		ag.errLast = ent
	} else {
		ag.errLast.next = ent
		ag.errLast = ent
	}
}

func (ag *asyncGroupAsync[N]) hasErr() bool {
	return (ag.errHead != nil)
}

func (ag *asyncGroupAsync[N]) makeErrs() map[N]errs.Err {
	m := make(map[N]errs.Err)
	for ent := ag.errHead; ent != nil; ent = ent.next {
		m[ent.name] = ent.err
	}
	return m
}

type asyncGroupSync struct {
	err errs.Err
}

func (ag *asyncGroupSync) Add(fn func() errs.Err) {
	ag.err = fn()
}
