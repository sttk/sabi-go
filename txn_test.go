package sabi_test

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"strings"
	"testing"
)

var (
	willFailToCreateBDaxConn = false
	willFailToCommitBDaxConn = false
)

var logs list.List

type (
	FailToCreateBDaxConn struct{}
	FailToCommitBDaxConn struct{}
	FailToRunLogic       struct{}
)

func clear() {
	willFailToCreateBDaxConn = false
	willFailToCommitBDaxConn = false
	sabi.ClearDaxBase()
	logs.Init()
}

////

type ADaxSrc struct {
	AMap map[string]string
}

func NewADaxSrc() ADaxSrc {
	return ADaxSrc{AMap: make(map[string]string)}
}

func (ds ADaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return &ADaxConn{AMap: ds.AMap}, sabi.Ok()
}

func (ds ADaxSrc) StartUp() sabi.Err {
	return sabi.Ok()
}

func (ds ADaxSrc) Shutdown() {
}

type ADaxConn struct {
	AMap map[string]string
}

func (conn *ADaxConn) Commit() sabi.Err {
	logs.PushBack("ADaxConn#Commit")
	return sabi.Ok()
}
func (conn *ADaxConn) Rollback() {
	logs.PushBack("ADaxConn#Rollback")
}
func (conn *ADaxConn) Close() {
	logs.PushBack("ADaxConn#Close")
}

type ADax struct {
	sabi.Dax
}

func (dax ADax) GetADaxConn(name string) (*ADaxConn, sabi.Err) {
	conn, err := dax.GetDaxConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*ADaxConn), err
}

type BDaxSrc struct {
	BMap map[string]string
}

func NewBDaxSrc() BDaxSrc {
	return BDaxSrc{BMap: make(map[string]string)}
}

func (ds BDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	if willFailToCreateBDaxConn {
		return nil, sabi.NewErr(FailToCreateBDaxConn{})
	}
	return &BDaxConn{BMap: ds.BMap}, sabi.Ok()
}

func (ds BDaxSrc) StartUp() sabi.Err {
	return sabi.Ok()
}

func (ds BDaxSrc) Shutdown() {
}

type BDaxConn struct {
	BMap map[string]string
}

func (conn *BDaxConn) Commit() sabi.Err {
	if willFailToCommitBDaxConn {
		return sabi.NewErr(FailToCommitBDaxConn{})
	}
	logs.PushBack("BDaxConn#Commit")
	return sabi.Ok()
}
func (conn *BDaxConn) Rollback() {
	logs.PushBack("BDaxConn#Rollback")
}
func (conn *BDaxConn) Close() {
	logs.PushBack("BDaxConn#Close")
}

type BDax struct {
	sabi.Dax
}

func (dax BDax) GetBDaxConn(name string) (*BDaxConn, sabi.Err) {
	conn, err := dax.GetDaxConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*BDaxConn), err
}

type CDaxSrc struct {
	CMap map[string]string
}

func NewCDaxSrc() CDaxSrc {
	return CDaxSrc{CMap: make(map[string]string)}
}

func (ds CDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return &CDaxConn{CMap: ds.CMap}, sabi.Ok()
}

func (ds CDaxSrc) StartUp() sabi.Err {
	return sabi.Ok()
}

func (ds CDaxSrc) Shutdown() {
}

type CDaxConn struct {
	CMap map[string]string
}

func (conn *CDaxConn) Commit() sabi.Err {
	logs.PushBack("CDaxConn#Commit")
	return sabi.Ok()
}
func (conn *CDaxConn) Rollback() {
	logs.PushBack("CDaxConn#Rollback")
}
func (conn *CDaxConn) Close() {
	logs.PushBack("CDaxConn#Close")
}

type CDax struct {
	sabi.Dax
}

func (dax CDax) GetCDaxConn(name string) (*CDaxConn, sabi.Err) {
	conn, err := dax.GetDaxConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*CDaxConn), err
}

////

type AGetDax struct {
	ADax
}

func (dax AGetDax) GetAData() (string, sabi.Err) {
	conn, err := dax.GetADaxConn("aaa")
	if !err.IsOk() {
		return "", err
	}
	data := conn.AMap["a"]
	return data, sabi.Ok()
}

