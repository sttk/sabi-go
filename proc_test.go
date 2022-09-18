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
	_, err := dax.GetFooConn("foo")
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
	_, err := dax.GetBarConn("bar")
	if !err.IsOk() {
		return err
	}
	dax.InnerMap()["result"] = data
	return sabi.Ok()
}

// ====== Procedure ========

func NewProc() sabi.Proc[MyDax] {
	base := sabi.NewConnBase()
	dax := struct {
		FooGetDataDax
		BarSetDataDax
	}{
		FooGetDataDax: NewFooGetDataDax(base),
		BarSetDataDax: NewBarSetDataDax(base),
	}
	return sabi.NewProc[MyDax](base, dax)
}

func TestProc_RunTxn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	sabi.AddGlobalConnCfg("foo", sabi.FooConnCfg{})
	sabi.SealGlobalConnCfgs()

	proc := NewProc()
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{})

	m, err := proc.RunTxn(GetAndSetDataLogic)
	assert.True(t, err.IsOk())
	assert.Equal(t, m["result"], "GETDATA")
}

func TestProc_RunTxn_failToGetConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	sabi.AddGlobalConnCfg("foo", sabi.FooConnCfg{})
	sabi.SealGlobalConnCfgs()

	proc := NewProc()
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{})

	sabi.WillFailToCreateFooConn = true

	m, err := proc.RunTxn(GetAndSetDataLogic)
	switch err.Reason().(type) {
	case sabi.FailToCreateConn:
		assert.Equal(t, err.Get("Name"), "foo")
		switch err.Cause().(sabi.Err).Reason().(type) {
		case sabi.InvalidConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}
	assert.Nil(t, m["result"])
}

func TestProc_RunTxn_failToCommitConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	sabi.AddGlobalConnCfg("foo", sabi.FooConnCfg{})
	sabi.SealGlobalConnCfgs()

	proc := NewProc()
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{})

	sabi.WillFailToCommitFooConn = true

	m, err := proc.RunTxn(GetAndSetDataLogic)
	switch err.Reason().(type) {
	case sabi.FailToCommitConn:
		errs := err.Get("Errors").(map[string]sabi.Err)
		assert.Equal(t, len(errs), 1)
		switch errs["foo"].Reason().(type) {
		case sabi.InvalidConn:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, m["result"], "GETDATA")
}
