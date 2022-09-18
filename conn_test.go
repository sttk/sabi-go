package sabi

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

var logs list.List
var WillFailToCreateFooConn bool = false
var WillFailToCommitFooConn bool = false

type /* error reason */ (
	InvalidConn struct{}
)

func Clear() {
	isGlobalConnCfgSealed = false
	globalConnCfgMap = make(map[string]ConnCfg)

	logs.Init()

	WillFailToCreateFooConn = false
	WillFailToCommitFooConn = false
}

type FooConn struct {
	Label string
}

func (conn *FooConn) Commit() Err {
	if WillFailToCommitFooConn {
		return ErrBy(InvalidConn{})
	}
	logs.PushBack("FooConn#Commit")
	return Ok()
}
func (conn *FooConn) Rollback() {
	logs.PushBack("FooConn#Rollback")
}
func (conn *FooConn) Close() {
	logs.PushBack("FooConn#Close")
}

type FooConnCfg struct {
	Label string
}

func (cfg FooConnCfg) CreateConn() (Conn, Err) {
	if WillFailToCreateFooConn {
		return nil, ErrBy(InvalidConn{})
	}
	return &FooConn{Label: cfg.Label}, Ok()
}

type BarConn struct {
	Label string
}

func (conn *BarConn) Commit() Err {
	logs.PushBack("BarConn#Commit")
	return Ok()
}
func (conn *BarConn) Rollback() {
	logs.PushBack("BarConn#Rollback")
}
func (conn *BarConn) Close() {
	logs.PushBack("BarConn#Close")
}

type BarConnCfg struct {
	Label string
}

func (cfg BarConnCfg) CreateConn() (Conn, Err) {
	return &BarConn{Label: cfg.Label}, Ok()
}

func TestAddGlobalConnCfg(t *testing.T) {
	Clear()
	defer Clear()

	assert.False(t, isGlobalConnCfgSealed)
	assert.Equal(t, len(globalConnCfgMap), 0)

	AddGlobalConnCfg("foo", FooConnCfg{})

	assert.False(t, isGlobalConnCfgSealed)
	assert.Equal(t, len(globalConnCfgMap), 1)

	AddGlobalConnCfg("bar", &BarConnCfg{})

	assert.False(t, isGlobalConnCfgSealed)
	assert.Equal(t, len(globalConnCfgMap), 2)
}

func TestSealGlobalConnCfgs(t *testing.T) {
	Clear()
	defer Clear()

	assert.False(t, isGlobalConnCfgSealed)
	assert.Equal(t, len(globalConnCfgMap), 0)

	AddGlobalConnCfg("foo", FooConnCfg{})

	assert.False(t, isGlobalConnCfgSealed)
	assert.Equal(t, len(globalConnCfgMap), 1)

	SealGlobalConnCfgs()

	assert.True(t, isGlobalConnCfgSealed)
	assert.Equal(t, len(globalConnCfgMap), 1)

	AddGlobalConnCfg("foo", FooConnCfg{})

	assert.True(t, isGlobalConnCfgSealed)
	assert.Equal(t, len(globalConnCfgMap), 1)

	isGlobalConnCfgSealed = false

	AddGlobalConnCfg("bar", &BarConnCfg{})

	assert.False(t, isGlobalConnCfgSealed)
	assert.Equal(t, len(globalConnCfgMap), 2)
}

func TestNewConnBase(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 0)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
}

func TestConnBase_addLocalConnCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 0)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)

	base.addLocalConnCfg("foo", FooConnCfg{})

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 1)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)

	base.addLocalConnCfg("bar", &BarConnCfg{})

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 2)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
}

func TestConnBase_begin(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	assert.False(t, isGlobalConnCfgSealed)
	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 0)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)

	base.addLocalConnCfg("foo", FooConnCfg{})

	assert.False(t, isGlobalConnCfgSealed)
	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 1)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)

	base.begin()

	assert.True(t, isGlobalConnCfgSealed)
	assert.True(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 1)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)

	base.addLocalConnCfg("bar", &BarConnCfg{})

	assert.True(t, isGlobalConnCfgSealed)
	assert.True(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 1)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)

	base.isLocalConnCfgSealed = false

	base.addLocalConnCfg("bar", &BarConnCfg{})

	assert.True(t, isGlobalConnCfgSealed)
	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 2)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
}

func TestConnBase_GetConn_withLocalConnCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	conn0, err0 := base.GetConn("foo")
	assert.Nil(t, conn0)
	switch err0.Reason().(type) {
	case ConnCfgIsNotFound:
		assert.Equal(t, err0.Get("Name"), "foo")
	default:
		assert.Fail(t, err0.Error())
	}

	base.addLocalConnCfg("foo", FooConnCfg{})

	conn1, err1 := base.GetConn("foo")
	assert.NotNil(t, conn1)
	assert.True(t, err1.IsOk())

	conn2, err2 := base.GetConn("foo")
	assert.Equal(t, conn2, conn1)
	assert.True(t, err2.IsOk())
}

