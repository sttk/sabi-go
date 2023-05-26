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

type InvalidValueError struct {
	Value string
}

func (e InvalidValueError) Error() string {
	return "InvalidValue{Value=" + e.Value + "}"
}

func TestNewErr_reasonIsValue(t *testing.T) {
	err := sabi.NewErr(InvalidValue{Value: "abc"})

	assert.Equal(t, err.Error(), "{reason=InvalidValue, Value=abc}")

	switch err.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, err.IsOk())
	assert.True(t, err.IsNotOk())
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
	assert.Nil(t, errors.Unwrap(err))

	assert.True(t, errors.Is(err, err))
	assert.True(t, errors.As(err, &err))

	e := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(err, e))
	assert.False(t, errors.As(err, &e))
}

func TestNewErr_reasonIsPointer(t *testing.T) {
	err := sabi.NewErr(&InvalidValue{Value: "abc"})

	assert.Equal(t, err.Error(), "{reason=InvalidValue, Value=abc}")

	switch err.Reason().(type) {
	case *InvalidValue:
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, err.IsOk())
	assert.True(t, err.IsNotOk())
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
	assert.Nil(t, errors.Unwrap(err))

	assert.True(t, errors.Is(err, err))
	assert.True(t, errors.As(err, &err))

	e := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(err, e))
	assert.False(t, errors.As(err, &e))
}

func TestNewErr_withCause(t *testing.T) {
	cause := errors.New("def")
	err := sabi.NewErr(InvalidValue{Value: "abc"}, cause)

	assert.Equal(t, err.Error(), "{reason=InvalidValue, Value=abc, cause=def}")

	switch err.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, err.IsOk())
	assert.True(t, err.IsNotOk())
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

	assert.True(t, errors.Is(err, err))
	assert.True(t, errors.As(err, &err))

	assert.True(t, errors.Is(err, cause))
	//assert.True(t, errors.As(err, cause)) --> compile error

	e := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(err, e))
	assert.False(t, errors.As(err, &e))
}

func TestNewErr_causeIsAlsoErr(t *testing.T) {
	cause := sabi.NewErr(FailToGetValue{Name: "foo"})
	err := sabi.NewErr(InvalidValue{Value: "abc"}, cause)

	assert.Equal(t, err.Error(), "{reason=InvalidValue, Value=abc, cause={reason=FailToGetValue, Name=foo}}")

	switch err.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, err.IsOk())
	assert.True(t, err.IsNotOk())
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

	assert.True(t, errors.Is(err, err))
	assert.True(t, errors.As(err, &err))

	assert.True(t, errors.Is(err, cause))
	assert.True(t, errors.As(err, &cause))

	e := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(err, e))
	assert.False(t, errors.As(err, &e))
}

func TestOk(t *testing.T) {
	err := sabi.Ok()

	assert.Equal(t, err.Error(), "{reason=nil}")
	assert.Nil(t, err.Reason())

	switch err.Reason().(type) {
	case nil:
	default:
		assert.Fail(t, err.Error())
	}

	assert.True(t, err.IsOk())
	assert.False(t, err.IsNotOk())
	assert.Equal(t, err.ReasonName(), "")
	assert.Equal(t, err.ReasonPackage(), "")
	assert.Nil(t, err.Get("Value"))
	assert.Nil(t, err.Get("value"))
	assert.Nil(t, err.Get("Name"))

	m := err.Situation()
	assert.Equal(t, len(m), 0)

	assert.Nil(t, err.Cause())
	assert.Nil(t, err.Unwrap())
	assert.Nil(t, errors.Unwrap(err))

	assert.True(t, errors.Is(err, err))
	assert.True(t, errors.As(err, &err))

	e := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(err, e))
	assert.False(t, errors.As(err, &e))
}

func TestErr_IfOk_ok(t *testing.T) {
	err := sabi.Ok()

	s := ""
	err2 := err.IfOk(func() sabi.Err {
		s = "executed."
		return sabi.NewErr(InvalidValue{Value: "x"})
	})

	assert.Equal(t, s, "executed.")
	assert.True(t, err.IsOk())
	assert.True(t, err2.IsNotOk())
	switch err2.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, err2.Error())
	}
}

func TestErr_IfOk_err(t *testing.T) {
	err := sabi.NewErr(InvalidValue{Value: "x"})

	s := ""
	err2 := err.IfOk(func() sabi.Err {
		s = "executed."
		return sabi.Ok()
	})

	assert.Equal(t, s, "")
	assert.True(t, err.IsNotOk())
	assert.True(t, err2.IsNotOk())
	switch err2.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, err2.Error())
	}
}