type BGetSetDax struct {
	BDax
}

func (dax BGetSetDax) GetBData() (string, sabi.Err) {
	conn, err := dax.GetBDaxConn("bbb")
	if !err.IsOk() {
		return "", err
	}
	data := conn.BMap["b"]
	return data, sabi.Ok()
}
func (dax BGetSetDax) SetBData(data string) sabi.Err {
	conn, err := dax.GetBDaxConn("bbb")
	if !err.IsOk() {
		return err
	}
	conn.BMap["b"] = data
	return sabi.Ok()
}

type CSetDax struct {
	CDax
}

func (dax CSetDax) SetCData(data string) sabi.Err {
	conn, err := dax.GetCDaxConn("ccc")
	if !err.IsOk() {
		return err
	}
	conn.CMap["c"] = data
	return sabi.Ok()
}

////

func TestRunTxn(t *testing.T) {
	clear()
	defer clear()

	type ABDax interface {
		GetAData() (string, sabi.Err)
		SetBData(data string) sabi.Err
	}

	base := sabi.NewDaxBase()
	base = struct {
		sabi.DaxBase
		AGetDax
		BGetSetDax
	}{
		DaxBase:    base,
		AGetDax:    AGetDax{ADax: ADax{Dax: base}},
		BGetSetDax: BGetSetDax{BDax: BDax{Dax: base}},
	}

	aDs := NewADaxSrc()
	bDs := NewBDaxSrc()
	base.AddLocalDaxSrc("aaa", aDs)
	base.AddLocalDaxSrc("bbb", bDs)

	aDs.AMap["a"] = "hello"
	err := sabi.RunTxn(base, func(dax ABDax) sabi.Err {
		data, err := dax.GetAData()
		if !err.IsOk() {
			return err
		}
		data = strings.ToUpper(data)
		err = dax.SetBData(data)
		return err
	})
	assert.True(t, err.IsOk())
	assert.Equal(t, bDs.BMap["b"], "HELLO")

	elem := logs.Front()
	if elem.Value == "ADaxConn#Commit" {
		assert.Equal(t, elem.Value, "ADaxConn#Commit")
		assert.Equal(t, elem.Next().Value, "BDaxConn#Commit")
	} else {
		assert.Equal(t, elem.Value, "BDaxConn#Commit")
		assert.Equal(t, elem.Next().Value, "ADaxConn#Commit")
	}
	elem = elem.Next().Next()
	if elem.Value == "ADaxConn#Close" {
		assert.Equal(t, elem.Value, "ADaxConn#Close")
		assert.Equal(t, elem.Next().Value, "BDaxConn#Close")
	} else {
		assert.Equal(t, elem.Value, "BDaxConn#Close")
		assert.Equal(t, elem.Next().Value, "ADaxConn#Close")
	}
}

func TestRunTxn_failToCreateDaxConn(t *testing.T) {
	clear()
	defer clear()

	willFailToCreateBDaxConn = true

	type ABDax interface {
		GetAData() (string, sabi.Err)
		SetBData(data string) sabi.Err
	}

	base := sabi.NewDaxBase()
	base = struct {
		sabi.DaxBase
		AGetDax
		BGetSetDax
	}{
		DaxBase:    base,
		AGetDax:    AGetDax{ADax: ADax{Dax: base}},
		BGetSetDax: BGetSetDax{BDax: BDax{Dax: base}},
	}

	aDs := NewADaxSrc()
	bDs := NewBDaxSrc()
	base.AddLocalDaxSrc("aaa", aDs)
	base.AddLocalDaxSrc("bbb", bDs)

	aDs.AMap["a"] = "hello"
	err := sabi.RunTxn(base, func(dax ABDax) sabi.Err {
		data, err := dax.GetAData()
		if !err.IsOk() {
			return err
		}
		data = strings.ToUpper(data)
		err = dax.SetBData(data)
		return err
	})
	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case sabi.FailToCreateDaxConn:
		reason := err.Reason().(sabi.FailToCreateDaxConn)
		assert.Equal(t, reason.Name, "bbb")
		cause := err.Cause().(sabi.Err)
		switch cause.Reason().(type) {
		case FailToCreateBDaxConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, bDs.BMap["b"], "")

	assert.Equal(t, logs.Front().Value, "ADaxConn#Rollback")
	assert.Equal(t, logs.Front().Next().Value, "ADaxConn#Close")
}

