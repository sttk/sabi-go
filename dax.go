// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"

	om "github.com/sttk/orderedmap"

	"github.com/sttk/sabi/errs"
)

type /* error reasons */ (
	// FailToSetupGlobalDaxSrcs is the error reason which indicates that some
	// DaxSrc(s) failed to set up.
	// The field Errors is the map of which keys are the registered names of
	// DaxSrc(s) that failed, and of which values are errs.Err having their error
	// reasons.
	FailToSetupGlobalDaxSrcs struct {
		Errors map[string]errs.Err
	}

	// FailToSetupLocalDaxSrc is the error reason which indicates that a local
	// DaxSrc failed to set up.
	// The field Name is the registered name of the DaxSrc failed.
	FailToSetupLocalDaxSrc struct {
		Name string
	}

	// DaxSrcIsNotFound is the error reason which indicates that a specified
	// DaxSrc is not found.
	// The field Name is the registered name of the DaxSrc not found.
	DaxSrcIsNotFound struct {
		Name string
	}

	// FailToCreateDaxConn is the error reason which indicates that it is failed
	// to create a new connection to a data store.
	// The field Name is the registered name of the DaxSrc failed to create a
	// DaxConn.
	FailToCreateDaxConn struct {
		Name string
	}

	// FailToCommitDaxConn is the error reason which indicates that some
	// connections failed to commit.
	// The field Errors is the map of which keys are registered names of DaxConn
	// which failed to commit, and of which values are errs.Err(s) having their
	// error reasons.
	FailToCommitDaxConn struct {
		Errors map[string]errs.Err
	}

	// CreatedDaxConnIsNil is the error reason which indicates that a DaxSrc
	// created a DaxConn instance but it is nil.
	// The field Name is the registered name of the DaxSrc that created a nil
	// DaxConn.
	CreatedDaxConnIsNil struct {
		Name string
	}

	// FailToCastDaxConn is the error reason which indicates that a DaxConn
	// instance of the specified name failed to cast to the destination type.
	// The field Name and FromType is the name and type name of the DaxConn
	// instance, and the field ToType is the type name of the destination type.
	FailToCastDaxConn struct {
		Name, FromType, ToType string
	}

	// FailToCastDaxBase is the error reason which indicates that a DaxBase instance
	// failed to cast to the destination type.
	// The field FromType is the type name of the DaxBase instance and the field
	// ToType is the type name of the destination type.
	FailToCastDaxBase struct {
		FromType, ToType string
	}
)

// DaxConn is the interface that represents a connection to a data store, and
// defines methods: Commit, Rollback and Close to work in a transaction
// process.
//
// Commit is the method for commiting updates in a transaction.
// IsCommitted is the method for check whether updates are already committed.
// Rollback is the method for rollback updates in a transaction.
// ForceBack is the method to revert updates forcely even if updates are
// already commited or this connection ooes not have rollback mechanism.
// If commting and rollbacking procedures are asynchronous, the argument
// AsyncGroup(s) are used to process them.
// Close is the method to close this connecton.
type DaxConn interface {
	Commit(ag AsyncGroup) errs.Err
	IsCommitted() bool
	Rollback(ag AsyncGroup)
	ForceBack(ag AsyncGroup)
	Close()
}

// DaxSrc is the interface that represents a data source which creates
// connections to a data store like database, etc.
// This interface declares three methods: Setup, Close, and CreateDaxConn.
//
// Setup is the method to connect to a data store and to prepare to create
// DaxConn objects which represents a connection to access data in the data
// store.
// If the set up procedure is asynchronous, the Setup method is implemented
// so as to use AsyncGroup.
// Close is the method to disconnect to a data store.
// CreateDaxConn is the method to create a DaxConn object.
type DaxSrc interface {
	Setup(ag AsyncGroup) errs.Err
	Close()
	CreateDaxConn() (DaxConn, errs.Err)
}

type daxSrcEntry struct {
	name    string
	ds      DaxSrc
	local   bool
	deleted bool
	prev    *daxSrcEntry
	next    *daxSrcEntry
}

