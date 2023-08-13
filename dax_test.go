package sabi

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sttk/sabi/errs"
)

var (
	Logs                       list.List
	WillFailToSetupFooDaxSrc   bool
	WillFailToSetupBarDaxSrc   bool
	WillFailToCreateFooDaxConn bool
	WillCreatedFooDaxConnBeNil bool
	WillFailToCommitFooDaxConn bool
	WillFailToCommitBarDaxConn bool
)

func Reset() {
	errs.FixCfg()

	isGlobalDaxSrcsFixed = false
	globalDaxSrcEntryList.head = nil
	globalDaxSrcEntryList.last = nil

	WillFailToSetupFooDaxSrc = false
	WillFailToSetupBarDaxSrc = false

	WillFailToCreateFooDaxConn = false
	WillCreatedFooDaxConnBeNil = false

	WillFailToCommitFooDaxConn = false
	WillFailToCommitBarDaxConn = false

	Logs.Init()
}

type (
	FailToSetupFooDaxSrc struct{}
	FailToSetupBarDaxSrc struct{}

	FailToCreateFooDaxConn struct{}
	FailToCommitFooDaxConn struct{}
	FailToCommitBarDaxConn struct{}
)

///

type FooDaxSrc struct{}

func (ds FooDaxSrc) Setup(ag AsyncGroup) errs.Err {
	if WillFailToSetupFooDaxSrc {
		return errs.New(FailToSetupFooDaxSrc{})
	}
	Logs.PushBack("FooDaxSrc#Setup")
	return errs.Ok()
}

func (ds FooDaxSrc) Close() {
	Logs.PushBack("FooDaxSrc#Close")
}

func (ds FooDaxSrc) CreateDaxConn() (DaxConn, errs.Err) {
	if WillFailToCreateFooDaxConn {
		return nil, errs.New(FailToCreateFooDaxConn{})
	}
	if WillCreatedFooDaxConnBeNil {
		return nil, errs.Ok()
	}
	Logs.PushBack("FooDaxSrc#CreateDaxConn")
	return FooDaxConn{client: &FooClient{}}, errs.Ok()
}

type FooClient struct {
	committed bool
}
type FooDaxConn struct {
	client *FooClient
}

func (conn FooDaxConn) Commit(ag AsyncGroup) errs.Err {
	if WillFailToCommitFooDaxConn {
		return errs.New(FailToCommitFooDaxConn{})
	}
	Logs.PushBack("FooDaxConn#Commit")
	conn.client.committed = true
	return errs.Ok()
}

func (conn FooDaxConn) IsCommitted() bool {
	return conn.client.committed
}

func (conn FooDaxConn) Rollback(ag AsyncGroup) {
	Logs.PushBack("FooDaxConn#Rollback")
}

func (conn FooDaxConn) ForceBack(ag AsyncGroup) {
	Logs.PushBack("FooDaxConn#ForceBack")
}

func (conn FooDaxConn) Close() {
	Logs.PushBack("FooDaxConn#Close")
}

type BarDaxSrc struct{}

func (ds *BarDaxSrc) Setup(ag AsyncGroup) errs.Err {
	ag.Add(func() errs.Err {
		if WillFailToSetupBarDaxSrc {
			return errs.New(FailToSetupBarDaxSrc{})
		}
		Logs.PushBack("BarDaxSrc#Setup")
		return errs.Ok()
	})
	return errs.Ok()
}
func (ds *BarDaxSrc) Close() {
	Logs.PushBack("BarDaxSrc#Close")
}
func (ds *BarDaxSrc) CreateDaxConn() (DaxConn, errs.Err) {
	Logs.PushBack("BarDaxSrc#CreateDaxConn")
	return &BarDaxConn{}, errs.Ok()
}

type BarDaxConn struct {
	committed bool
}

