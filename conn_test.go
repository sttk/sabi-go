package sabi

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

var logs list.List
var willFailToCreateFooConn bool = false
var willFailToCommitFooConn bool = false

type /* error reason */ (
	InvalidConn struct{}
)

func Clear() {
	isGlobalConnCfgSealed = false
	globalConnCfgMap = make(map[string]ConnCfg)

	logs.Init()

	willFailToCreateFooConn = false
	willFailToCommitFooConn = false
}

type FooConn struct {
	Label string
}

func (conn *FooConn) Commit() Err {
	if willFailToCommitFooConn {
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
	if willFailToCreateFooConn {
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
