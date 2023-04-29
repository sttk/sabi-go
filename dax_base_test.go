package sabi

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var logs list.List
var WillFailToSetUpFooDaxSrc bool = false
var WillFailToCreateFooDaxConn bool = false
var WillFailToCommitFooDaxConn bool = false

type /* error reason */ (
	InvalidDaxConn struct{}
)

func ClearDaxBase() {
	isGlobalDaxSrcsFixed = false
	globalDaxSrcMap = make(map[string]DaxSrc)

	logs.Init()

	WillFailToSetUpFooDaxSrc = false
	WillFailToCreateFooDaxConn = false
	WillFailToCommitFooDaxConn = false
}

type FooDaxConn struct {
	Label string
}

func (conn *FooDaxConn) Commit() Err {
	if WillFailToCommitFooDaxConn {
		return NewErr(InvalidDaxConn{})
	}
	logs.PushBack("FooDaxConn#Commit")
	return Ok()
}

func (conn *FooDaxConn) Rollback() {
	logs.PushBack("FooDaxConn#Rollback")
}

func (conn *FooDaxConn) Close() {
	logs.PushBack("FooDaxConn#Close")
}

type FooDaxSrc struct {
	Label string
}

func (ds FooDaxSrc) CreateDaxConn() (DaxConn, Err) {
	if WillFailToCreateFooDaxConn {
		return nil, NewErr(InvalidDaxConn{})
	}
	logs.PushBack("FooDaxSrc#CreateDaxConn")
	return &FooDaxConn{Label: ds.Label}, Ok()
}

func (ds FooDaxSrc) SetUp() Err {
	if WillFailToSetUpFooDaxSrc {
		return NewErr(InvalidDaxConn{})
	}
	logs.PushBack("FooDaxSrc#SetUp")
	return Ok()
}

func (ds FooDaxSrc) End() {
	logs.PushBack("FooDaxSrc#End")
}

type BarDaxConn struct {
	Label string
	store map[string]string
}

func (conn *BarDaxConn) Commit() Err {
	logs.PushBack("BarDaxConn#Commit")
	return Ok()
}

func (conn *BarDaxConn) Rollback() {
	logs.PushBack("BarDaxConn#Rollback")
}

func (conn *BarDaxConn) Close() {
	logs.PushBack("BarDaxConn#Close")
}

func (conn *BarDaxConn) Store(name, value string) {
	conn.store[name] = value
}

type BarDaxSrc struct {
	Label string
	Store map[string]string
}

func (ds BarDaxSrc) CreateDaxConn() (DaxConn, Err) {
	logs.PushBack("BarDaxSrc#CreateDaxConn")
	return &BarDaxConn{Label: ds.Label, store: ds.Store}, Ok()
}

func (ds BarDaxSrc) SetUp() Err {
	logs.PushBack("BarDaxSrc#SetUp")
	return Ok()
}

func (ds BarDaxSrc) End() {
	logs.PushBack("BarDaxSrc#End")
}

func TestAddGlobalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 0)

	AddGlobalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)

	isGlobalDaxSrcsFixed = true

	AddGlobalDaxSrc("bar", &BarDaxSrc{})

	assert.True(t, isGlobalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)

	isGlobalDaxSrcsFixed = false

	AddGlobalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 2)
}

func TestStartUpGlobalDaxSrcs_and_ShutdownGlobalDaxSrcs(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 0)

	AddGlobalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)

	AddGlobalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 2)

	if err := StartUpGlobalDaxSrcs(); !err.IsOk() {
		t.Logf("err = %v\n", err)
		return
	}
	assert.True(t, isGlobalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 2)

	elem := logs.Front()
	if elem.Value == "FooDaxSrc#SetUp" {
		assert.Equal(t, elem.Value, "FooDaxSrc#SetUp")
		assert.Equal(t, elem.Next().Value, "BarDaxSrc#SetUp")
	} else {
		assert.Equal(t, elem.Value, "BarDaxSrc#SetUp")
		assert.Equal(t, elem.Next().Value, "FooDaxSrc#SetUp")
	}
	assert.Nil(t, elem.Next().Next())

	defer func() {
		ShutdownGlobalDaxSrcs()

		elem = logs.Front()
		if elem.Value == "FooDaxSrc#SetUp" {
			assert.Equal(t, elem.Value, "FooDaxSrc#SetUp")
			assert.Equal(t, elem.Next().Value, "BarDaxSrc#SetUp")
		} else {
			assert.Equal(t, elem.Value, "BarDaxSrc#SetUp")
			assert.Equal(t, elem.Next().Value, "FooDaxSrc#SetUp")
		}
		elem = elem.Next().Next()
		if elem.Value == "FooDaxSrc#End" {
			assert.Equal(t, elem.Value, "FooDaxSrc#End")
			assert.Equal(t, elem.Next().Value, "BarDaxSrc#End")
		} else {
			assert.Equal(t, elem.Value, "BarDaxSrc#End")
			assert.Equal(t, elem.Next().Value, "FooDaxSrc#End")
		}
		assert.Nil(t, elem.Next().Next())
	}()
}

