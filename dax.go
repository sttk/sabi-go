// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"
)

type /* error reasons */ (
	// DaxSrcIsNotFound is an error reason which indicates that a specified data
	// source instance is not found.
	// The field Name is a registered name of a DaxSrc is not found.
	DaxSrcIsNotFound struct {
		Name string
	}

	// FailToCreateDaxConn is an error reason which indicates that it failed to
	// create a new connection to a data source.
	// The field Name is a registered name of a DaxSrc which failed to create a
	// DaxConn.
	FailToCreateDaxConn struct {
		Name string
	}

	// FailToCommitDaxConn is an error reason interface which indicates that some
	// connections failed to commit.
	// The field Errors is a map of which keys are registered names of DaxConn
	// which failed to commit, and of which values are Err instances holding
	// their error reasons.
	FailToCommitDaxConn struct {
		Errors map[string]Err
	}
)

// DaxConn is an interface which represents a connection to a data source, and
// requires methods: #Commit, #Rollback and #Close to work in a tranaction
// process.
type DaxConn interface {
	Commit() Err
	Rollback()
	Close()
}

// DaxSrc is an interface which represents a data source like database, etc.,
// and creates a DaxConn to the data source.
// This requires a method: #CreateDaxConn to do so.
type DaxSrc interface {
	CreateDaxConn() (DaxConn, Err)
}

// Dax is an interface for a set of data accesses, and requires a method:
// #GetConn which gets a connection to an external data access.
type Dax interface {
	GetDaxConn(name string) (DaxConn, Err)
}

var (
	isGlobalDaxSrcsFixed bool              = false
	globalDaxSrcMap      map[string]DaxSrc = make(map[string]DaxSrc)
	globalDaxSrcMutex    sync.Mutex
)

// AddGlobalDaxSrc registers a global DaxSrc with its name to make enable to
// use DaxSrc in all transactions.
func AddGlobalDaxSrc(name string, ds DaxSrc) {
	globalDaxSrcMutex.Lock()
	defer globalDaxSrcMutex.Unlock()

	if !isGlobalDaxSrcsFixed {
		globalDaxSrcMap[name] = ds
	}
}

// FixGlobalDaxSrcs makes unable to register any further global DaxSrc.
func FixGlobalDaxSrcs() {
	isGlobalDaxSrcsFixed = true
}

// DaxBase is a structure type which manages multiple DaxSrc and those DaxConn,
// and also work as an implementation of Dax interface.
type DaxBase struct {
	isLocalDaxSrcsFixed bool
	localDaxSrcMap      map[string]DaxSrc
	daxConnMap          map[string]DaxConn
	daxConnMutex        sync.Mutex
}

// NewDaxBase is a function which creates a new DaxBase.
func NewDaxBase() *DaxBase {
	return &DaxBase{
		isLocalDaxSrcsFixed: false,
		localDaxSrcMap:      make(map[string]DaxSrc),
		daxConnMap:          make(map[string]DaxConn),
	}
}

// AddLocalDaxSrc is a method which registers a local DaxSrc with a specified
// name.
func (base *DaxBase) AddLocalDaxSrc(name string, ds DaxSrc) {
	base.daxConnMutex.Lock()
	defer base.daxConnMutex.Unlock()

	if !base.isLocalDaxSrcsFixed {
		base.localDaxSrcMap[name] = ds
	}
}

// GetDaxConn gets a DaxConn which is a connection to a data source by
// specified name.
// If a DaxConn is found, this method returns it, but not found, creates a new
// one with a local or global DaxSrc associated with same name.
// If there are both local and global DaxSrc with same name, the local DaxSrc
// is used.
func (base *DaxBase) GetDaxConn(name string) (DaxConn, Err) {
	conn := base.daxConnMap[name]
	if conn != nil {
		return conn, Ok()
	}

	ds := base.localDaxSrcMap[name]
	if ds == nil {
		ds = globalDaxSrcMap[name]
	}
	if ds == nil {
		return nil, ErrBy(DaxSrcIsNotFound{Name: name})
	}

	base.daxConnMutex.Lock()
	defer base.daxConnMutex.Unlock()

	conn = base.daxConnMap[name]
	if conn != nil {
		return conn, Ok()
	}

	var err Err
	conn, err = ds.CreateDaxConn()
	if !err.IsOk() {
		return nil, ErrBy(FailToCreateDaxConn{Name: name}, err)
	}

	base.daxConnMap[name] = conn

	return conn, Ok()
}

func (base *DaxBase) begin() {
	base.isLocalDaxSrcsFixed = true
	isGlobalDaxSrcsFixed = true
}

type namedErr struct {
	name string
	err  Err
}

func (base *DaxBase) commit() Err {
	ch := make(chan namedErr)

	for name, conn := range base.daxConnMap {
		go func(name string, conn DaxConn, ch chan namedErr) {
			err := conn.Commit()
			ne := namedErr{name: name, err: err}
			ch <- ne
		}(name, conn, ch)
	}

	errs := make(map[string]Err)
	n := len(base.daxConnMap)
	for i := 0; i < n; i++ {
		select {
		case ne := <-ch:
			if !ne.err.IsOk() {
				errs[ne.name] = ne.err
			}
		}
	}

	if len(errs) > 0 {
		return ErrBy(FailToCommitDaxConn{Errors: errs})
	}

	return Ok()
}

func (base *DaxBase) rollback() {
	var wg sync.WaitGroup
	wg.Add(len(base.daxConnMap))

	for _, conn := range base.daxConnMap {
		go func(conn DaxConn) {
			defer wg.Done()
			conn.Rollback()
		}(conn)
	}

	wg.Wait()
}

func (base *DaxBase) close() {
	var wg sync.WaitGroup
	wg.Add(len(base.daxConnMap))

	for _, conn := range base.daxConnMap {
		go func(conn DaxConn) {
			defer wg.Done()
			conn.Close()
		}(conn)
	}

	wg.Wait()

	base.isLocalDaxSrcsFixed = false
}