type daxSrcEntryList struct {
	head *daxSrcEntry
	last *daxSrcEntry
}

var (
	isGlobalDaxSrcsFixed  bool = false
	globalDaxSrcEntryList daxSrcEntryList
)

// Uses is the method that registers a global DaxSrc with its name to enable to
// use DaxConn created by the argument DaxSrc in all transactions.
//
// If a DaxSrc is tried to register with a name already registered, it is
// ignored and a DaxSrc registered with the same name first is used.
// And this method ignore adding new DaxSrc(s) after Setup or first beginning
// of Proc or Txn.
func Uses(name string, ds DaxSrc) errs.Err {
	if isGlobalDaxSrcsFixed {
		return errs.Ok()
	}

	ent := &daxSrcEntry{name: name, ds: ds}

	if globalDaxSrcEntryList.head == nil {
		globalDaxSrcEntryList.head = ent
		globalDaxSrcEntryList.last = ent
	} else {
		ent.prev = globalDaxSrcEntryList.last
		globalDaxSrcEntryList.last.next = ent
		globalDaxSrcEntryList.last = ent
	}

	return errs.Ok()
}

// Setup is the function that make the globally registered DaxSrc usable.
//
// This function forbids adding more global DaxSrc(s), and calls each Setup
// method of all registered DaxSrc(s).
// If one of DaxSrc(s) fails to execute synchronous Setup, this function stops
// other setting up and returns an errs.Err containing the error reason of
// that failure.
//
// If one of DaxSrc(s) fails to execute asynchronous Setup, this function
// continue to other setting up and returns an errs.Err containing the error
// reason of that failure and other errors if any.
func Setup() errs.Err {
	isGlobalDaxSrcsFixed = true
	errs.FixCfg()

	var ag asyncGroupAsync[string]

	for ent := globalDaxSrcEntryList.head; ent != nil; ent = ent.next {
		ag.name = ent.name
		err := ent.ds.Setup(&ag)
		if err.IsNotOk() {
			ag.wait()
			Close()
			ag.addErr(ag.name, err)
			return errs.New(FailToSetupGlobalDaxSrcs{Errors: ag.makeErrs()})
		}
	}

	ag.wait()

	if ag.hasErr() {
		Close()
		return errs.New(FailToSetupGlobalDaxSrcs{Errors: ag.makeErrs()})
	}

	return errs.Ok()
}

// Close is the function that closes and frees each resource of registered
// global DaxSrc(s).
// This function should always be called before an application ends.
func Close() {
	for ent := globalDaxSrcEntryList.head; ent != nil; ent = ent.next {
		ent.ds.Close()
	}
}

// StartApp is the function that calls Setup function, the argument function,
// and Close function in order.
// If Setup function or the argument function fails, this function stops
// calling other functions and return an errs.Err containing the error
// reaason..
//
// This function is a macro-like function aimed at reducing the amount of
// coding.
func StartApp(app func() errs.Err) errs.Err {
	err := Setup()
	if err.IsNotOk() {
		return err
	}
	defer Close()

	return app()
}

// Dax is the interface for a set of data access methods.
//
// This interface is embedded by Dax implementations for data
// stores, and each Dax implementation defines data access methods to each
// data store.
// In data access methods, DacConn instances connected to data stores can be
// got with GetConn function.
//
// At the same time, this interface is embedded by Dax interfaces for logics.
// The each Dax interface declares only methods used by each logic.
type Dax interface {
	getDaxConn(name string) (DaxConn, errs.Err)
}

// DaxBase is the interface that declares the methods to manage DaxSrc(s).
// And this interface declarees unexported methods to process a transaction.
//
// Close is the method to close and free all local DaxSrc(s).
// Uses is the method to register and setup a local DaxSrc with an argument
// name.
// Uses_ is the method that creates a runner function which runs #Uses method.
// Disuses is the method to close and remove a local DaxSrc specified by
// an argument name.
// Disuses_ is the method that creates a runner function which runs #Disuses
// method.
type DaxBase interface {
	Dax

	Close()
	Uses(name string, ds DaxSrc) errs.Err
	Uses_(name string, ds DaxSrc) func() errs.Err
	Disuses(name string)
	Disuses_(name string) func() errs.Err

	begin()
	commit() errs.Err
	rollback()
	end()
}