func TestDaxBase_SetUpLocalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 0)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	err := base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.True(t, err.IsOk())

	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	base.(*daxBaseImpl).isLocalDaxSrcsFixed = true

	err = base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})
	assert.True(t, err.IsOk())

	assert.True(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	base.(*daxBaseImpl).isLocalDaxSrcsFixed = false

	err = base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})
	assert.True(t, err.IsOk())

	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 2)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)
}

func TestDaxBase_FreeLocalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 0)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	err := base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.True(t, err.IsOk())

	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	err = base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})
	assert.True(t, err.IsOk())

	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 2)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	base.FreeLocalDaxSrc("bar")

	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	base.(*daxBaseImpl).isLocalDaxSrcsFixed = true

	base.FreeLocalDaxSrc("foo")

	assert.True(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	base.(*daxBaseImpl).isLocalDaxSrcsFixed = false

	base.FreeLocalDaxSrc("foo")

	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 0)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)
}

func TestDaxBase_begin(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 0)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 0)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	AddGlobalDaxSrc("foo", FooDaxSrc{})
	err := base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.True(t, err.IsOk())

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	base.begin()

	assert.True(t, isGlobalDaxSrcsFixed)
	assert.True(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	AddGlobalDaxSrc("bar", &BarDaxSrc{})
	err = base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})
	assert.True(t, err.IsOk())

	assert.True(t, isGlobalDaxSrcsFixed)
	assert.True(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	base.(*daxBaseImpl).isLocalDaxSrcsFixed = false

	assert.True(t, isGlobalDaxSrcsFixed)
	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	AddGlobalDaxSrc("bar", &BarDaxSrc{})
	err = base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})
	assert.True(t, err.IsOk())

	assert.True(t, isGlobalDaxSrcsFixed)
	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 2)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	isGlobalDaxSrcsFixed = false

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 1)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 2)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)

	AddGlobalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, isGlobalDaxSrcsFixed)
	assert.False(t, base.(*daxBaseImpl).isLocalDaxSrcsFixed)
	assert.Equal(t, len(globalDaxSrcMap), 2)
	assert.Equal(t, len(base.(*daxBaseImpl).localDaxSrcMap), 2)
	assert.Equal(t, len(base.(*daxBaseImpl).daxConnMap), 0)
}

func TestDaxBase_GetDaxConn_withLocalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	conn, err := base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case DaxSrcIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	err = base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.True(t, err.IsOk())

	conn, err = base.GetDaxConn("foo")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	var conn2 DaxConn
	conn2, err = base.GetDaxConn("foo")
	assert.Equal(t, conn2, conn)
	assert.True(t, err.IsOk())
}

func TestDaxBase_GetDaxConn_withGlobalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	conn, err := base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case DaxSrcIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	AddGlobalDaxSrc("foo", FooDaxSrc{})

	conn, err = base.GetDaxConn("foo")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	var conn2 DaxConn
	conn2, err = base.GetDaxConn("foo")
	assert.Equal(t, conn2, conn)
	assert.True(t, err.IsOk())
}

func TestDaxBase_GetDaxConn_localDsIsTakenPriorityOfGlobalDs(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	conn, err := base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case DaxSrcIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	AddGlobalDaxSrc("foo", FooDaxSrc{Label: "global"})
	ShutdownGlobalDaxSrcs()

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{Label: "local"})

	conn, err = base.GetDaxConn("foo")
	assert.Equal(t, conn.(*FooDaxConn).Label, "local")
	assert.True(t, err.IsOk())
}

func TestDaxBase_GetDaxConn_failToCreateDaxConn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	WillFailToCreateFooDaxConn = true
	defer func() { WillFailToCreateFooDaxConn = false }()

	base := NewDaxBase()
	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})

	conn, err := base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case FailToCreateDaxConn:
		assert.Equal(t, err.Get("Name"), "foo")
		switch err.Cause().(Err).Reason().(type) {
		case InvalidDaxConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}
}

func TestDaxBase_commit(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})
	base.begin()

	fooConn, fooErr := base.GetDaxConn("foo")
	assert.NotNil(t, fooConn)
	assert.True(t, fooErr.IsOk())

	barConn, barErr := base.GetDaxConn("bar")
	assert.NotNil(t, barConn)
	assert.True(t, barErr.IsOk())

	err := base.commit()
	assert.True(t, err.IsOk())

	assert.Equal(t, logs.Len(), 6)
	elem := logs.Front()
	assert.Equal(t, elem.Value, "FooDaxSrc#SetUp")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxSrc#SetUp")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "FooDaxSrc#CreateDaxConn")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxSrc#CreateDaxConn")
	elem = elem.Next()
	if elem.Value == "FooDaxConn#Commit" {
		assert.Equal(t, elem.Value, "FooDaxConn#Commit")
		elem = elem.Next()
		assert.Equal(t, elem.Value, "BarDaxConn#Commit")
	} else {
		assert.Equal(t, elem.Value, "BarDaxConn#Commit")
		elem = elem.Next()
		assert.Equal(t, elem.Value, "FooDaxConn#Commit")
	}
	elem = elem.Next()
	assert.Nil(t, elem)
}

