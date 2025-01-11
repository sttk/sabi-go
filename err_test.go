package sabi_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sttk/sabi"
)

type /* error reasons */ (
	InvalidValue struct {
		Name  string
		Value string
	}

	FailToGetValue struct {
		Name string
	}
)

type InvalidValueError struct {
	Name  string
	Value string
}

func (e InvalidValueError) Error() string {
	return "InvalidValue { Name: " + e.Name + ", Value: " + e.Value + " }"
}

///

func TestErr(t *testing.T) {

	t.Run("NewErr", func(t *testing.T) {
		t.Run("reason is a value", func(t *testing.T) {
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.InvalidValue { Name: foo, Value: abc } }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is a pointer", func(t *testing.T) {
			err := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.InvalidValue { Name: foo, Value: abc } }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is nil", func(t *testing.T) {
			err := sabi.NewErr(nil)

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = nil }")
			assert.Nil(t, err.Cause())
		})

		t.Run("cause is an `error`", func(t *testing.T) {
			cause := errors.New("def")
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.InvalidValue { Name: foo, Value: abc }, cause = def }")
			assert.Equal(t, err.Cause(), cause)
		})

		t.Run("cause is a custom error ", func(t *testing.T) {
			cause := InvalidValueError{Name: "bar", Value: "def"}
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.InvalidValue { Name: foo, Value: abc }, cause = InvalidValue { Name: bar, Value: def } }")
			assert.Equal(t, err.Cause(), cause)
		})

		t.Run("cause is also a sabi.Err", func(t *testing.T) {
			cause := sabi.NewErr(FailToGetValue{Name: "foo"})
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.InvalidValue { Name: foo, Value: abc }, cause = github.com/sttk/sabi.Err { reason = github.com/sttk/sabi_test.FailToGetValue { Name: foo } } }")
			assert.Equal(t, err.Cause(), cause)
		})

		t.Run("reason is nil but cause is specified", func(t *testing.T) {
			cause := errors.New("def")
			err := sabi.NewErr(nil, cause)

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = nil, cause = def }")
			assert.Equal(t, err.Cause(), cause)
		})

		t.Run("reason is pointer to nil", func(t *testing.T) {
			var reason error = nil
			err := sabi.NewErr(&reason)

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = <nil> }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is a boolean", func(t *testing.T) {
			err := sabi.NewErr(true)

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = true }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is a number", func(t *testing.T) {
			err := sabi.NewErr(123)

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = 123 }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is a string", func(t *testing.T) {
			err := sabi.NewErr("abc")

			assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = abc }")
			assert.Nil(t, err.Cause())
		})
	})

	t.Run("Ok", func(t *testing.T) {
		err := sabi.Ok()

		assert.Equal(t, err.Error(), "github.com/sttk/sabi.Err { reason = nil }")
		assert.Nil(t, err.Cause())
	})

	t.Run("Switch expression for reason", func(t *testing.T) {
		t.Run("reason is a value", func(t *testing.T) {
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})

			switch err.Reason().(type) {
			case InvalidValue:
			default:
				assert.Fail(t, err.Error())
			}
		})

		t.Run("reason is a pointer", func(t *testing.T) {
			err := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})

			switch err.Reason().(type) {
			case *InvalidValue:
			default:
				assert.Fail(t, err.Error())
			}
		})
	})

	t.Run("IsOk, IsNotOk", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			err := sabi.Ok()

			assert.True(t, err.IsOk())
			assert.False(t, err.IsNotOk())
		})

		t.Run("reason is a value", func(t *testing.T) {
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})

			assert.False(t, err.IsOk())
			assert.True(t, err.IsNotOk())
		})

		t.Run("reason is a pointer", func(t *testing.T) {
			err := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})

			assert.False(t, err.IsOk())
			assert.True(t, err.IsNotOk())
		})
	})

	t.Run("apply errors.Is", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			err := sabi.Ok()
			assert.Nil(t, err.Unwrap())

			er0 := sabi.Ok()
			er1 := errors.New("def")
			er2 := InvalidValueError{Value: "def"}
			er3 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			er4 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})

			assert.True(t, errors.Is(err, err))
			assert.True(t, errors.Is(err, er0))
			assert.False(t, errors.Is(err, er1))
			assert.False(t, errors.Is(err, er2))
			assert.False(t, errors.Is(err, er3))
			assert.False(t, errors.Is(err, er4))
		})

		t.Run("reason is a value and with no cause", func(t *testing.T) {
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			assert.Nil(t, err.Unwrap())

			er0 := sabi.Ok()
			er1 := errors.New("def")
			er2 := InvalidValueError{Value: "def"}
			er3 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			er4 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, er0))
			assert.False(t, errors.Is(err, er1))
			assert.False(t, errors.Is(err, er2))
			assert.True(t, errors.Is(err, er3))
			assert.False(t, errors.Is(err, er4))
		})

		t.Run("reason is a pointer and with no cause", func(t *testing.T) {
			err := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})
			assert.Nil(t, err.Unwrap())

			er0 := sabi.Ok()
			er1 := errors.New("def")
			er2 := InvalidValueError{Value: "def"}
			er3 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			er4 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, er0))
			assert.False(t, errors.Is(err, er1))
			assert.False(t, errors.Is(err, er2))
			assert.False(t, errors.Is(err, er3))
			assert.False(t, errors.Is(err, er4))
		})

		t.Run("reason is another type and with no cause", func(t *testing.T) {
			err := sabi.NewErr(FailToGetValue{Name: "foo"})
			assert.Nil(t, err.Unwrap())

			er0 := sabi.Ok()
			er1 := errors.New("def")
			er2 := InvalidValueError{Value: "def"}
			er3 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			er4 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, er0))
			assert.False(t, errors.Is(err, er1))
			assert.False(t, errors.Is(err, er2))
			assert.False(t, errors.Is(err, er3))
			assert.False(t, errors.Is(err, er4))
		})

		t.Run("reason is a value and with cause", func(t *testing.T) {
			cause := errors.New("def")
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, cause)
			assert.Equal(t, err.Unwrap(), cause)

			er0 := sabi.Ok()
			er1 := errors.New("def")
			er2 := InvalidValueError{Value: "def"}
			er3 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			er4 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})
			er5 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, er1)
			er6 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"}, er1)
			er7 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, cause)
			er8 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, er0))
			assert.False(t, errors.Is(err, er1))
			assert.False(t, errors.Is(err, er2))
			assert.False(t, errors.Is(err, er3))
			assert.False(t, errors.Is(err, er4))
			assert.False(t, errors.Is(err, er5))
			assert.False(t, errors.Is(err, er6))
			assert.True(t, errors.Is(err, er7))
			assert.False(t, errors.Is(err, er8))
		})

		t.Run("reason is a pointer and with cause", func(t *testing.T) {
			cause := errors.New("def")
			err := sabi.NewErr(&InvalidValue{Value: "abc"}, cause)
			assert.Equal(t, err.Unwrap(), cause)

			er0 := sabi.Ok()
			er1 := errors.New("def")
			er2 := InvalidValueError{Value: "def"}
			er3 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			er4 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})
			er5 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, er1)
			er6 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"}, er5)
			er7 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, cause)
			er8 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, er0))
			assert.False(t, errors.Is(err, er1))
			assert.False(t, errors.Is(err, er2))
			assert.False(t, errors.Is(err, er3))
			assert.False(t, errors.Is(err, er4))
			assert.False(t, errors.Is(err, er5))
			assert.False(t, errors.Is(err, er6))
			assert.False(t, errors.Is(err, er7))
			assert.False(t, errors.Is(err, er8))
		})

		t.Run("reason is another type and with cause", func(t *testing.T) {
			err := sabi.NewErr(FailToGetValue{Name: "foo"})
			assert.Nil(t, err.Unwrap())

			er0 := sabi.Ok()
			er1 := errors.New("def")
			er2 := InvalidValueError{Name: "foo", Value: "def"}
			er3 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			er4 := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"})

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, er0))
			assert.False(t, errors.Is(err, er1))
			assert.False(t, errors.Is(err, er2))
			assert.False(t, errors.Is(err, er3))
			assert.False(t, errors.Is(err, er4))
		})
	})

	t.Run("apply errors.As", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			err := sabi.Ok()
			assert.Nil(t, err.Unwrap())

			var er0 sabi.Err
			//var er1 error
			var er2 InvalidValueError

			assert.True(t, errors.As(err, &er0))
			assert.Equal(t, er0.Error(), err.Error())
			//assert.False(t, errors.As(err, er1)) // -> compile error
			assert.False(t, errors.As(err, &er2))
		})

		t.Run("reason is a value and with no cause", func(t *testing.T) {
			err := sabi.NewErr(FailToGetValue{Name: "foo"})
			assert.Nil(t, err.Unwrap())

			var er0 sabi.Err
			//var er1 error
			var er2 InvalidValueError

			assert.True(t, errors.As(err, &er0))
			assert.Equal(t, er0.Error(), err.Error())
			//assert.False(t, errors.As(err, er1)) // -> compile error
			assert.False(t, errors.As(err, &er2))
		})

		t.Run("reason is a pointer and with no cause", func(t *testing.T) {
			err := sabi.NewErr(&FailToGetValue{Name: "foo"})
			assert.Nil(t, err.Unwrap())

			var er0 sabi.Err
			//var er1 error
			var er2 InvalidValueError

			assert.True(t, errors.As(err, &er0))
			assert.Equal(t, er0.Error(), err.Error())
			//assert.False(t, errors.As(err, er1)) // -> compile error
			assert.False(t, errors.As(err, &er2))
		})

		t.Run("reason is a value and with cause", func(t *testing.T) {
			cause := errors.New("def")
			err := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"}, cause)
			assert.Equal(t, err.Unwrap(), cause)

			var er0 sabi.Err
			//var er1 error
			var er2 InvalidValueError

			assert.True(t, errors.As(err, &er0))
			assert.Equal(t, er0.Error(), err.Error())
			//assert.False(t, errors.As(err, er1)) // -> compile error
			assert.False(t, errors.As(err, &er2))
		})

		t.Run("reason is a pointer and with cause", func(t *testing.T) {
			cause := errors.New("def")
			err := sabi.NewErr(&InvalidValue{Name: "foo", Value: "abc"}, cause)
			assert.Equal(t, err.Unwrap(), cause)

			var er0 sabi.Err
			//var er1 error
			var er2 InvalidValueError

			assert.True(t, errors.As(err, &er0))
			assert.Equal(t, er0.Error(), err.Error())
			//assert.False(t, errors.As(err, er1)) // -> compile error
			assert.False(t, errors.As(err, &er2))
		})
	})

	t.Run("IfOkThen", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			err := sabi.Ok()

			var done bool
			er2 := err.IfOkThen(func() sabi.Err {
				done = true
				return sabi.Ok()
			})
			assert.True(t, er2.IsOk())
			assert.True(t, done)

			er3 := sabi.NewErr(InvalidValue{Name: "foo", Value: "abc"})
			er4 := err.IfOkThen(func() sabi.Err {
				return er3
			})
			assert.Equal(t, er4, er3)
		})

		t.Run("error", func(t *testing.T) {
			err := sabi.NewErr(FailToGetValue{Name: "foo"})

			var done bool
			er2 := err.IfOkThen(func() sabi.Err {
				done = true
				return sabi.Ok()
			})
			assert.Equal(t, er2, err)
			assert.False(t, done)
		})
	})
}