type daxBaseImpl struct {
	DaxBase

	isLocalDaxSrcsFixed  bool
	localDaxSrcEntryList daxSrcEntryList

	daxSrcEntryMap map[string]*daxSrcEntry
	agSync         asyncGroupSync

	daxConnMap   om.Map[string, DaxConn]
	daxConnMutex sync.Mutex
}

// NewDaxBase is the function that creates a new DaxBase instance.
func NewDaxBase() DaxBase {
	isGlobalDaxSrcsFixed = true
	errs.FixCfg()

	base := &daxBaseImpl{
		daxSrcEntryMap: make(map[string]*daxSrcEntry),
		daxConnMap:     om.New[string, DaxConn](),
	}

	for ent := globalDaxSrcEntryList.last; ent != nil; ent = ent.prev {
		base.daxSrcEntryMap[ent.name] = ent
	}

	return base
}

func (base *daxBaseImpl) Close() {
	if base.isLocalDaxSrcsFixed {
		return
	}

	for ent := base.localDaxSrcEntryList.head; ent != nil; ent = ent.next {
		if !ent.deleted {
			ent.deleted = true
			ent.ds.Close()
		}
	}

	base.localDaxSrcEntryList.head = nil
	base.localDaxSrcEntryList.last = nil
}

func (base *daxBaseImpl) Uses(name string, ds DaxSrc) errs.Err {
	if base.isLocalDaxSrcsFixed {
		return errs.Ok()
	}

	err := ds.Setup(&(base.agSync))
	if err.IsNotOk() {
		return errs.New(FailToSetupLocalDaxSrc{Name: name}, err)
	}

	if base.agSync.err.IsNotOk() {
		return errs.New(FailToSetupLocalDaxSrc{Name: name}, base.agSync.err)
	}

	ent := &daxSrcEntry{name: name, ds: ds, local: true}

	if base.localDaxSrcEntryList.head == nil {
		base.localDaxSrcEntryList.head = ent
		base.localDaxSrcEntryList.last = ent
	} else {
		ent.prev = base.localDaxSrcEntryList.last
		base.localDaxSrcEntryList.last.next = ent
		base.localDaxSrcEntryList.last = ent
	}

	base.daxSrcEntryMap[ent.name] = ent

	return errs.Ok()
}

func (base *daxBaseImpl) Uses_(name string, ds DaxSrc) func() errs.Err {
	return func() errs.Err {
		return base.Uses(name, ds)
	}
}

func (base *daxBaseImpl) Disuses(name string) {
	if base.isLocalDaxSrcsFixed {
		return
	}

	ent := base.daxSrcEntryMap[name]
	if ent != nil && ent.local {
		ent.deleted = true

		if ent.prev != nil {
			ent.prev.next = ent.next
		} else {
			base.localDaxSrcEntryList.head = ent.next
		}

		if ent.next != nil {
			ent.next.prev = ent.prev
		} else {
			base.localDaxSrcEntryList.last = ent.prev
		}

		ent.ds.Close()
	}
}

func (base *daxBaseImpl) Disuses_(name string) func() errs.Err {
	return func() errs.Err {
		base.Disuses(name)
		return errs.Ok()
	}
}

func (base *daxBaseImpl) begin() {
	base.isLocalDaxSrcsFixed = true
}

func (base *daxBaseImpl) commit() errs.Err {
	var ag asyncGroupAsync[string]

	for ent := base.daxConnMap.Front(); ent != nil; ent = ent.Next() {
		ag.name = ent.Key()
		err := ent.Value().Commit(&ag)
		if err.IsNotOk() {
			ag.wait()
			ag.addErr(ent.Key(), err)
			return errs.New(FailToCommitDaxConn{Errors: ag.makeErrs()})
		}
	}

	ag.wait()

	if ag.hasErr() {
		return errs.New(FailToCommitDaxConn{Errors: ag.makeErrs()})
	}

	return errs.Ok()
}

