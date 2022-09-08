// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"
)

type /* error reasons */ (
	// ConnCfgIsNotFound is an error reason which indicates that a connection
	// configuration to an external data source is not found
	ConnCfgIsNotFound struct {
		Name string
	}

	// FailToCreateConn is an error reason which indicates that it is failed
	// to create a new connection to an external data source.
	FailToCreateConn struct {
		Name string
	}

	// FailToCommitConn is an error reason which indicates that some
	// connections to external data sources failed to commit.
	FailToCommitConn struct {
		Errors map[string]Err
	}

	// FailToRollbackConn is an error reason which indicates that some
	// connections to external data sources failed to rollback.
	FailToRollbackConn struct {
		Errors map[string]Err
	}

	// FailToCloseConn is an error reason which indicates that some
	// connections to external data sources failed to close.
	FailToCloseConn struct {
		Errors map[string]Err
	}
)

// Xio is an interface for a set of inputs/outputs, and requires 2 methods:
// #GetConn which gets a connection to an external data sourc, and #innerMap
// which gets a map to communicate data among multiple inputs/outputs.
type Xio interface {
	GetConn(name string) (Conn, Err)
	InnerMap() map[string]any
}

// XioBase is a structure type which is used as an implementation of Xio
// interface, and manages one or more ConnCfg and Conn used in
// a transaction.
type XioBase struct {
	isLocalConnCfgSealed bool
	localConnCfgMap      map[string]ConnCfg
	connMap              map[string]Conn
	connMutex            sync.Mutex
	innerMap             map[string]any
	error                Err
}

// NewXioBase is creates a new XioBase.
func NewXioBase() *XioBase {
	return &XioBase{
		isLocalConnCfgSealed: false,
		localConnCfgMap:      make(map[string]ConnCfg),
		connMap:              make(map[string]Conn),
		innerMap:             make(map[string]any),
		error:                Ok(),
	}
}

// AddLocalConnCfg registers a transaction-local ConnCfg with its name.
func (base *XioBase) AddLocalConnCfg(name string, cfg ConnCfg) {
	base.connMutex.Lock()
	defer base.connMutex.Unlock()

	if !base.isLocalConnCfgSealed {
		base.localConnCfgMap[name] = cfg
	}
}

// SealLocalConnCfgs makes unable to register any further transaction-local
// ConnCfg.
func (base *XioBase) SealLocalConnCfgs() {
	base.isLocalConnCfgSealed = true
	isGlobalConnCfgSealed = true
}

// GetConn gets a Conn which is a connection to an external data source by
// specified name. If a Conn is not found, this method creates new one with
// a local or global ConnCfg associated with same name.
func (base *XioBase) GetConn(name string) (Conn, Err) {
	conn := base.connMap[name]
	if conn != nil {
		return conn, Ok()
	}

	cfg := base.localConnCfgMap[name]
	if cfg == nil {
		cfg = globalConnCfgMap[name]
	}
	if cfg == nil {
		return nil, ErrBy(ConnCfgIsNotFound{Name: name})
	}

	base.connMutex.Lock()
	defer base.connMutex.Unlock()

	var err Err
	conn, err = cfg.CreateConn()
	if !err.IsOk() {
		return nil, ErrBy(FailToCreateConn{Name: name}, err)
	}

	base.connMap[name] = conn

	return conn, Ok()
}

// InnerMap gets a singular map in a transaction for communicating among
// multiple Xio data operations.
func (base *XioBase) InnerMap() map[string]any {
	return base.innerMap
}

type namedErr struct {
	name string
	err  Err
}

func (base *XioBase) commit() Err {
	ch := make(chan namedErr)

	for name, conn := range base.connMap {
		go func(name string, conn Conn, ch chan namedErr) {
			err := conn.Commit()
			ne := namedErr{name: name, err: err}
			ch <- ne
		}(name, conn, ch)
	}

	errs := make(map[string]Err)
	n := len(base.connMap)
	for i := 0; i < n; i++ {
		select {
		case ne := <-ch:
			if !ne.err.IsOk() {
				errs[ne.name] = ne.err
			}
		}
	}

	if len(errs) > 0 {
		return ErrBy(FailToCommitConn{Errors: errs})
	}

	return Ok()
}

func (base *XioBase) rollback() {
	var wg sync.WaitGroup
	wg.Add(len(base.connMap))

	for _, conn := range base.connMap {
		go func(conn Conn) {
			defer wg.Done()
			conn.Rollback()
		}(conn)
	}

	wg.Wait()
}

func (base *XioBase) close() {
	var wg sync.WaitGroup
	wg.Add(len(base.connMap))

	for _, conn := range base.connMap {
		go func(conn Conn) {
			defer wg.Done()
			conn.Close()
		}(conn)
	}

	wg.Wait()
}