func (conn *BarDaxConn) Commit(ag AsyncGroup) errs.Err {
	ag.Add(func() errs.Err {
		if WillFailToCommitBarDaxConn {
			return errs.New(FailToCommitBarDaxConn{})
		}
		Logs.PushBack("BarDaxConn#Commit")
		conn.committed = true
		return errs.Ok()
	})
	return errs.Ok()
}
func (conn *BarDaxConn) IsCommitted() bool {
	return conn.committed
}
func (conn *BarDaxConn) Rollback(ag AsyncGroup) {
	Logs.PushBack("BarDaxConn#Rollback")
}
func (conn *BarDaxConn) ForceBack(ag AsyncGroup) {
	Logs.PushBack("BarDaxConn#ForceBack")
}
func (conn *BarDaxConn) Close() {
	Logs.PushBack("BarDaxConn#Close")
}

///

func TestUses_ok(t *testing.T) {
	Reset()
	defer Reset()

	Uses("cliargs", FooDaxSrc{})

	ent0 := globalDaxSrcEntryList.head
	assert.Equal(t, ent0.name, "cliargs")
	assert.IsType(t, ent0.ds, FooDaxSrc{})
	assert.False(t, ent0.local)
	assert.Nil(t, ent0.prev)
	assert.Nil(t, ent0.next)

	Uses("database", &FooDaxSrc{})

	ent0 = globalDaxSrcEntryList.head
	assert.Equal(t, ent0.name, "cliargs")
	assert.IsType(t, ent0.ds, FooDaxSrc{})
	assert.False(t, ent0.local)
	assert.Nil(t, ent0.prev)

	ent1 := ent0.next
	assert.Equal(t, ent1.name, "database")
	assert.IsType(t, ent1.ds, &FooDaxSrc{})
	assert.False(t, ent1.local)
	assert.Equal(t, ent1.prev, ent0)
	assert.Nil(t, ent1.next)

	Uses("file", &BarDaxSrc{})

	ent0 = globalDaxSrcEntryList.head
	assert.Equal(t, ent0.name, "cliargs")
	assert.IsType(t, ent0.ds, FooDaxSrc{})
	assert.False(t, ent0.local)
	assert.Nil(t, ent0.prev)

	ent1 = ent0.next
	assert.Equal(t, ent1.name, "database")
	assert.IsType(t, ent1.ds, &FooDaxSrc{})
	assert.False(t, ent1.local)
	assert.Equal(t, ent1.prev, ent0)

	ent2 := ent1.next
	assert.Equal(t, ent2.name, "file")
	assert.IsType(t, ent2.ds, &BarDaxSrc{})
	assert.False(t, ent2.local)
	assert.Equal(t, ent2.prev, ent1)
}

func TestUses_nameAlreadyExists(t *testing.T) {
	Reset()
	defer Reset()

	Uses("database", FooDaxSrc{})

	ent0 := globalDaxSrcEntryList.head
	assert.Equal(t, ent0.name, "database")
	assert.IsType(t, ent0.ds, FooDaxSrc{})
	assert.False(t, ent0.local)
	assert.Nil(t, ent0.prev)
	assert.Nil(t, ent0.next)

	Uses("database", &FooDaxSrc{})

	ent0 = globalDaxSrcEntryList.head
	assert.Equal(t, ent0.name, "database")
	assert.IsType(t, ent0.ds, FooDaxSrc{})
	assert.False(t, ent0.local)
	assert.Nil(t, ent0.prev)

	ent1 := ent0.next
	assert.Equal(t, ent1.name, "database")
	assert.IsType(t, ent1.ds, &FooDaxSrc{})
	assert.False(t, ent1.local)
	assert.Equal(t, ent1.prev, ent0)
	assert.Nil(t, ent1.next)
}

func TestSetup_zeroDs(t *testing.T) {
	Reset()
	defer Reset()

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Nil(t, Logs.Front())

	err := Setup()
	assert.True(t, err.IsOk())
	defer Close()

	assert.True(t, isGlobalDaxSrcsFixed)
	assert.Nil(t, Logs.Front())
}

