package sabi_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"testing"
)

func TestTxn_Run(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	proc := NewProc()
	proc.AddLocalConnCfg("foo", sabi.FooConnCfg{})
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{})

	txn := proc.NewTxn(GetAndSetDataLogic)

	m, err := txn.Run()
	assert.True(t, err.IsOk())
	assert.Equal(t, m["result"], "GETDATA")
}

func TestTxn_Run_failToGetConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	proc := NewProc()
	proc.AddLocalConnCfg("foo", sabi.FooConnCfg{})
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{})

	txn := proc.NewTxn(GetAndSetDataLogic)

	sabi.WillFailToCreateFooConn = true

	m, err := txn.Run()
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

func TestTxn_Run_failToCommitConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	proc := NewProc()
	proc.AddLocalConnCfg("foo", sabi.FooConnCfg{})
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{})

	txn := proc.NewTxn(GetAndSetDataLogic)

	sabi.WillFailToCommitFooConn = true

	m, err := txn.Run()
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
