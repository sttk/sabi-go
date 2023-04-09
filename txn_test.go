package sabi_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"strings"
	"testing"
)

// ====== Logic part =======

type MyDax interface {
	GetData() (string, sabi.Err)
	SetData(data string) sabi.Err
}

func GetAndSetDataLogic(dax MyDax) sabi.Err {
	data, err := dax.GetData()
	if !err.IsOk() {
		return err
	}
	data = strings.ToUpper(data)
	return dax.SetData(data)
}

// ====== Data access part ========

type FooGetDataDax struct {
	sabi.FooDax
}

func NewFooGetDataDax(base sabi.Dax) FooGetDataDax {
	return FooGetDataDax{FooDax: sabi.NewFooDax(base)}
}

func (dax FooGetDataDax) GetData() (string, sabi.Err) {
	_, err := dax.GetFooDaxConn("foo")
	if !err.IsOk() {
		return "", err
	}
	data := "GetData"
	return data, sabi.Ok()
}

type BarSetDataDax struct {
	sabi.BarDax
}

func NewBarSetDataDax(base sabi.Dax) BarSetDataDax {
	return BarSetDataDax{BarDax: sabi.NewBarDax(base)}
}

func (dax BarSetDataDax) SetData(data string) sabi.Err {
	conn, err := dax.GetBarDaxConn("bar")
	if !err.IsOk() {
		return err
	}
	conn.Store("result", data)
	return sabi.Ok()
}

// ====== Procedure ========

func NewDaxBase() sabi.DaxBase {
	base := sabi.NewDaxBase()
	return struct {
		sabi.DaxBase
		FooGetDataDax
		BarSetDataDax
	}{
		DaxBase:       base,
		FooGetDataDax: NewFooGetDataDax(base),
		BarSetDataDax: NewBarSetDataDax(base),
	}
}

func TestRunTxn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	sabi.AddGlobalDaxSrc("foo", sabi.FooDaxSrc{})
	sabi.FixGlobalDaxSrcs()

	store := make(map[string]string)

	base := NewDaxBase()
	base.AddLocalDaxSrc("bar", &sabi.BarDaxSrc{Store: store})

	err := sabi.RunTxn(base, GetAndSetDataLogic)
	assert.True(t, err.IsOk())

	assert.Equal(t, store["result"], "GETDATA")
}

func TestRunTxn_failToGetDaxConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	sabi.AddGlobalDaxSrc("foo", sabi.FooDaxSrc{})
	sabi.FixGlobalDaxSrcs()

	store := make(map[string]string)

	base := NewDaxBase()
	base.AddLocalDaxSrc("bar", &sabi.BarDaxSrc{Store: store})

	sabi.WillFailToCreateFooDaxConn = true

	err := sabi.RunTxn(base, GetAndSetDataLogic)
	switch err.Reason().(type) {
	case sabi.FailToCreateDaxConn:
		assert.Equal(t, err.Get("Name"), "foo")
		switch err.Cause().(sabi.Err).Reason().(type) {
		case sabi.InvalidDaxConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, store["result"], "")
}

func TestProc_RunTxn_failToCommitDaxConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	sabi.AddGlobalDaxSrc("foo", sabi.FooDaxSrc{})
	sabi.FixGlobalDaxSrcs()

	store := make(map[string]string)

	base := NewDaxBase()
	base.AddLocalDaxSrc("bar", &sabi.BarDaxSrc{Store: store})

	sabi.WillFailToCommitFooDaxConn = true

	err := sabi.RunTxn(base, GetAndSetDataLogic)
	switch err.Reason().(type) {
	case sabi.FailToCommitDaxConn:
		errs := err.Get("Errors").(map[string]sabi.Err)
		assert.Equal(t, len(errs), 1)
		switch errs["foo"].Reason().(type) {
		case sabi.InvalidDaxConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, store["result"], "GETDATA")
}

func TestTxn_Run(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	store := make(map[string]string)

	base := NewDaxBase()
	base.AddLocalDaxSrc("foo", sabi.FooDaxSrc{})
	base.AddLocalDaxSrc("bar", &sabi.BarDaxSrc{Store: store})

	txn := sabi.Txn(base, GetAndSetDataLogic)

	err := txn.Run()
	assert.True(t, err.IsOk())

	assert.Equal(t, store["result"], "GETDATA")
}

func TestTxn_Run_failToGetConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	store := make(map[string]string)

	base := NewDaxBase()
	base.AddLocalDaxSrc("foo", sabi.FooDaxSrc{})
	base.AddLocalDaxSrc("bar", &sabi.BarDaxSrc{Store: store})

	txn := sabi.Txn(base, GetAndSetDataLogic)

	sabi.WillFailToCreateFooDaxConn = true

	err := txn.Run()
	switch err.Reason().(type) {
	case sabi.FailToCreateDaxConn:
		assert.Equal(t, err.Get("Name"), "foo")
		switch err.Cause().(sabi.Err).Reason().(type) {
		case sabi.InvalidDaxConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, store["result"], "")
}

func TestTxn_Run_failToCommitConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	store := make(map[string]string)

	base := NewDaxBase()
	base.AddLocalDaxSrc("foo", sabi.FooDaxSrc{})
	base.AddLocalDaxSrc("bar", &sabi.BarDaxSrc{Store: store})

	txn := sabi.Txn(base, GetAndSetDataLogic)

	sabi.WillFailToCommitFooDaxConn = true

	err := txn.Run()
	switch err.Reason().(type) {
	case sabi.FailToCommitDaxConn:
		errs := err.Get("Errors").(map[string]sabi.Err)
		assert.Equal(t, len(errs), 1)
		switch errs["foo"].Reason().(type) {
		case sabi.InvalidDaxConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, store["result"], "GETDATA")
}
