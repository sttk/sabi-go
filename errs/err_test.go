package errs_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sttk/sabi/errs"
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

func TestErr_New_reasonIsValue(t *testing.T) {
	e := errs.New(InvalidValue{Value: "abc"})

	assert.Equal(t, e.Error(), "{reason=InvalidValue, Value=abc}")

	switch e.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, e.Error())
	}

	assert.False(t, e.IsOk())
	assert.True(t, e.IsNotOk())
	assert.Equal(t, e.ReasonName(), "InvalidValue")
	assert.Equal(t, e.ReasonPackage(), "github.com/sttk/sabi/errs_test")
	assert.Equal(t, e.Get("Value"), "abc")
	assert.Nil(t, e.Get("value"))
	assert.Nil(t, e.Get("Name"))

	m := e.Situation()
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m["Value"], "abc")
	assert.Nil(t, m["value"])

	assert.Nil(t, e.Cause())
	assert.Nil(t, e.Unwrap())
	assert.Nil(t, errors.Unwrap(e))

	assert.True(t, errors.Is(e, e))
	assert.True(t, errors.As(e, &e))

	er := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(e, er))
	assert.False(t, errors.As(e, &er))
}

func TestErr_New_reasonIsPointer(t *testing.T) {
	e := errs.New(&InvalidValue{Value: "abc"})

	assert.Equal(t, e.Error(), "{reason=InvalidValue, Value=abc}")

	switch e.Reason().(type) {
	case *InvalidValue:
	default:
		assert.Fail(t, e.Error())
	}

	assert.False(t, e.IsOk())
	assert.True(t, e.IsNotOk())
	assert.Equal(t, e.ReasonName(), "InvalidValue")
	assert.Equal(t, e.ReasonPackage(), "github.com/sttk/sabi/errs_test")
	assert.Equal(t, e.Get("Value"), "abc")
	assert.Nil(t, e.Get("value"))
	assert.Nil(t, e.Get("Name"))

	m := e.Situation()
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m["Value"], "abc")

	assert.Nil(t, e.Cause())
	assert.Nil(t, e.Unwrap())
	assert.Nil(t, errors.Unwrap(e))

	assert.True(t, errors.Is(e, e))
	assert.True(t, errors.As(e, &e))

	er := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(e, er))
	assert.False(t, errors.As(e, &er))
}

func TestErr_New_withCause(t *testing.T) {
	cause := errors.New("def")
	e := errs.New(InvalidValue{Value: "abc"}, cause)

	assert.Equal(t, e.Error(), "{reason=InvalidValue, Value=abc, cause=def}")

	switch e.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, e.Error())
	}

	assert.False(t, e.IsOk())
	assert.True(t, e.IsNotOk())
	assert.Equal(t, e.ReasonName(), "InvalidValue")
	assert.Equal(t, e.ReasonPackage(), "github.com/sttk/sabi/errs_test")
	assert.Equal(t, e.Get("Value"), "abc")
	assert.Nil(t, e.Get("value"))
	assert.Nil(t, e.Get("Name"))

	m := e.Situation()
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m["Value"], "abc")

	assert.Equal(t, e.Cause(), cause)
	assert.Equal(t, e.Unwrap(), cause)
	assert.Equal(t, errors.Unwrap(e), cause)

	assert.True(t, errors.Is(e, e))
	assert.True(t, errors.As(e, &e))

	assert.True(t, errors.Is(e, cause))
	//assert.True(t, errors.As(e, &cause)) // --> compile error
}

