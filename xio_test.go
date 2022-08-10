package sabi

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func ClearGlobalXioConnCfgs() {
	isGlobalXioConnCfgSealed = false
	globalXioConnCfgMap = make(map[string]XioConnCfg)
}

type /* error reasons */ (
	InvalidXioConn struct{}
)

var willFailToCreateXioConn bool = false

type FooXioConn struct {
	Label string
}

func (conn *FooXioConn) Commit() Err {
	return Ok()
}

func (conn *FooXioConn) Rollback() Err {
	return Ok()
}

func (conn *FooXioConn) Close() Err {
	return Ok()
}

type FooXio struct {
	Xio
}

func NewFooXio(xio Xio) FooXio {
	return FooXio{Xio: xio}
}

func (xio FooXio) GetFooConn(name string) (*FooXioConn, Err) {
	conn, err := xio.GetConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*FooXioConn), Ok()
}

type FooXioConnCfg struct {
	Label string
}

func (cfg FooXioConnCfg) NewConn() (XioConn, Err) {
	if willFailToCreateXioConn {
		return nil, ErrBy(InvalidXioConn{})
	}
	return &FooXioConn{Label: cfg.Label}, Ok()
}

type BarXioConn struct {
}

func (conn *BarXioConn) Commit() Err {
	return Ok()
}

func (conn *BarXioConn) Rollback() Err {
	return Ok()
}

func (conn *BarXioConn) Close() Err {
	return Ok()
}

type BarXio struct {
	Xio
}

func NewBarXio(xio Xio) BarXio {
	return BarXio{Xio: xio}
}

func (xio BarXio) GetBarConn(name string) (*BarXioConn, Err) {
	conn, err := xio.GetConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*BarXioConn), Ok()
}

type BarXioConnCfg struct {
	Label string
}

func (cfg *BarXioConnCfg) NewConn() (XioConn, Err) {
	return &BarXioConn{}, Ok()
}

func TestAddGlobalXioConnCfg(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	assert.False(t, isGlobalXioConnCfgSealed)
	assert.Equal(t, len(globalXioConnCfgMap), 0)

	AddGlobalXioConnCfg("foo", FooXioConnCfg{})

	assert.False(t, isGlobalXioConnCfgSealed)
	assert.Equal(t, len(globalXioConnCfgMap), 1)

	AddGlobalXioConnCfg("bar", &BarXioConnCfg{})

	assert.False(t, isGlobalXioConnCfgSealed)
	assert.Equal(t, len(globalXioConnCfgMap), 2)
}

func TestSealGlobalXioConnCfgs(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	assert.False(t, isGlobalXioConnCfgSealed)
	assert.Equal(t, len(globalXioConnCfgMap), 0)

	AddGlobalXioConnCfg("foo", FooXioConnCfg{})

	assert.False(t, isGlobalXioConnCfgSealed)
	assert.Equal(t, len(globalXioConnCfgMap), 1)

	SealGlobalXioConnCfgs()

	assert.True(t, isGlobalXioConnCfgSealed)
	assert.Equal(t, len(globalXioConnCfgMap), 1)

	AddGlobalXioConnCfg("foo", FooXioConnCfg{})

	assert.True(t, isGlobalXioConnCfgSealed)
	assert.Equal(t, len(globalXioConnCfgMap), 1)

	isGlobalXioConnCfgSealed = false

	AddGlobalXioConnCfg("bar", &BarXioConnCfg{})

	assert.False(t, isGlobalXioConnCfgSealed)
	assert.Equal(t, len(globalXioConnCfgMap), 2)
}

func TestNewXioBase(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	base := NewXioBase()

	assert.False(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 0)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())
}

func TestXioBase_AddLocalXioConnCfg(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	base := NewXioBase()

	assert.False(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 0)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.AddLocalXioConnCfg("foo-local", FooXioConnCfg{})

	assert.False(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 1)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.AddLocalXioConnCfg("bar-local", &BarXioConnCfg{})

	assert.False(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 2)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())
}

