// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"
)

type /* error reasons */ (
	// ConnCfgIsNotFound is an error reason which indicates that a conection
	// configuration to an external data source is not found.
	ConnCfgIsNotFound struct {
		Name string
	}

	// FailToCreateConn is an error reason which indicates that it is failed
	// to create a new connection to an external data source.
	FailToCreateConn struct {
		Name string
	}

	// FailToCommitConn is an error reason which indicates that some connections
	// to external data sources failed to commit.
	FailToCommitConn struct {
		Errors map[string]Err
	}
)

// Conn is an interface which represents a connection to an external data
// source and requires a methods: #Commit, #Roolback and #Close to work in a
// transaction process.
type Conn interface {
	Commit() Err
	Rollback()
	Close()
}

// ConnCfg is an interface which requires a method: #CreateConn which creates
// a connection to a data source with configuration parameters.
type ConnCfg interface {
	CreateConn() (Conn, Err)
}

var (
	isGlobalConnCfgSealed bool               = false
	globalConnCfgMap      map[string]ConnCfg = make(map[string]ConnCfg)
	globalConnCfgMutex    sync.Mutex
)

// AddGlobalConnCfg registers a global ConnCfg with its name to make enable
// to use ConnCfg in all processes..
func AddGlobalConnCfg(name string, cfg ConnCfg) {
	globalConnCfgMutex.Lock()
	defer globalConnCfgMutex.Unlock()

	if !isGlobalConnCfgSealed {
		globalConnCfgMap[name] = cfg
	}
}

// SealGlobalConnCfgs makes unable to register any further global ConnCfg.
func SealGlobalConnCfgs() {
	isGlobalConnCfgSealed = true
}

// ConnBase is a structure type which manages multiple Conn and ConnCfg, and
// also work as an implementation of Dax interface.
type ConnBase struct {
	isLocalConnCfgSealed bool
	localConnCfgMap      map[string]ConnCfg
	connMap              map[string]Conn
	innerMap             map[string]any
	connMutex            sync.Mutex
}

// NewConnBase is a function which creates a new ConnBase.
func NewConnBase() *ConnBase {
	return &ConnBase{
		isLocalConnCfgSealed: false,
		localConnCfgMap:      make(map[string]ConnCfg),
		connMap:              make(map[string]Conn),
		innerMap:             make(map[string]any),
	}
}

func (base *ConnBase) addLocalConnCfg(name string, cfg ConnCfg) {
	base.connMutex.Lock()
	defer base.connMutex.Unlock()

	if !base.isLocalConnCfgSealed {
		base.localConnCfgMap[name] = cfg
	}
}

// GetConn gets a Conn which is a connection to an external data source by
// specified name. If a Conn is not found, this method creates new one with
// a local or global ConnCfg associated with same name.
func (base *ConnBase) GetConn(name string) (Conn, Err) {
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

// InnerMap gets a singular map in a ConnBase/Proc for communicating among
// multiple data accesses.
func (base *ConnBase) InnerMap() map[string]any {
	return base.innerMap
}

func (base *ConnBase) begin() {
	base.isLocalConnCfgSealed = true
	isGlobalConnCfgSealed = true
}

type namedErr struct {
	name string
	err  Err
}

func (base *ConnBase) commit() Err {
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

func (base *ConnBase) rollback() {
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

func (base *ConnBase) close() {
	var wg sync.WaitGroup
	wg.Add(len(base.connMap))

	for _, conn := range base.connMap {
		go func(conn Conn) {
			defer wg.Done()
			conn.Close()
		}(conn)
	}

	wg.Wait()

	base.isLocalConnCfgSealed = false
}