func TestErr_New_causeIsCustomError(t *testing.T) {
	cause := InvalidValueError{Value: "def"}
	e := errs.New(InvalidValue{Value: "abc"}, cause)

	assert.Equal(t, e.Error(), "{reason=InvalidValue, Value=abc, cause=InvalidValue{Value=def}}")

	switch e.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, e.Error())
	}

	assert.False(t, e.IsOk())
	assert.True(t, e.IsNotOk())
	assert.Equal(t, e.ReasonName(), "InvalidValue")
	assert.Equal(t, e.ReasonPackage(), "github.com/sttk/sabi/errs_test")
	assert.Equal(t, e.Get("Value"), "abc")
	assert.Nil(t, e.Get("value"))
	assert.Nil(t, e.Get("Name"))

	m := e.Situation()
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m["Value"], "abc")

	assert.Equal(t, e.Cause(), cause)
	assert.Equal(t, e.Unwrap(), cause)
	assert.Equal(t, errors.Unwrap(e), cause)

	assert.True(t, errors.Is(e, e))
	assert.True(t, errors.As(e, &e))

	assert.True(t, errors.Is(e, cause))
	assert.True(t, errors.As(e, &cause))
}

func TestErr_New_causeIsAlsoErr(t *testing.T) {
	cause := errs.New(FailToGetValue{Name: "foo"})
	e := errs.New(InvalidValue{Value: "abc"}, cause)

	assert.Equal(t, e.Error(), "{reason=InvalidValue, Value=abc, cause={reason=FailToGetValue, Name=foo}}")

	switch e.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, e.Error())
	}

	assert.False(t, e.IsOk())
	assert.True(t, e.IsNotOk())
	assert.Equal(t, e.ReasonName(), "InvalidValue")
	assert.Equal(t, e.ReasonPackage(), "github.com/sttk/sabi/errs_test")
	assert.Equal(t, e.Get("Value"), "abc")
	assert.Equal(t, e.Get("Name"), "foo")
	assert.Nil(t, e.Get("value"))
	assert.Nil(t, e.Get("name"))

	m := e.Situation()
	assert.Equal(t, len(m), 2)
	assert.Equal(t, m["Value"], "abc")
	assert.Equal(t, m["Name"], "foo")

	assert.Equal(t, e.Cause(), cause)
	assert.Equal(t, e.Unwrap(), cause)
	assert.Equal(t, errors.Unwrap(e), cause)

	assert.True(t, errors.Is(e, e))
	assert.True(t, errors.As(e, &e))

	assert.True(t, errors.Is(e, cause))
	assert.True(t, errors.As(e, &cause))

	er := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(e, er))
	assert.False(t, errors.As(e, &er))
}

func TestErr_Ok(t *testing.T) {
	e := errs.Ok()

	assert.Equal(t, e.Error(), "{reason=nil}")
	assert.Nil(t, e.Reason())

	switch e.Reason().(type) {
	case nil:
	default:
		assert.Fail(t, e.Error())
	}

	assert.True(t, e.IsOk())
	assert.False(t, e.IsNotOk())
	assert.Equal(t, e.ReasonName(), "")
	assert.Equal(t, e.ReasonPackage(), "")
	assert.Nil(t, e.Get("Value"))
	assert.Nil(t, e.Get("Name"))

	m := e.Situation()
	assert.Equal(t, len(m), 0)

	assert.Nil(t, e.Cause())
	assert.Nil(t, e.Unwrap())
	assert.Nil(t, errors.Unwrap(e))

	assert.True(t, errors.Is(e, e))
	assert.True(t, errors.As(e, &e))

	er := InvalidValueError{Value: "aaa"}
	assert.False(t, errors.Is(e, er))
	assert.False(t, errors.As(e, &er))
}

func TestErr_IfOk_ok(t *testing.T) {
	e := errs.Ok()

	s := ""
	e1 := e.IfOk(func() errs.Err {
		s = "executed."
		return errs.New(InvalidValue{Value: "x"})
	})

	assert.Equal(t, s, "executed.")
	assert.True(t, e.IsOk())
	assert.True(t, e1.IsNotOk())

	switch e1.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, e1.Error())
	}
}

func TestErr_IfOk_error(t *testing.T) {
	e := errs.New(InvalidValue{Value: "x"})

	s := ""
	e1 := e.IfOk(func() errs.Err {
		s = "executed."
		return errs.Ok()
	})

	assert.Equal(t, s, "")
	assert.True(t, e.IsNotOk())
	assert.True(t, e1.IsNotOk())

	switch e1.Reason().(type) {
	case InvalidValue:
	default:
		assert.Fail(t, e1.Error())
	}
}