func TestXioBase_SealLocalXioConnCfg(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	base := NewXioBase()

	assert.False(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 0)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.AddLocalXioConnCfg("foo-local", FooXioConnCfg{})

	assert.False(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 1)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.SealLocalXioConnCfgs()

	assert.True(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 1)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.AddLocalXioConnCfg("bar-local", &BarXioConnCfg{})

	assert.True(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 1)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())

	base.isLocalXioConnCfgSealed = false

	base.AddLocalXioConnCfg("bar-local", &BarXioConnCfg{})

	assert.False(t, base.isLocalXioConnCfgSealed)
	assert.Equal(t, len(base.localXioConnCfgMap), 2)
	assert.Equal(t, len(base.xioConnMap), 0)
	assert.Equal(t, len(base.innerMap), 0)
	assert.True(t, base.error.IsOk())
}

func TestXioBase_GetConn_withLocalXioConnCfg(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	base := NewXioBase()

	conn0, err0 := base.GetConn("foo")
	assert.Nil(t, conn0)
	switch err0.Reason().(type) {
	case XioConnCfgIsNotFound:
		assert.Equal(t, err0.Get("Name"), "foo")
	default:
		assert.Fail(t, err0.Error())
	}

	base.AddLocalXioConnCfg("foo", FooXioConnCfg{})

	conn1, err1 := base.GetConn("foo")
	assert.NotNil(t, conn1)
	assert.True(t, err1.IsOk())

	conn2, err2 := base.GetConn("foo")
	assert.Equal(t, conn2, conn1)
	assert.True(t, err2.IsOk())
}

func TestXioBase_GetConn_withGlobalXioConnCfg(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	base := NewXioBase()

	conn0, err0 := base.GetConn("foo")
	assert.Nil(t, conn0)
	switch err0.Reason().(type) {
	case XioConnCfgIsNotFound:
		assert.Equal(t, err0.Get("Name"), "foo")
	default:
		assert.Fail(t, err0.Error())
	}

	AddGlobalXioConnCfg("foo", FooXioConnCfg{})
	SealGlobalXioConnCfgs()

	conn1, err1 := base.GetConn("foo")
	assert.NotNil(t, conn1)
	assert.True(t, err1.IsOk())

	conn2, err2 := base.GetConn("foo")
	assert.Equal(t, conn2, conn1)
	assert.True(t, err2.IsOk())
}

func TestXioBase_GetConn_localCfgIsTakenPriorityOfGlobalCfg(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	base := NewXioBase()

	conn, err := base.GetConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case XioConnCfgIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	AddGlobalXioConnCfg("foo", FooXioConnCfg{Label: "global"})
	SealGlobalXioConnCfgs()

	base.AddLocalXioConnCfg("foo", FooXioConnCfg{Label: "local"})

	conn, err = base.GetConn("foo")
	assert.Equal(t, conn.(*FooXioConn).Label, "local")
	assert.True(t, err.IsOk())
}

func TestXioBase_GetConn_failToCreateConn(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	willFailToCreateXioConn = true
	defer func() { willFailToCreateXioConn = false }()

	base := NewXioBase()
	base.AddLocalXioConnCfg("foo", FooXioConnCfg{})

	conn, err := base.GetConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case FailToCreateXioConn:
		assert.Equal(t, err.Get("Name"), "foo")
		switch err.Cause().(Err).Reason().(type) {
		case InvalidXioConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}
}

func TestXioBase_GetConn_ForEachDataSrc(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

	base := NewXioBase()
	base.AddLocalXioConnCfg("foo", FooXioConnCfg{})
	base.AddLocalXioConnCfg("bar", &BarXioConnCfg{})
	base.SealLocalXioConnCfgs()

	fooXio := NewFooXio(base)
	fooConn, err0 := fooXio.GetFooConn("foo")
	assert.True(t, err0.IsOk())
	assert.Equal(t, reflect.TypeOf(fooConn).String(), "*sabi.FooXioConn")

	barXio := NewBarXio(base)
	barConn, err1 := barXio.GetBarConn("bar")
	assert.True(t, err1.IsOk())
	assert.Equal(t, reflect.TypeOf(barConn).String(), "*sabi.BarXioConn")
}

func TestXioBase_InnerMap(t *testing.T) {
	ClearGlobalXioConnCfgs()
	defer ClearGlobalXioConnCfgs()

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