func TestDaxBase_commit_failed(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})

	base.begin()

	fooConn, fooErr := base.GetDaxConn("foo")
	assert.NotNil(t, fooConn)
	assert.True(t, fooErr.IsOk())

	barConn, barErr := base.GetDaxConn("bar")
	assert.NotNil(t, barConn)
	assert.True(t, barErr.IsOk())

	WillFailToCommitFooDaxConn = true

	err := base.commit()
	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case FailToCommitDaxConn:
		m := err.Get("Errors").(map[string]Err)
		assert.Equal(t, m["foo"].ReasonName(), "InvalidDaxConn")
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, logs.Len(), 5)
	elem := logs.Front()
	assert.Equal(t, elem.Value, "FooDaxSrc#SetUp")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxSrc#SetUp")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "FooDaxSrc#CreateDaxConn")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxSrc#CreateDaxConn")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxConn#Commit")
	elem = elem.Next()
	assert.Nil(t, elem)
}

func TestDaxBase_rollback(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})
	base.begin()

	fooConn, fooErr := base.GetDaxConn("foo")
	assert.NotNil(t, fooConn)
	assert.True(t, fooErr.IsOk())

	barConn, barErr := base.GetDaxConn("bar")
	assert.NotNil(t, barConn)
	assert.True(t, barErr.IsOk())

	base.rollback()

	assert.Equal(t, logs.Len(), 6)
	elem := logs.Front()
	assert.Equal(t, elem.Value, "FooDaxSrc#SetUp")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxSrc#SetUp")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "FooDaxSrc#CreateDaxConn")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxSrc#CreateDaxConn")
	elem = elem.Next()
	if elem.Value == "FooDaxConn#Rollback" {
		assert.Equal(t, elem.Value, "FooDaxConn#Rollback")
		elem = elem.Next()
		assert.Equal(t, elem.Value, "BarDaxConn#Rollback")
	} else {
		assert.Equal(t, elem.Value, "BarDaxConn#Rollback")
		elem = elem.Next()
		assert.Equal(t, elem.Value, "FooDaxConn#Rollback")
	}
	elem = elem.Next()
	assert.Nil(t, elem)
}

func TestDaxBase_close(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})
	base.begin()

	fooConn, fooErr := base.GetDaxConn("foo")
	assert.NotNil(t, fooConn)
	assert.True(t, fooErr.IsOk())

	barConn, barErr := base.GetDaxConn("bar")
	assert.NotNil(t, barConn)
	assert.True(t, barErr.IsOk())

	base.end()

	assert.Equal(t, logs.Len(), 6)
	elem := logs.Front()
	assert.Equal(t, elem.Value, "FooDaxSrc#SetUp")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxSrc#SetUp")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "FooDaxSrc#CreateDaxConn")
	elem = elem.Next()
	assert.Equal(t, elem.Value, "BarDaxSrc#CreateDaxConn")
	elem = elem.Next()
	if elem.Value == "FooDaxConn#Close" {
		assert.Equal(t, elem.Value, "FooDaxConn#Close")
		elem = elem.Next()
		assert.Equal(t, elem.Value, "BarDaxConn#Close")
	} else {
		assert.Equal(t, elem.Value, "BarDaxConn#Close")
		elem = elem.Next()
		assert.Equal(t, elem.Value, "FooDaxConn#Close")
	}
	elem = elem.Next()
	assert.Nil(t, elem)
}

type FooDax struct {
	Dax
}

func NewFooDax(dax Dax) FooDax {
	return FooDax{Dax: dax}
}

func (dax FooDax) GetFooDaxConn(name string) (*FooDaxConn, Err) {
	conn, err := dax.GetDaxConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*FooDaxConn), Ok()
}

type BarDax struct {
	Dax
}

func NewBarDax(dax Dax) BarDax {
	return BarDax{Dax: dax}
}

func (dax BarDax) GetBarDaxConn(name string) (*BarDaxConn, Err) {
	conn, err := dax.GetDaxConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*BarDaxConn), Ok()
}

func TestDax_GetXxxConn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := NewDaxBase()
	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})

	base.begin()

	fooDax := NewFooDax(base)
	fooConn, fooErr := fooDax.GetFooDaxConn("foo")
	assert.True(t, fooErr.IsOk())
	assert.Equal(t, reflect.TypeOf(fooConn).String(), "*sabi.FooDaxConn")

	barDax := NewBarDax(base)
	barConn, barErr := barDax.GetBarDaxConn("bar")
	assert.True(t, barErr.IsOk())
	assert.Equal(t, reflect.TypeOf(barConn).String(), "*sabi.BarDaxConn")
}