func TestRunTxn_failToCommitDaxConn(t *testing.T) {
	clear()
	defer clear()

	willFailToCommitBDaxConn = true

	type ABDax interface {
		GetAData() (string, sabi.Err)
		SetBData(data string) sabi.Err
	}

	base := sabi.NewDaxBase()
	base = struct {
		sabi.DaxBase
		AGetDax
		BGetSetDax
	}{
		DaxBase:    base,
		AGetDax:    AGetDax{ADax: ADax{Dax: base}},
		BGetSetDax: BGetSetDax{BDax: BDax{Dax: base}},
	}

	aDs := NewADaxSrc()
	bDs := NewBDaxSrc()
	base.AddLocalDaxSrc("aaa", aDs)
	base.AddLocalDaxSrc("bbb", bDs)

	aDs.AMap["a"] = "hello"
	err := sabi.RunTxn(base, func(dax ABDax) sabi.Err {
		data, err := dax.GetAData()
		if !err.IsOk() {
			return err
		}
		data = strings.ToUpper(data)
		err = dax.SetBData(data)
		return err
	})
	assert.False(t, err.IsOk())

	switch err.Reason().(type) {
	case sabi.FailToCommitDaxConn:
		reason := err.Reason().(sabi.FailToCommitDaxConn)
		assert.Equal(t, len(reason.Errors), 1)
		err = reason.Errors["bbb"]
		switch err.Reason().(type) {
		case FailToCommitBDaxConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, logs.Front().Value, "ADaxConn#Commit")

	elem := logs.Front().Next()
	if elem.Value == "ADaxConn#Rollback" {
		assert.Equal(t, elem.Value, "ADaxConn#Rollback")
		assert.Equal(t, elem.Next().Value, "BDaxConn#Rollback")
	} else {
		assert.Equal(t, elem.Value, "BDaxConn#Rollback")
		assert.Equal(t, elem.Next().Value, "ADaxConn#Rollback")
	}
	elem = elem.Next().Next()
	if elem.Value == "ADaxConn#Close" {
		assert.Equal(t, elem.Value, "ADaxConn#Close")
		assert.Equal(t, elem.Next().Value, "BDaxConn#Close")
	} else {
		assert.Equal(t, elem.Value, "BDaxConn#Close")
		assert.Equal(t, elem.Next().Value, "ADaxConn#Close")
	}
}

func TestRunTxn_Run_errorInLogic(t *testing.T) {
	clear()
	defer clear()

	willFailToCommitBDaxConn = true

	type ABDax interface {
		GetAData() (string, sabi.Err)
		SetBData(data string) sabi.Err
	}

	base := sabi.NewDaxBase()
	base = struct {
		sabi.DaxBase
		AGetDax
		BGetSetDax
	}{
		DaxBase:    base,
		AGetDax:    AGetDax{ADax: ADax{Dax: base}},
		BGetSetDax: BGetSetDax{BDax: BDax{Dax: base}},
	}

	aDs := NewADaxSrc()
	bDs := NewBDaxSrc()
	base.AddLocalDaxSrc("aaa", aDs)
	base.AddLocalDaxSrc("bbb", bDs)

	aDs.AMap["a"] = "hello"
	err := sabi.RunTxn(base, func(dax ABDax) sabi.Err {
		_, err := dax.GetAData()
		if !err.IsOk() {
			return err
		}
		return sabi.NewErr(FailToRunLogic{})
	})
	assert.False(t, err.IsOk())

	switch err.Reason().(type) {
	case FailToRunLogic:
	default:
		assert.Fail(t, err.Error())
	}

	elem := logs.Front()
	assert.Equal(t, elem.Value, "ADaxConn#Rollback")
	assert.Equal(t, elem.Next().Value, "ADaxConn#Close")
}