func (base *daxBaseImpl) rollback() {
	var ag asyncGroupAsync[string]

	for ent := base.daxConnMap.Front(); ent != nil; ent = ent.Next() {
		conn := ent.Value()
		if conn.IsCommitted() {
			ent.Value().ForceBack(&ag)
		} else {
			ent.Value().Rollback(&ag)
		}
	}

	ag.wait()
}

func (base *daxBaseImpl) end() {
	for {
		ent := base.daxConnMap.FrontAndLdelete()
		if ent == nil {
			break
		}
		ent.Value().Close()
	}

	base.isLocalDaxSrcsFixed = false
}

func (base *daxBaseImpl) getDaxConn(name string) (DaxConn, errs.Err) {
	conn, loaded := base.daxConnMap.Load(name)
	if loaded {
		return conn, errs.Ok()
	}

	base.daxConnMutex.Lock()
	defer base.daxConnMutex.Unlock()

	conn, _, e := base.daxConnMap.LoadOrStoreFunc(name, func() (DaxConn, error) {
		ent, exists := base.daxSrcEntryMap[name]
		if !exists {
			return nil, errs.New(DaxSrcIsNotFound{Name: name})
		}

		if ent.deleted && ent.local {
			for gEnt := globalDaxSrcEntryList.head; gEnt != nil; gEnt = gEnt.next {
				if gEnt.name == name {
					base.daxSrcEntryMap[ent.name] = gEnt
					ent = gEnt
					break
				}
			}
			if ent.deleted {
				return nil, errs.New(DaxSrcIsNotFound{Name: name})
			}
		}

		conn, err := ent.ds.CreateDaxConn()
		if err.IsNotOk() {
			return nil, errs.New(FailToCreateDaxConn{Name: name}, err)
		}
		if conn == nil {
			return nil, errs.New(CreatedDaxConnIsNil{Name: name})
		}
		return conn, nil
	})

	if e != nil {
		return nil, e.(errs.Err)
	}
	return conn, errs.Ok()
}

// GetDaxConn is the function to cast type of DaxConn instance.
// If the cast failed, this function returns an errs.Err of the reason:
// FailToCastDaxConn with the DaxConn name and type names of source and
// destination.
func GetDaxConn[C DaxConn](dax Dax, name string) (C, errs.Err) {
	conn, err := dax.getDaxConn(name)
	if err.IsOk() {
		casted, ok := conn.(C)
		if ok {
			return casted, err
		}

		from := typeNameOf(conn)
		to := typeNameOfTypeParam[C]()
		err = errs.New(FailToCastDaxConn{Name: name, FromType: from, ToType: to})
	}

	return *new(C), err
}

// Txn is the function that executes logic functions in a transaction.
//
// First, this function casts the argument DaxBase to the type specified with
// the function's type parameter.
// Next, this function begins the transaction, and the argument logic functions
// are executed..
// Then, if no error occurs, this function commits all updates in the
// transaction, otherwise rollbacks them.
// If there are commit errors after some DaxConn(s) are commited, or there are
// DaxConn(s) which don't have rollback mechanism, this function executes
// ForceBack methods of those DaxConn(s).
// And after that, this function ends the transaction.
//
// During a transaction, it is denied to add or remove any local DaxSrc(s).
func Txn[D any](base DaxBase, logics ...func(dax D) errs.Err) errs.Err {
	dax, ok := base.(D)
	if !ok {
		from := typeNameOf(&base)[1:]
		to := typeNameOfTypeParam[D]()
		return errs.New(FailToCastDaxBase{FromType: from, ToType: to})
	}

	base.begin()
	defer base.end()

	err := errs.Ok()

	for _, logic := range logics {
		err = logic(dax)
		if err.IsNotOk() {
			break
		}
	}

	if err.IsOk() {
		err = base.commit()
	}

	if err.IsNotOk() {
		base.rollback()
	}

	return err
}

// Txn_ is the function that creates a runner function which runs a Txn
// function.
func Txn_[D any](base DaxBase, logics ...func(dax D) errs.Err) func() errs.Err {
	return func() errs.Err {
		return Txn[D](base, logics...)
	}
}
