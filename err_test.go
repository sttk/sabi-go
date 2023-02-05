package sabi_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"testing"
)

type /* error reasons */ (
	InvalidValue struct {
		Value string
	}
	FailToGetValue struct {
		Name string
	}
)

func TestNewErr_reasonIsValue(t *testing.T) {
	err := sabi.NewErr(InvalidValue{Value: "abc"})

	assert.Equal(t, err.Error(), "{reason=InvalidValue, Value=abc}")
	assert.Equal(t, err.FileName(), "err_test.go")
	assert.Equal(t, err.LineNumber(), 20)

	switch err.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, err.IsOk())
	assert.Equal(t, err.ReasonName(), "InvalidValue")
	assert.Equal(t, err.ReasonPackage(), "github.com/sttk-go/sabi_test")
	assert.Equal(t, err.Get("Value"), "abc")
	assert.Nil(t, err.Get("value"))
	assert.Nil(t, err.Get("Name"))

	m := err.Situation()
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m["Value"], "abc")
	assert.Nil(t, m["value"])

	assert.Nil(t, err.Cause())
	assert.Nil(t, err.Unwrap())
}

func TestNewErr_reasonIsPointer(t *testing.T) {
	err := sabi.NewErr(&InvalidValue{Value: "abc"})

	assert.Equal(t, err.Error(), "{reason=InvalidValue, Value=abc}")
	assert.Equal(t, err.FileName(), "err_test.go")
	assert.Equal(t, err.LineNumber(), 49)

	switch err.Reason().(type) {
	case *InvalidValue:
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, err.IsOk())
	assert.Equal(t, err.ReasonName(), "InvalidValue")
	assert.Equal(t, err.ReasonPackage(), "github.com/sttk-go/sabi_test")
	assert.Equal(t, err.Get("Value"), "abc")
	assert.Nil(t, err.Get("value"))
	assert.Nil(t, err.Get("Name"))

	m := err.Situation()
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m["Value"], "abc")
	assert.Nil(t, m["value"])

	assert.Nil(t, err.Cause())
	assert.Nil(t, err.Unwrap())
}

func TestNewErr_withCause(t *testing.T) {
	cause := errors.New("def")
	err := sabi.NewErr(InvalidValue{Value: "abc"}, cause)

	assert.Equal(t, err.Error(), "{reason=InvalidValue, Value=abc, cause=def}")
	assert.Equal(t, err.FileName(), "err_test.go")
	assert.Equal(t, err.LineNumber(), 79)

	switch err.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, err.IsOk())
	assert.Equal(t, err.ReasonName(), "InvalidValue")
	assert.Equal(t, err.ReasonPackage(), "github.com/sttk-go/sabi_test")
	assert.Equal(t, err.Get("Value"), "abc")
	assert.Nil(t, err.Get("value"))
	assert.Nil(t, err.Get("Name"))

	m := err.Situation()
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m["Value"], "abc")

	assert.Equal(t, err.Cause(), cause)
	assert.Equal(t, err.Unwrap(), cause)
	assert.Equal(t, errors.Unwrap(err), cause)
}

func TestNewErr_causeIsAlsoErr(t *testing.T) {
	cause := sabi.NewErr(FailToGetValue{Name: "foo"})
	err := sabi.NewErr(InvalidValue{Value: "abc"}, cause)

	assert.Equal(t, err.Error(), "{reason=InvalidValue, Value=abc, cause={reason=FailToGetValue, Name=foo}}")
	assert.Equal(t, err.FileName(), "err_test.go")
	assert.Equal(t, err.LineNumber(), 109)

	switch err.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, err.IsOk())
	assert.Equal(t, err.ReasonName(), "InvalidValue")
	assert.Equal(t, err.ReasonPackage(), "github.com/sttk-go/sabi_test")

	assert.Equal(t, err.Get("Value"), "abc")
	assert.Equal(t, err.Get("Name"), "foo")
	assert.Nil(t, err.Get("value"))

	m := err.Situation()
	assert.Equal(t, len(m), 2)
	assert.Equal(t, m["Value"], "abc")
	assert.Equal(t, m["Name"], "foo")

	assert.Equal(t, err.Cause(), cause)
	assert.Equal(t, err.Unwrap(), cause)
	assert.Equal(t, errors.Unwrap(err), cause)
}

func TestOk(t *testing.T) {
	err := sabi.Ok()

	assert.Equal(t, err.Error(), "{reason=NoError}")
	assert.Equal(t, err.FileName(), "")
	assert.Equal(t, err.LineNumber(), 0)

	switch err.Reason().(type) {
	case sabi.NoError:
	default:
		assert.Fail(t, err.Error())
	}

	assert.True(t, err.IsOk())
	assert.Equal(t, err.ReasonName(), "NoError")
	assert.Equal(t, err.ReasonPackage(), "github.com/sttk-go/sabi")
	assert.Nil(t, err.Get("Value"))
	assert.Nil(t, err.Get("value"))
	assert.Nil(t, err.Get("Name"))

	m := err.Situation()
	assert.Equal(t, len(m), 0)

	assert.Nil(t, err.Cause())
	assert.Nil(t, err.Unwrap())
}
