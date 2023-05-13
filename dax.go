// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"reflect"
	"sync"
)

type /* error reasons */ (
	// FailToStartUpGlobalDaxSrcs is an error reason which indicates that some
	// dax sources failed to start up.
	// The field: Errors is a map of which keys are the registered names of
	// DaxSrc(s) which failed to start up, and of which values are Err having
	// their error reasons.
	FailToStartUpGlobalDaxSrcs struct {
		Errors map[string]Err
	}

	// DaxSrcIsNotFound is an error reason which indicates that a specified data
	// source is not found.
	// The field: Name is a registered name of a DaxSrc not found.
	DaxSrcIsNotFound struct {
		Name string
	}

	// FailToCreateDaxConn is an error reason which indicates that it's failed to
	// create a new connection to a data store.
	// The field: Name is a registered name of a DaxSrc which failed to create a
	// DaxConn.
	FailToCreateDaxConn struct {
		Name string
	}

	// FailToCastConn is an error reason which indicates that it's failed to
	// cast type of a DaxConn.
	// The field: Name is a registered name of a DaxSrc which created to a target
	// DaxConn.
	// And the field: Type is a destination type name.
	FailToCastDaxConn struct {
		Name, FromType, ToType string
	}

	// FailToCommitDaxConn is an error reason which indicates that some
	// connections failed to commit.
	// The field: Errors is a map of which keys are the registered names of
	// DaxConn which failed to commit, and of which values are Err having their
	// error reasons.
	FailToCommitDaxConn struct {
		Errors map[string]Err
	}
)

// DaxConn is an interface which represents a connection to a data store, and
// defines methods: Commit, Rollback and Close to work in a tranaction process.
type DaxConn interface {
	Commit() Err
	Rollback()
	Close()
}

// DaxSrc is an interface which represents a data connection source for a data
// store like database, etc., and creates a DaxConn which is a connection to a
// data store.
// This interface defines a method: CreateDaxConn to creates a DaxConn instance
// and returns its pointer.
// This interface also defines methods: SetUp and End, which makes available
// and frees this dax source.
type DaxSrc interface {
	CreateDaxConn() (DaxConn, Err)
	SetUp() Err
	End()
}

// Dax is an interface for a set of data access methods.
// This method gets a DaxConn which is a connection to a data store by
// specified name.
// If a DaxConn is found, this method returns it, but not found, creates a new
// one with a local or global DaxSrc associated with same name.
// If there are both local and global DaxSrc with same name, the local DaxSrc
// is used.
type Dax interface {
	getDaxConn(name string) (DaxConn, Err)
}

var (
	isGlobalDaxSrcsFixed bool              = false
	globalDaxSrcMap      map[string]DaxSrc = make(map[string]DaxSrc)
	globalDaxSrcMutex    sync.Mutex
)

// AddGlobalDaxSrc registers a global DaxSrc with its name to make enable to
// use DaxSrc in all transactions.
// This method ignores to add a global DaxSrc when its name is already
// registered.
// In addition, this method ignores to add any more global DaxSrc(s) after
// calling FixGlobalDaxSrcs function.
func AddGlobalDaxSrc(name string, ds DaxSrc) {
	globalDaxSrcMutex.Lock()
	defer globalDaxSrcMutex.Unlock()

	if !isGlobalDaxSrcsFixed {
		_, exists := globalDaxSrcMap[name]
		if !exists {
			globalDaxSrcMap[name] = ds
		}
	}
}

// StartUpGlobalDaxSrcs is a function to forbid adding more global dax sources
// and to make available the registered global dax sources by calling Setup
// method of each DaxSrc.
// If even one DaxSrc fail to execute its SstUp method, this function
// executes Free methods of all global DaxSrc(s) and returns sabi.Err.
func StartUpGlobalDaxSrcs() Err {
	isGlobalDaxSrcsFixed = true

	ch := make(chan namedErr)

	for name, ds := range globalDaxSrcMap {
		go func(name string, ds DaxSrc, ch chan namedErr) {
			err := ds.SetUp()
			ne := namedErr{name: name, err: err}
			ch <- ne
		}(name, ds, ch)
	}

	errs := make(map[string]Err)
	n := len(globalDaxSrcMap)
	for i := 0; i < n; i++ {
		select {
		case ne := <-ch:
			if !ne.err.IsOk() {
				errs[ne.name] = ne.err
			}
		}
	}

	if len(errs) > 0 {
		ShutdownGlobalDaxSrcs()
		return NewErr(FailToStartUpGlobalDaxSrcs{Errors: errs})
	}

	return Ok()
}

// ShutdownGlobalDaxSrcs is a function to terminate all global dax sources.
func ShutdownGlobalDaxSrcs() {
	var wg sync.WaitGroup
	wg.Add(len(globalDaxSrcMap))

	for _, ds := range globalDaxSrcMap {
		go func(ds DaxSrc) {
			defer wg.Done()
			ds.End()
		}(ds)
	}

	wg.Wait()
}

