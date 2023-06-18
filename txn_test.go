package sabi_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sttk/sabi"
)

func TestRunTxn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

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
		AGetDax:    AGetDax{Dax: base},
		BGetSetDax: BGetSetDax{Dax: base},
	}

	aDs := NewADaxSrc()
	bDs := NewBDaxSrc()

	err := base.SetUpLocalDaxSrc("aaa", aDs)
	assert.True(t, err.IsOk())
	err = base.SetUpLocalDaxSrc("bbb", bDs)
	assert.True(t, err.IsOk())

	aDs.AMap["a"] = "hello"
	err = sabi.RunTxn(base, func(dax ABDax) sabi.Err {
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

	log := Logs.Front()
	if log.Value == "ADaxConn#Commit" {
		assert.Equal(t, log.Value, "ADaxConn#Commit")
		log = log.Next()
		assert.Equal(t, log.Value, "BDaxConn#Commit")
	} else {
		assert.Equal(t, log.Value, "BDaxConn#Commit")
		log = log.Next()
		assert.Equal(t, log.Value, "ADaxConn#Commit")
	}
	log = log.Next()
	if log.Value == "ADaxConn#Close" {
		assert.Equal(t, log.Value, "ADaxConn#Close")
		log = log.Next()
		assert.Equal(t, log.Value, "BDaxConn#Close")
	} else {
		assert.Equal(t, log.Value, "BDaxConn#Close")
		log = log.Next()
		assert.Equal(t, log.Value, "ADaxConn#Close")
	}
	log = log.Next()
	assert.Nil(t, log)
}

func TestRunTxn_failToCreateDaxConn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	WillFailToCreateBDaxConn = true

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
		AGetDax:    AGetDax{Dax: base},
		BGetSetDax: BGetSetDax{Dax: base},
	}

	aDs := NewADaxSrc()
	bDs := NewBDaxSrc()

	err := base.SetUpLocalDaxSrc("aaa", aDs)
	assert.True(t, err.IsOk())
	err = base.SetUpLocalDaxSrc("bbb", bDs)
	assert.True(t, err.IsOk())

	aDs.AMap["a"] = "hello"
	err = sabi.RunTxn(base, func(dax ABDax) sabi.Err {
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

	assert.Equal(t, Logs.Front().Value, "ADaxConn#Rollback")
	assert.Equal(t, Logs.Front().Next().Value, "ADaxConn#Close")
}

func TestRunTxn_failToCommitDaxConn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	WillFailToCommitBDaxConn = true

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
		AGetDax:    AGetDax{Dax: base},
		BGetSetDax: BGetSetDax{Dax: base},
	}

	aDs := NewADaxSrc()
	bDs := NewBDaxSrc()

	err := base.SetUpLocalDaxSrc("aaa", aDs)
	assert.True(t, err.IsOk())
	err = base.SetUpLocalDaxSrc("bbb", bDs)
	assert.True(t, err.IsOk())

	aDs.AMap["a"] = "hello"
	err = sabi.RunTxn(base, func(dax ABDax) sabi.Err {
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

	log := Logs.Front()
	assert.Equal(t, log.Value, "ADaxConn#Commit")
	log = log.Next()
	if log.Value == "ADaxConn#Rollback" {
		assert.Equal(t, log.Value, "ADaxConn#Rollback")
		log = log.Next()
		assert.Equal(t, log.Value, "BDaxConn#Rollback")
	} else {
		assert.Equal(t, log.Value, "BDaxConn#Rollback")
		log = log.Next()
		assert.Equal(t, log.Value, "ADaxConn#Rollback")
	}
	log = log.Next()
	if log.Value == "ADaxConn#Close" {
		assert.Equal(t, log.Value, "ADaxConn#Close")
		log = log.Next()
		assert.Equal(t, log.Value, "BDaxConn#Close")
	} else {
		assert.Equal(t, log.Value, "BDaxConn#Close")
		log = log.Next()
		assert.Equal(t, log.Value, "ADaxConn#Close")
	}
	log = log.Next()
	assert.Nil(t, log)
}

func TestRunTxn_Run_errorInLogic(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	WillFailToCommitBDaxConn = true

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
		AGetDax:    AGetDax{Dax: base},
		BGetSetDax: BGetSetDax{Dax: base},
	}

	aDs := NewADaxSrc()
	bDs := NewBDaxSrc()

	err := base.SetUpLocalDaxSrc("aaa", aDs)
	assert.True(t, err.IsOk())
	err = base.SetUpLocalDaxSrc("bbb", bDs)
	assert.True(t, err.IsOk())

	aDs.AMap["a"] = "hello"
	err = sabi.RunTxn(base, func(dax ABDax) sabi.Err {
		_, err = dax.GetAData()
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

	log := Logs.Front()
	assert.Equal(t, log.Value, "ADaxConn#Rollback")
	log = log.Next()
	assert.Equal(t, log.Value, "ADaxConn#Close")
}