func TestConnBase_GetConn_withGlobalConnCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	conn0, err0 := base.GetConn("foo")
	assert.Nil(t, conn0)
	switch err0.Reason().(type) {
	case ConnCfgIsNotFound:
		assert.Equal(t, err0.Get("Name"), "foo")
	default:
		assert.Fail(t, err0.Error())
	}

	AddGlobalConnCfg("foo", FooConnCfg{})

	conn1, err1 := base.GetConn("foo")
	assert.NotNil(t, conn1)
	assert.True(t, err1.IsOk())

	conn2, err2 := base.GetConn("foo")
	assert.Equal(t, conn2, conn1)
	assert.True(t, err2.IsOk())
}

func TestConnBase_GetConn_localCfgIsTakenPriorityOfGlobalCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	conn, err := base.GetConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case ConnCfgIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	AddGlobalConnCfg("foo", FooConnCfg{Label: "global"})
	SealGlobalConnCfgs()

	base.addLocalConnCfg("foo", FooConnCfg{Label: "local"})

	conn, err = base.GetConn("foo")
	assert.Equal(t, conn.(*FooConn).Label, "local")
	assert.True(t, err.IsOk())
}

func TestConnBase_GetConn_failToCreateConn(t *testing.T) {
	Clear()
	defer Clear()

	WillFailToCreateFooConn = true
	defer func() { WillFailToCreateFooConn = false }()

	base := NewConnBase()
	base.addLocalConnCfg("foo", FooConnCfg{})

	conn, err := base.GetConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case FailToCreateConn:
		assert.Equal(t, err.Get("Name"), "foo")
		switch err.Cause().(Err).Reason().(type) {
		case InvalidConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}
}

func TestConnBase_InnerMap(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	m := base.InnerMap()
	assert.Nil(t, m["param"])
	m["param"] = 123

	m = base.InnerMap()
	assert.Equal(t, m["param"], 123)
	m["param"] = 456

	m = base.InnerMap()
	assert.Equal(t, m["param"], 456)
}

func TestConnBase_commit(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	base.addLocalConnCfg("foo", FooConnCfg{})
	base.addLocalConnCfg("bar", &BarConnCfg{})
	base.begin()

	fooConn, fooErr := base.GetConn("foo")
	assert.NotNil(t, fooConn)
	assert.True(t, fooErr.IsOk())

	barConn, barErr := base.GetConn("bar")
	assert.NotNil(t, barConn)
	assert.True(t, barErr.IsOk())

	err := base.commit()
	assert.True(t, err.IsOk())

	assert.Equal(t, logs.Len(), 2)
	if logs.Front().Value == "FooConn#Commit" {
		assert.Equal(t, logs.Front().Value, "FooConn#Commit")
		assert.Equal(t, logs.Back().Value, "BarConn#Commit")
	} else {
		assert.Equal(t, logs.Front().Value, "BarConn#Commit")
		assert.Equal(t, logs.Back().Value, "FooConn#Commit")
	}
}

func TestConnBase_commit_failed(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	base.addLocalConnCfg("foo", FooConnCfg{})
	base.addLocalConnCfg("bar", &BarConnCfg{})
	base.begin()

	fooConn, fooErr := base.GetConn("foo")
	assert.NotNil(t, fooConn)
	assert.True(t, fooErr.IsOk())

	barConn, barErr := base.GetConn("bar")
	assert.NotNil(t, barConn)
	assert.True(t, barErr.IsOk())

	WillFailToCommitFooConn = true

	err := base.commit()
	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case FailToCommitConn:
		m := err.Get("Errors").(map[string]Err)
		assert.Equal(t, m["foo"].ReasonName(), "InvalidConn")
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, logs.Len(), 1)
	assert.Equal(t, logs.Back().Value, "BarConn#Commit")
}

func TestConnBase_rollback(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	base.addLocalConnCfg("foo", FooConnCfg{})
	base.addLocalConnCfg("bar", &BarConnCfg{})
	base.begin()

	fooConn, fooErr := base.GetConn("foo")
	assert.NotNil(t, fooConn)
	assert.True(t, fooErr.IsOk())

	barConn, barErr := base.GetConn("bar")
	assert.NotNil(t, barConn)
	assert.True(t, barErr.IsOk())

	base.rollback()

	assert.Equal(t, logs.Len(), 2)
	if logs.Front().Value == "FooConn#Rollback" {
		assert.Equal(t, logs.Front().Value, "FooConn#Rollback")
		assert.Equal(t, logs.Back().Value, "BarConn#Rollback")
	} else {
		assert.Equal(t, logs.Front().Value, "BarConn#Rollback")
		assert.Equal(t, logs.Back().Value, "FooConn#Rollback")
	}
}

func TestConnBase_close(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()

	base.addLocalConnCfg("foo", FooConnCfg{})
	base.addLocalConnCfg("bar", &BarConnCfg{})
	base.begin()

	fooConn, fooErr := base.GetConn("foo")
	assert.NotNil(t, fooConn)
	assert.True(t, fooErr.IsOk())

	barConn, barErr := base.GetConn("bar")
	assert.NotNil(t, barConn)
	assert.True(t, barErr.IsOk())

	base.close()

	assert.Equal(t, logs.Len(), 2)
	if logs.Front().Value == "FooConn#Close" {
		assert.Equal(t, logs.Front().Value, "FooConn#Close")
		assert.Equal(t, logs.Back().Value, "BarConn#Close")
	} else {
		assert.Equal(t, logs.Front().Value, "BarConn#Close")
		assert.Equal(t, logs.Back().Value, "FooConn#Close")
	}
}
