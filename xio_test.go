package sabi

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type FooXio struct {
	Xio
}

func NewFooXio(xio Xio) FooXio {
	return FooXio{Xio: xio}
}

func (xio FooXio) GetFooConn(name string) (*FooConn, Err) {
	conn, err := xio.GetConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*FooConn), Ok()
}

type BarXio struct {
	Xio
}

func NewBarXio(xio Xio) BarXio {
	return BarXio{Xio: xio}
}

func (xio BarXio) GetBarConn(name string) (*BarConn, Err) {
	conn, err := xio.GetConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*BarConn), Ok()
}

func TestNewXioBase(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 0)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())
}

func TestXioBase_AddLocalConnCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 0)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.AddLocalConnCfg("foo-local", FooConnCfg{})

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 1)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.AddLocalConnCfg("bar-local", &BarConnCfg{})

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 2)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())
}

func TestXioBase_SealLocalConnCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 0)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.AddLocalConnCfg("foo-local", FooConnCfg{})

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 1)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.SealLocalConnCfgs()

	assert.True(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 1)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.AddLocalConnCfg("bar-local", &BarConnCfg{})

	assert.True(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 1)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.isLocalConnCfgSealed = false

	base.AddLocalConnCfg("bar-local", &BarConnCfg{})

	assert.False(t, base.isLocalConnCfgSealed)
	assert.Equal(t, len(base.localConnCfgMap), 2)
	assert.Equal(t, len(base.connMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())
}

func TestXioBase_GetConn_withLocalConnCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	conn0, err0 := base.GetConn("foo")
	assert.Nil(t, conn0)
	switch err0.Reason().(type) {
	case ConnCfgIsNotFound:
		assert.Equal(t, err0.Get("Name"), "foo")
	default:
		assert.Fail(t, err0.Error())
	}

	base.AddLocalConnCfg("foo", FooConnCfg{})

	conn1, err1 := base.GetConn("foo")
	assert.NotNil(t, conn1)
	assert.True(t, err1.IsOk())

	conn2, err2 := base.GetConn("foo")
	assert.Equal(t, conn2, conn1)
	assert.True(t, err2.IsOk())
}

func TestXioBase_GetConn_withGlobalConnCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	conn0, err0 := base.GetConn("foo")
	assert.Nil(t, conn0)
	switch err0.Reason().(type) {
	case ConnCfgIsNotFound:
		assert.Equal(t, err0.Get("Name"), "foo")
	default:
		assert.Fail(t, err0.Error())
	}

	AddGlobalConnCfg("foo", FooConnCfg{})
	SealGlobalConnCfgs()

	conn1, err1 := base.GetConn("foo")
	assert.NotNil(t, conn1)
	assert.True(t, err1.IsOk())

	conn2, err2 := base.GetConn("foo")
	assert.Equal(t, conn2, conn1)
	assert.True(t, err2.IsOk())
}

func TestXioBase_GetConn_localCfgIsTakenPriorityOfGlobalCfg(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

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

	base.AddLocalConnCfg("foo", FooConnCfg{Label: "local"})

	conn, err = base.GetConn("foo")
	assert.Equal(t, conn.(*FooConn).Label, "local")
	assert.True(t, err.IsOk())
}

func TestXioBase_GetConn_failToCreateConn(t *testing.T) {
	Clear()
	defer Clear()

	willFailToCreateFooConn = true
	defer func() { willFailToCreateFooConn = false }()

	base := NewXioBase()
	base.AddLocalConnCfg("foo", FooConnCfg{})

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

func TestXioBase_GetConn_ForEachDataSrc(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()
	base.AddLocalConnCfg("foo", FooConnCfg{})
	base.AddLocalConnCfg("bar", &BarConnCfg{})
	base.SealLocalConnCfgs()

	fooXio := NewFooXio(base)
	fooConn, err0 := fooXio.GetFooConn("foo")
	assert.True(t, err0.IsOk())
	assert.Equal(t, reflect.TypeOf(fooConn).String(), "*sabi.FooConn")

	barXio := NewBarXio(base)
	barConn, err1 := barXio.GetBarConn("bar")
	assert.True(t, err1.IsOk())
	assert.Equal(t, reflect.TypeOf(barConn).String(), "*sabi.BarConn")
}

func TestXioBase_InnerMap(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	m := base.InnerMap()
	assert.Nil(t, m["param"])
	m["param"] = 123

	m = base.InnerMap()
	assert.Equal(t, m["param"], 123)
	m["param"] = 456

	m = base.InnerMap()
	assert.Equal(t, m["param"], 456)
}

func TestXioBase_commit(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	base.AddLocalConnCfg("foo", FooConnCfg{})
	base.AddLocalConnCfg("bar", &BarConnCfg{})
	base.SealLocalConnCfgs()

	fooXio := NewFooXio(base)
	fooConn, _ := fooXio.GetFooConn("foo")
	assert.NotNil(t, fooConn)

	barXio := NewBarXio(base)
	barConn, _ := barXio.GetBarConn("bar")
	assert.NotNil(t, barConn)

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

func TestXioBase_commit_failed(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	base.AddLocalConnCfg("foo", FooConnCfg{})
	base.AddLocalConnCfg("bar", &BarConnCfg{})
	base.SealLocalConnCfgs()

	fooXio := NewFooXio(base)
	fooConn, _ := fooXio.GetFooConn("foo")
	assert.NotNil(t, fooConn)

	barXio := NewBarXio(base)
	barConn, _ := barXio.GetBarConn("bar")
	assert.NotNil(t, barConn)

	willFailToCommitFooConn = true

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

func TestXioBase_Rollback(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	base.AddLocalConnCfg("foo", FooConnCfg{})
	base.AddLocalConnCfg("bar", &BarConnCfg{})
	base.SealLocalConnCfgs()

	fooXio := NewFooXio(base)
	fooConn, _ := fooXio.GetFooConn("foo")
	assert.NotNil(t, fooConn)

	barXio := NewBarXio(base)
	barConn, _ := barXio.GetBarConn("bar")
	assert.NotNil(t, barConn)

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

func TestXioBase_Close(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()

	base.AddLocalConnCfg("foo", FooConnCfg{})
	base.AddLocalConnCfg("bar", &BarConnCfg{})
	base.SealLocalConnCfgs()

	fooXio := NewFooXio(base)
	fooConn, _ := fooXio.GetFooConn("foo")
	assert.NotNil(t, fooConn)

	barXio := NewBarXio(base)
	barConn, _ := barXio.GetBarConn("bar")
	assert.NotNil(t, barConn)

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