// DaxBase is an interface which works as a front of an implementation as a
// base of data connection sources, and defines methods: SetUpLocalDaxSrc and
// FreeLocalDaxSrc.
//
// SetUpLocalDaxSrc method registered a DaxSrc with a name in this
// implementation, but  ignores to add a local DaxSrc when its name is already
// registered.
// In addition, this method ignores to add local DaxSrc(s) while the
// transaction is processing.
//
// This interface inherits Dax interface to get a DaxConn by a name.
// Also, this has unexported methods for a transaction process.
type DaxBase interface {
	Dax
	SetUpLocalDaxSrc(name string, ds DaxSrc) Err
	FreeLocalDaxSrc(name string)
	FreeAllLocalDaxSrcs()
	begin()
	commit() Err
	rollback()
	end()
}

type daxBaseImpl struct {
	isLocalDaxSrcsFixed bool
	localDaxSrcMap      map[string]DaxSrc
	daxConnMap          map[string]DaxConn
	daxConnMutex        sync.Mutex
}

// NewDaxBase is a function which creates a new DaxBase.
func NewDaxBase() DaxBase {
	return &daxBaseImpl{
		isLocalDaxSrcsFixed: false,
		localDaxSrcMap:      make(map[string]DaxSrc),
		daxConnMap:          make(map[string]DaxConn),
	}
}

func (base *daxBaseImpl) SetUpLocalDaxSrc(name string, ds DaxSrc) Err {
	base.daxConnMutex.Lock()
	defer base.daxConnMutex.Unlock()

	if !base.isLocalDaxSrcsFixed {
		_, exists := base.localDaxSrcMap[name]
		if !exists {
			err := ds.SetUp()
			if !err.IsOk() {
				return err
			}
			base.localDaxSrcMap[name] = ds
		}
	}

	return Ok()
}

func (base *daxBaseImpl) FreeLocalDaxSrc(name string) {
	base.daxConnMutex.Lock()
	defer base.daxConnMutex.Unlock()

	if !base.isLocalDaxSrcsFixed {
		ds, exists := base.localDaxSrcMap[name]
		if exists {
			delete(base.localDaxSrcMap, name)
			ds.End()
		}
	}
}

func (base *daxBaseImpl) FreeAllLocalDaxSrcs() {
	base.daxConnMutex.Lock()
	defer base.daxConnMutex.Unlock()

	if !base.isLocalDaxSrcsFixed {
		for _, ds := range base.localDaxSrcMap {
			ds.End()
		}

		base.localDaxSrcMap = make(map[string]DaxSrc)
	}
}

func (base *daxBaseImpl) getDaxConn(name string) (DaxConn, Err) {
	conn := base.daxConnMap[name]
	if conn != nil {
		return conn, Ok()
	}

	ds := base.localDaxSrcMap[name]
	if ds == nil {
		ds = globalDaxSrcMap[name]
	}
	if ds == nil {
		return nil, NewErr(DaxSrcIsNotFound{Name: name})
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
		return nil, NewErr(FailToCreateDaxConn{Name: name}, err)
	}

	base.daxConnMap[name] = conn

	return conn, Ok()
}

func (base *daxBaseImpl) begin() {
	base.isLocalDaxSrcsFixed = true
	isGlobalDaxSrcsFixed = true
}

type namedErr struct {
	name string
	err  Err
}

func (base *daxBaseImpl) commit() Err {
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
		return NewErr(FailToCommitDaxConn{Errors: errs})
	}

	return Ok()
}

func (base *daxBaseImpl) rollback() {
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

func (base *daxBaseImpl) end() {
	var wg sync.WaitGroup
	wg.Add(len(base.daxConnMap))

	for _, conn := range base.daxConnMap {
		go func(conn DaxConn) {
			defer wg.Done()
			conn.Close()
		}(conn)
	}

	base.daxConnMap = make(map[string]DaxConn)

	wg.Wait()

	base.isLocalDaxSrcsFixed = false
}

// GetDaxConn is a function to cast type of DaxConn instance.
// If it's failed to cast to a destination type, this function returns an Err
// of a reason: FailToGetDaxConn.
func GetDaxConn[C any](dax Dax, name string) (C, Err) {
	conn, err := dax.getDaxConn(name)
	if err.IsOk() {
		casted, ok := conn.(C)
		if ok {
			return casted, err
		}

		var from string
		t := reflect.TypeOf(conn)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			from = "*" + t.Name() + " (" + t.PkgPath() + ")"
		} else {
			from = t.Name() + " (" + t.PkgPath() + ")"
		}

		var to string
		t = reflect.TypeOf(casted)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			to = "*" + t.Name() + " (" + t.PkgPath() + ")"
		} else {
			to = t.Name() + " (" + t.PkgPath() + ")"
		}
		err = NewErr(FailToCastDaxConn{Name: name, FromType: from, ToType: to})
	}

	return *new(C), err
}