func TestSetup_oneDs(t *testing.T) {
	Reset()
	defer Reset()

	Uses("cliargs", FooDaxSrc{})

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Nil(t, Logs.Front())

	err := Setup()
	assert.True(t, err.IsOk())
	defer Close()

	assert.True(t, isGlobalDaxSrcsFixed)
	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Nil(t, log)
}

func TestSetup_multipleDs(t *testing.T) {
	Reset()
	defer Reset()

	Uses("cliargs", FooDaxSrc{})
	Uses("database", &BarDaxSrc{})

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Nil(t, Logs.Front())

	err := Setup()
	assert.True(t, err.IsOk())
	defer Close()

	assert.True(t, isGlobalDaxSrcsFixed)
	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Nil(t, log)
}

func TestSetup_cannotAddAfterSetup(t *testing.T) {
	Reset()
	defer Reset()

	Uses("cliargs", FooDaxSrc{})

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Nil(t, Logs.Front())

	err := Setup()
	assert.True(t, err.IsOk())

	assert.True(t, isGlobalDaxSrcsFixed)
	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Nil(t, log)

	ent := globalDaxSrcEntryList.head
	assert.IsType(t, ent.ds, FooDaxSrc{})
	assert.Nil(t, ent.next)

	Uses("database", &FooDaxSrc{})

	assert.True(t, isGlobalDaxSrcsFixed)
	log = Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Nil(t, log)

	ent = globalDaxSrcEntryList.head
	assert.IsType(t, ent.ds, FooDaxSrc{})
	assert.Nil(t, ent.next)
}

