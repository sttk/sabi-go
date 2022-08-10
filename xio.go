// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"
)

type /* error reasons */ (
	// XioConnCfgIsNotFound is an error reason which indicates that a connection
	// configuration to an external data source is not found
	XioConnCfgIsNotFound struct {
		Name string
	}

	// FailToCreateXioConn is an error reason which indicates that it is failed
	// to create a new connection to an external data source.
	FailToCreateXioConn struct {
		Name string
	}

	// FailToCommitXioConn is an error reason which indicates that some
	// connections to external data sources failed to commit.
	FailToCommitXioConn struct {
		Errors map[string]Err
	}

	// FailToRollbackXioConn is an error reason which indicates that some
	// connections to external data sources failed to rollback.
	FailToRollbackXioConn struct {
		Errors map[string]Err
	}

	// FailToCloseXioConn is an error reason which indicates that some
	// connections to external data sources failed to close.
	FailToCloseXioConn struct {
		Errors map[string]Err
	}
)

// Xio is an interface for a set of inputs/outputs, and requires 2 methods:
// #GetConn which gets a connection to an external data sourc, and #innerMap
// which gets a map to communicate data among multiple inputs/outputs.
type Xio interface {
	GetConn(name string) (XioConn, Err)
	InnerMap() map[string]any
}

// XioConn is an interface which represents a connection to an external data
// source and requires a methods: #Commit, #Rollback, and #Close to work
// in a transaction process.
type XioConn interface {
	Commit() Err
	Rollback() Err
	Close() Err
}

// XioConnCfg is an interface which requires a method: #NewConn which creates
// a connection to a data source with configuration parameters.
type XioConnCfg interface {
	NewConn() (XioConn, Err)
}

var (
	isGlobalXioConnCfgSealed bool                  = false
	globalXioConnCfgMap      map[string]XioConnCfg = make(map[string]XioConnCfg)
	globalXioConnCfgMutex    sync.Mutex
)

// AddGlobalXioConnCfg registers a global XioConnCfg with its name to make
// enable to use XioConn in all transactions.
func AddGlobalXioConnCfg(name string, cfg XioConnCfg) {
	globalXioConnCfgMutex.Lock()
	defer globalXioConnCfgMutex.Unlock()

	if !isGlobalXioConnCfgSealed {
		globalXioConnCfgMap[name] = cfg
	}
}

// SealGlobalXioConnCfgs makes unable to register any further global
// XioConnCfg.
func SealGlobalXioConnCfgs() {
	isGlobalXioConnCfgSealed = true
}

// XioBase is a structure type which is used as an implementation of Xio
// interface, and manages one or more XioConnCfg and XioConn used in
// a transaction.
type XioBase struct {
	isLocalXioConnCfgSealed bool
	localXioConnCfgMap      map[string]XioConnCfg
	xioConnMap              map[string]XioConn
	xioConnMutex            sync.Mutex
	innerMap                map[string]any
	error                   Err
}

// NewXioBase is creates a new XioBase.
func NewXioBase() *XioBase {
	return &XioBase{
		isLocalXioConnCfgSealed: false,
		localXioConnCfgMap:      make(map[string]XioConnCfg),
		xioConnMap:              make(map[string]XioConn),
		innerMap:                make(map[string]any),
		error:                   Ok(),
	}
}

// AddLocalXioConnCfg registers a transaction-local XioConnCfg with its name.
func (base *XioBase) AddLocalXioConnCfg(name string, cfg XioConnCfg) {
	base.xioConnMutex.Lock()
	defer base.xioConnMutex.Unlock()

	if !base.isLocalXioConnCfgSealed {
		base.localXioConnCfgMap[name] = cfg
	}
}

// SealLocalXioConnCfgs makes unable to register any further transaction-local
// XioConnCfg.
func (base *XioBase) SealLocalXioConnCfgs() {
	base.isLocalXioConnCfgSealed = true
	isGlobalXioConnCfgSealed = true
}

// GetConn gets a XioConn which is a connection to an external data source by
// specified name. If a XioConn is not found, this method creates new one with
// a local or global XioConnCfg associated with same name.
func (base *XioBase) GetConn(name string) (XioConn, Err) {
	conn := base.xioConnMap[name]
	if conn != nil {
		return conn, Ok()
	}

	cfg := base.localXioConnCfgMap[name]
	if cfg == nil {
		cfg = globalXioConnCfgMap[name]
	}
	if cfg == nil {
		return nil, ErrBy(XioConnCfgIsNotFound{Name: name})
	}

	base.xioConnMutex.Lock()
	defer base.xioConnMutex.Unlock()

	var err Err
	conn, err = cfg.NewConn()
	if !err.IsOk() {
		return nil, ErrBy(FailToCreateXioConn{Name: name}, err)
	}

	base.xioConnMap[name] = conn

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

	for name, conn := range base.xioConnMap {
		go func(name string, conn XioConn, ch chan namedErr) {
			err := conn.Commit()
			ne := namedErr{name: name, err: err}
			ch <- ne
		}(name, conn, ch)
	}

	errs := make(map[string]Err)
	n := len(base.xioConnMap)
	for i := 0; i < n; i++ {
		select {
		case ne := <-ch:
			if !ne.err.IsOk() {
				errs[ne.name] = ne.err
			}
		}
	}

	if len(errs) > 0 {
		return ErrBy(FailToCommitXioConn{Errors: errs})
	}

	return Ok()
}

func (base *XioBase) rollback() Err {
	ch := make(chan namedErr)

	for name, conn := range base.xioConnMap {
		go func(name string, conn XioConn, ch chan namedErr) {
			err := conn.Rollback()
			ne := namedErr{name: name, err: err}
			ch <- ne
		}(name, conn, ch)
	}

	errs := make(map[string]Err)
	n := len(base.xioConnMap)
	for i := 0; i < n; i++ {
		select {
		case ne := <-ch:
			if !ne.err.IsOk() {
				errs[ne.name] = ne.err
			}
		}
	}

	if len(errs) > 0 {
		return ErrBy(FailToRollbackXioConn{Errors: errs})
	}

	return Ok()
}

func (base *XioBase) close() Err {
	ch := make(chan namedErr)

	for name, conn := range base.xioConnMap {
		go func(name string, conn XioConn, ch chan namedErr) {
			err := conn.Close()
			ne := namedErr{name: name, err: err}
			ch <- ne
		}(name, conn, ch)
	}

	errs := make(map[string]Err)
	n := len(base.xioConnMap)
	for i := 0; i < n; i++ {
		select {
		case ne := <-ch:
			if !ne.err.IsOk() {
				errs[ne.name] = ne.err
			}
		}
	}

	if len(errs) > 0 {
		return ErrBy(FailToCloseXioConn{Errors: errs})
	}

	return Ok()
}
