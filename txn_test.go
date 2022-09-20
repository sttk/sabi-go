package sabi_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"testing"
)

func TestTxn_Run(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	store := make(map[string]string)

	proc := NewProc()
	proc.AddLocalConnCfg("foo", sabi.FooConnCfg{})
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{Store: store})

	txn := proc.Txn(GetAndSetDataLogic)

	err := txn.Run()
	assert.True(t, err.IsOk())

	assert.Equal(t, store["result"], "GETDATA")
}

func TestTxn_Run_failToGetConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	store := make(map[string]string)

	proc := NewProc()
	proc.AddLocalConnCfg("foo", sabi.FooConnCfg{})
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{Store: store})

	txn := proc.Txn(GetAndSetDataLogic)

	sabi.WillFailToCreateFooConn = true

	err := txn.Run()
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

	assert.Equal(t, store["result"], "")
}

func TestTxn_Run_failToCommitConn(t *testing.T) {
	sabi.Clear()
	defer sabi.Clear()

	store := make(map[string]string)

	proc := NewProc()
	proc.AddLocalConnCfg("foo", sabi.FooConnCfg{})
	proc.AddLocalConnCfg("bar", &sabi.BarConnCfg{Store: store})

	txn := proc.Txn(GetAndSetDataLogic)

	sabi.WillFailToCommitFooConn = true

	err := txn.Run()
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

	assert.Equal(t, store["result"], "GETDATA")
}