func TestSetup_error_sync(t *testing.T) {
	Reset()
	defer Reset()

	WillFailToSetupFooDaxSrc = true

	Uses("cliargs", FooDaxSrc{})

	err := Setup()
	assert.True(t, err.IsNotOk())
	assert.IsType(t, err.Reason(), FailToSetupGlobalDaxSrcs{})
	errmap := err.Reason().(FailToSetupGlobalDaxSrcs).Errors
	assert.Equal(t, len(errmap), 1)
	assert.IsType(t, errmap["cliargs"].Reason(), FailToSetupFooDaxSrc{})

	assert.True(t, isGlobalDaxSrcsFixed)

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestSetup_error_async(t *testing.T) {
	Reset()
	defer Reset()

	WillFailToSetupBarDaxSrc = true

	Uses("cliargs", &BarDaxSrc{})

	err := Setup()
	assert.True(t, err.IsNotOk())
	assert.IsType(t, err.Reason(), FailToSetupGlobalDaxSrcs{})
	errmap := err.Reason().(FailToSetupGlobalDaxSrcs).Errors
	assert.Equal(t, len(errmap), 1)
	assert.IsType(t, errmap["cliargs"].Reason(), FailToSetupBarDaxSrc{})

	assert.True(t, isGlobalDaxSrcsFixed)
	log := Logs.Front()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestSetup_error_asyncAndSync(t *testing.T) {
	Reset()
	defer Reset()

	WillFailToSetupBarDaxSrc = true
	WillFailToSetupFooDaxSrc = true

	Uses("cliargs", &BarDaxSrc{})
	Uses("database", FooDaxSrc{})

	err := Setup()
	assert.True(t, err.IsNotOk())
	assert.IsType(t, err.Reason(), FailToSetupGlobalDaxSrcs{})
	errmap := err.Reason().(FailToSetupGlobalDaxSrcs).Errors
	assert.Equal(t, len(errmap), 2)
	assert.IsType(t, errmap["cliargs"].Reason(), FailToSetupBarDaxSrc{})
	assert.IsType(t, errmap["database"].Reason(), FailToSetupFooDaxSrc{})

	assert.True(t, isGlobalDaxSrcsFixed)
	log := Logs.Front()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestStartApp_ok(t *testing.T) {
	Reset()
	defer Reset()

	Uses("database", FooDaxSrc{})

	app := func() errs.Err {
		Logs.PushBack("run app")
		return errs.Ok()
	}

	err := StartApp(app)
	assert.True(t, err.IsOk())

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "run app")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestStartApp_error_failToSetup(t *testing.T) {
	Reset()
	defer Reset()

	WillFailToSetupFooDaxSrc = true

	Uses("database", FooDaxSrc{})

	app := func() errs.Err {
		Logs.PushBack("run app")
		return errs.Ok()
	}

	err := StartApp(app)
	assert.True(t, err.IsNotOk())
	assert.IsType(t, err.Reason(), FailToSetupGlobalDaxSrcs{})
	errmap := err.Reason().(FailToSetupGlobalDaxSrcs).Errors
	assert.Equal(t, len(errmap), 1)
	assert.IsType(t, errmap["database"].Reason(), FailToSetupFooDaxSrc{})

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestStartApp_error_appReturnsError(t *testing.T) {
	Reset()
	defer Reset()

	Uses("database", FooDaxSrc{})

	type FailToDoSomething struct{}

	app := func() errs.Err {
		return errs.New(FailToDoSomething{})
	}

	err := StartApp(app)
	assert.True(t, err.IsNotOk())
	assert.IsType(t, err.Reason(), FailToDoSomething{})

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestNewDaxBase_withNoGlobalDs(t *testing.T) {
	Reset()
	defer Reset()

	assert.False(t, isGlobalDaxSrcsFixed)

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	assert.True(t, isGlobalDaxSrcsFixed)

	assert.False(t, base.isLocalDaxSrcsFixed)
	assert.Nil(t, base.localDaxSrcEntryList.head)
	assert.Nil(t, base.localDaxSrcEntryList.last)
	assert.Equal(t, len(base.daxSrcEntryMap), 0)
	assert.Equal(t, base.daxConnMap.Len(), 0)

	log := Logs.Front()
	assert.Nil(t, log)
}

func TestNewDaxBase_withSomeGlobalDs(t *testing.T) {
	Reset()
	defer Reset()

	assert.False(t, isGlobalDaxSrcsFixed)

	Uses("cliargs", FooDaxSrc{})
	Uses("database", &BarDaxSrc{})

	func() {
		base := NewDaxBase().(*daxBaseImpl)
		defer base.Close()

		assert.True(t, isGlobalDaxSrcsFixed)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.IsType(t, base.daxSrcEntryMap["cliargs"].ds, FooDaxSrc{})
		assert.IsType(t, base.daxSrcEntryMap["database"].ds, &BarDaxSrc{})
		assert.Equal(t, base.daxConnMap.Len(), 0)
	}()

	log := Logs.Front()
	assert.Nil(t, log)
}

func TestDax_Uses_ok(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		err := base.Uses("cliargs", FooDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 1)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		ent := base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})

		err = base.Uses("database", &BarDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		ent = ent.next
		assert.Nil(t, ent)

		base.Close()

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.True(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.True(t, base.daxSrcEntryMap["database"].deleted)
		assert.Equal(t, base.daxConnMap.Len(), 0)
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDax_Uses_doNothingWhileFixed(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)
		defer base.Close()

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		base.begin()

		err := base.Uses("cliargs", FooDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		base.end()
	}()

	log := Logs.Front()
	assert.Nil(t, log)
}

func TestDax_Uses_failToSetupLocalDaxSrc_sync(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)
		defer base.Close()

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		WillFailToSetupFooDaxSrc = true

		err := base.Uses("cliargs", FooDaxSrc{})
		assert.True(t, err.IsNotOk())
		switch r := err.Reason().(type) {
		case FailToSetupLocalDaxSrc:
			assert.Equal(t, r.Name, "cliargs")
		default:
			assert.Fail(t, err.Error())
		}

		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		ent := base.localDaxSrcEntryList.head
		assert.Nil(t, ent)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)
	}()

	log := Logs.Front()
	assert.Nil(t, log)
}

func TestDax_Uses_failToSetupLocalDaxSrc_async(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)
		defer base.Close()

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		WillFailToSetupBarDaxSrc = true

		err := base.Uses("database", &BarDaxSrc{})
		assert.True(t, err.IsNotOk())
		switch r := err.Reason().(type) {
		case FailToSetupLocalDaxSrc:
			assert.Equal(t, r.Name, "database")
		default:
			assert.Fail(t, err.Error())
		}

		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		ent := base.localDaxSrcEntryList.head
		assert.Nil(t, ent)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)
	}()

	log := Logs.Front()
	assert.Nil(t, log)
}

func TestDax_Uses_createRunner(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		runner := base.Uses_("cliargs", FooDaxSrc{})
		err := runner()
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 1)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		ent := base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})

		runner = base.Uses_("database", &BarDaxSrc{})
		err = runner()
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		ent = ent.next
		assert.Nil(t, ent)

		base.Close()

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.True(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.True(t, base.daxSrcEntryMap["database"].deleted)
		assert.Equal(t, base.daxConnMap.Len(), 0)
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDax_Disuses_ok(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		err := base.Uses("cliargs", FooDaxSrc{})
		assert.True(t, err.IsOk())

		err = base.Uses("database", &BarDaxSrc{})
		assert.True(t, err.IsOk())

		err = base.Uses("file", &FooDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 3)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		assert.False(t, base.daxSrcEntryMap["file"].deleted)
		ent := base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		base.Disuses("database")

		assert.Equal(t, len(base.daxSrcEntryMap), 3)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.True(t, base.daxSrcEntryMap["database"].deleted)
		assert.False(t, base.daxSrcEntryMap["file"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		base.Disuses("cliargs")

		assert.Equal(t, len(base.daxSrcEntryMap), 3)
		assert.True(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.True(t, base.daxSrcEntryMap["database"].deleted)
		assert.False(t, base.daxSrcEntryMap["file"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, &FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		base.Disuses("file")

		assert.Equal(t, len(base.daxSrcEntryMap), 3)
		assert.True(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.True(t, base.daxSrcEntryMap["database"].deleted)
		assert.True(t, base.daxSrcEntryMap["file"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.Nil(t, ent)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, base.daxConnMap.Len(), 0)
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDax_Disuses_doNothingWhileFixed(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		err := base.Uses("cliargs", FooDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 1)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		ent := base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)

		err = base.Uses("database", &BarDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		assert.False(t, base.isLocalDaxSrcsFixed)
		base.begin()
		assert.True(t, base.isLocalDaxSrcsFixed)

		base.Disuses("cliargs")

		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		base.Disuses("database")

		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDax_Disuses_createRunner(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		err := base.Uses("cliargs", FooDaxSrc{})
		assert.True(t, err.IsOk())

		err = base.Uses("database", &BarDaxSrc{})
		assert.True(t, err.IsOk())

		err = base.Uses("file", &FooDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 3)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		assert.False(t, base.daxSrcEntryMap["file"].deleted)
		ent := base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		runner := base.Disuses_("database")
		err = runner()
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 3)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.True(t, base.daxSrcEntryMap["database"].deleted)
		assert.False(t, base.daxSrcEntryMap["file"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		runner = base.Disuses_("cliargs")
		err = runner()
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 3)
		assert.True(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.True(t, base.daxSrcEntryMap["database"].deleted)
		assert.False(t, base.daxSrcEntryMap["file"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, &FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		runner = base.Disuses_("file")
		err = runner()
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 3)
		assert.True(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.True(t, base.daxSrcEntryMap["database"].deleted)
		assert.True(t, base.daxSrcEntryMap["file"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.Nil(t, ent)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, base.daxConnMap.Len(), 0)
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDax_Close_doNothingWhileFixed(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase().(*daxBaseImpl)

		assert.False(t, base.isLocalDaxSrcsFixed)
		assert.Nil(t, base.localDaxSrcEntryList.head)
		assert.Nil(t, base.localDaxSrcEntryList.last)
		assert.Equal(t, len(base.daxSrcEntryMap), 0)
		assert.Equal(t, base.daxConnMap.Len(), 0)

		err := base.Uses("cliargs", FooDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 1)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		ent := base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)

		err = base.Uses("database", &BarDaxSrc{})
		assert.True(t, err.IsOk())

		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)

		assert.False(t, base.isLocalDaxSrcsFixed)
		base.begin()
		assert.True(t, base.isLocalDaxSrcsFixed)

		base.Close()

		assert.Equal(t, len(base.daxSrcEntryMap), 2)
		assert.False(t, base.daxSrcEntryMap["cliargs"].deleted)
		assert.False(t, base.daxSrcEntryMap["database"].deleted)
		ent = base.localDaxSrcEntryList.head
		assert.IsType(t, ent.ds, FooDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.IsType(t, ent.ds, &BarDaxSrc{})
		assert.False(t, ent.deleted)
		ent = ent.next
		assert.Nil(t, ent)
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Nil(t, log)
}

func TestGetDaxConn_ok(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	err := base.Uses("cliargs", FooDaxSrc{})
	assert.True(t, err.IsOk())
	err = base.Uses("database", &BarDaxSrc{})
	assert.True(t, err.IsOk())
	err = base.Uses("file", &FooDaxSrc{})
	assert.True(t, err.IsOk())

	base.begin()
	defer base.end()

	conn1, err := GetDaxConn[FooDaxConn](base, "cliargs")
	assert.True(t, err.IsOk())
	assert.IsType(t, conn1, FooDaxConn{})

	conn2, err := GetDaxConn[*BarDaxConn](base, "database")
	assert.True(t, err.IsOk())
	assert.IsType(t, conn2, &BarDaxConn{})

	conn3, err := GetDaxConn[FooDaxConn](base, "file")
	assert.True(t, err.IsOk())
	assert.IsType(t, conn3, FooDaxConn{})
}

func TestGetDaxConn_daxConnIsAlreadyCreated(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	err := base.Uses("cliargs", FooDaxSrc{})
	assert.True(t, err.IsOk())

	base.begin()
	defer base.end()

	conn1, err := GetDaxConn[FooDaxConn](base, "cliargs")
	assert.True(t, err.IsOk())
	assert.IsType(t, conn1, FooDaxConn{})

	conn2, err := GetDaxConn[FooDaxConn](base, "cliargs")
	assert.True(t, err.IsOk())
	assert.IsType(t, conn2, FooDaxConn{})
	assert.Equal(t, &conn1, &conn2)
}

func TestGetDaxConn_daxSrcIsNotFound(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	base.begin()
	defer base.end()

	_, err := GetDaxConn[FooDaxConn](base, "cliargs")
	switch r := err.Reason().(type) {
	case DaxSrcIsNotFound:
		assert.Equal(t, r.Name, "cliargs")
	default:
		assert.Fail(t, err.Error())
	}
}

func TestGetDaxConn_daxSrcIsDisused(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	err := base.Uses("cliargs", FooDaxSrc{})
	assert.True(t, err.IsOk())

	func() {
		base.begin()
		defer base.end()

		conn, err := GetDaxConn[FooDaxConn](base, "cliargs")
		assert.True(t, err.IsOk())
		assert.IsType(t, conn, FooDaxConn{})
	}()

	base.Disuses("cliargs")

	func() {
		base.begin()
		defer base.end()

		_, err := GetDaxConn[FooDaxConn](base, "cliargs")
		switch r := err.Reason().(type) {
		case DaxSrcIsNotFound:
			assert.Equal(t, r.Name, "cliargs")
		default:
			assert.Fail(t, err.Error())
		}
	}()
}

func TestGetDaxConn_localDaxSrcIsDisusedButGlobalDaxSrcExists(t *testing.T) {
	Reset()
	defer Reset()

	Uses("database", &BarDaxSrc{})

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	err := base.Uses("database", FooDaxSrc{})
	assert.True(t, err.IsOk())

	func() {
		base.begin()
		defer base.end()

		conn, err := GetDaxConn[FooDaxConn](base, "database")
		assert.True(t, err.IsOk())
		assert.IsType(t, conn, FooDaxConn{})
	}()

	base.Disuses("database")

	func() {
		base.begin()
		defer base.end()

		conn, err := GetDaxConn[*BarDaxConn](base, "database")
		assert.True(t, err.IsOk())
		assert.IsType(t, conn, &BarDaxConn{})
	}()
}

func TestGetDaxConn_failToCreateDaxConn(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	err := base.Uses("database", FooDaxSrc{})
	assert.True(t, err.IsOk())

	WillFailToCreateFooDaxConn = true

	func() {
		base.begin()
		defer base.end()

		_, err := GetDaxConn[FooDaxConn](base, "database")
		switch r := err.Reason().(type) {
		case FailToCreateDaxConn:
			assert.Equal(t, r.Name, "database")
			assert.IsType(t, err.Cause().(errs.Err).Reason(), FailToCreateFooDaxConn{})
		default:
			assert.Fail(t, err.Error())
		}
	}()
}

func TestGetDaxConn_createdDaxConnIsNil(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	err := base.Uses("database", FooDaxSrc{})
	assert.True(t, err.IsOk())

	WillCreatedFooDaxConnBeNil = true

	func() {
		base.begin()
		defer base.end()

		_, err := GetDaxConn[FooDaxConn](base, "database")
		switch r := err.Reason().(type) {
		case CreatedDaxConnIsNil:
			assert.Equal(t, r.Name, "database")
		default:
			assert.Fail(t, err.Error())
		}
	}()
}

func TestGetDaxConn_failToCastDaxConn(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase().(*daxBaseImpl)
	defer base.Close()

	err := base.Uses("database", FooDaxSrc{})
	assert.True(t, err.IsOk())

	func() {
		base.begin()
		defer base.end()

		_, err := GetDaxConn[*BarDaxConn](base, "database")
		switch r := err.Reason().(type) {
		case FailToCastDaxConn:
			assert.Equal(t, r.Name, "database")
			assert.Equal(t, r.FromType, "sabi.FooDaxConn")
			assert.Equal(t, r.ToType, "*sabi.BarDaxConn")
		default:
			assert.Fail(t, err.Error())
		}
	}()
}

func TestTxn_zeroLogic(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase()
	defer base.Close()

	err := Txn[Dax](base)
	assert.True(t, err.IsOk())
}

func TestTxn_oneLogic(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase()
		defer base.Close()

		err := base.Uses("database", FooDaxSrc{})
		assert.True(t, err.IsOk())

		err = Txn(base, func(dax Dax) errs.Err {
			_, err := GetDaxConn[FooDaxConn](dax, "database")
			assert.True(t, err.IsOk())
			return errs.Ok()
		})
		assert.True(t, err.IsOk())
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Commit")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestTxn_twoLogic(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase()
		defer base.Close()

		err := base.Uses("database", FooDaxSrc{})
		assert.True(t, err.IsOk())
		err = err.IfOk(base.Uses_("file", &BarDaxSrc{}))
		assert.True(t, err.IsOk())

		err = Txn(base, func(dax Dax) errs.Err {
			_, err := GetDaxConn[FooDaxConn](dax, "database")
			assert.True(t, err.IsOk())
			return errs.Ok()
		}, func(dax Dax) errs.Err {
			_, err := GetDaxConn[*BarDaxConn](dax, "file")
			assert.True(t, err.IsOk())
			return errs.Ok()
		})
		assert.True(t, err.IsOk())
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Commit")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Commit")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

type LogicDax interface {
	Dax
	getData() string
}

func TestTxn_failToCastDaxBase(t *testing.T) {
	Reset()
	defer Reset()

	func() {
		base := NewDaxBase()
		defer base.Close()

		err := Txn(base, func(dax LogicDax) errs.Err {
			return errs.Ok()
		})
		switch r := err.Reason().(type) {
		case FailToCastDaxBase:
			assert.Equal(t, r.FromType, "sabi.DaxBase")
			assert.Equal(t, r.ToType, "sabi.LogicDax")
		default:
			assert.Fail(t, err.Error())
		}
	}()
}

func TestTxn_failToRunLogic(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase()
	defer base.Close()

	type FailToDoSomething struct{}

	func() {
		base := NewDaxBase()
		defer base.Close()

		err := base.Uses("database", FooDaxSrc{})
		assert.True(t, err.IsOk())
		err = err.IfOk(base.Uses_("file", &BarDaxSrc{}))
		assert.True(t, err.IsOk())

		err = Txn(base, func(dax Dax) errs.Err {
			_, err := GetDaxConn[FooDaxConn](dax, "database")
			assert.True(t, err.IsOk())
			Logs.PushBack("run logic 1")
			return errs.New(FailToDoSomething{})
		}, func(dax Dax) errs.Err {
			_, err := GetDaxConn[*BarDaxConn](dax, "file")
			assert.True(t, err.IsOk())
			Logs.PushBack("run logic 2")
			return errs.Ok()
		})
		switch err.Reason().(type) {
		case FailToDoSomething:
		default:
			assert.Fail(t, err.Error())
		}
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "run logic 1")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Rollback")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestTxn_failToCommit_sync(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase()
	defer base.Close()

	func() {
		base := NewDaxBase()
		defer base.Close()

		err := base.Uses("database", FooDaxSrc{})
		assert.True(t, err.IsOk())
		err = base.Uses("file", &BarDaxSrc{})
		assert.True(t, err.IsOk())

		WillFailToCommitFooDaxConn = true

		err = Txn(base, func(dax Dax) errs.Err {
			_, err := GetDaxConn[FooDaxConn](dax, "database")
			assert.True(t, err.IsOk())
			Logs.PushBack("run logic 1")
			return errs.Ok()
		}, func(dax Dax) errs.Err {
			_, err := GetDaxConn[*BarDaxConn](dax, "file")
			assert.True(t, err.IsOk())
			Logs.PushBack("run logic 2")
			return errs.Ok()
		})
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "run logic 1")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "run logic 2")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Rollback")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Rollback")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestTxn_failToCommit_async(t *testing.T) {
	Reset()
	defer Reset()

	base := NewDaxBase()
	defer base.Close()

	func() {
		base := NewDaxBase()
		defer base.Close()

		err := base.Uses("database", FooDaxSrc{})
		assert.True(t, err.IsOk())
		err = base.Uses("file", &BarDaxSrc{})
		assert.True(t, err.IsOk())

		WillFailToCommitBarDaxConn = true

		err = Txn(base, func(dax Dax) errs.Err {
			_, err := GetDaxConn[FooDaxConn](dax, "database")
			assert.True(t, err.IsOk())
			Logs.PushBack("run logic 1")
			return errs.Ok()
		}, func(dax Dax) errs.Err {
			_, err := GetDaxConn[*BarDaxConn](dax, "file")
			assert.True(t, err.IsOk())
			Logs.PushBack("run logic 2")
			return errs.Ok()
		})
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "run logic 1")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "run logic 2")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Commit")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#ForceBack")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Rollback")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}

func TestTxn_runner(t *testing.T) {
	Reset()
	defer Reset()

	Reset()
	defer Reset()

	func() {
		base := NewDaxBase()
		defer base.Close()

		err := base.Uses("database", FooDaxSrc{}).
			IfOk(base.Uses_("file", &BarDaxSrc{})).
			IfOk(Txn_(base, func(dax Dax) errs.Err {
				_, err := GetDaxConn[FooDaxConn](dax, "database")
				assert.True(t, err.IsOk())
				return errs.Ok()
			}, func(dax Dax) errs.Err {
				_, err := GetDaxConn[*BarDaxConn](dax, "file")
				assert.True(t, err.IsOk())
				return errs.Ok()
			}))

		assert.True(t, err.IsOk())
	}()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Setup")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Commit")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Commit")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#Close")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#Close")
	log = log.Next()
	assert.Nil(t, log)
}
