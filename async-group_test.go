package sabi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sttk/sabi/errs"
)

func TestAsyncGroup_asyncGroupSync_ok(t *testing.T) {
	var ag asyncGroupSync
	assert.True(t, ag.err.IsOk())

	exec := false
	fn := func() errs.Err {
		exec = true
		return errs.Ok()
	}

	ag.Add(fn)
	assert.True(t, ag.err.IsOk())
	assert.True(t, exec)
}

func TestAsyncGroup_asyncGroupSync_error(t *testing.T) {
	var ag asyncGroupSync
	assert.True(t, ag.err.IsOk())

	type FailToDoSomething struct{}

	exec := false
	fn := func() errs.Err {
		exec = true
		return errs.New(FailToDoSomething{})
	}

	ag.Add(fn)
	switch ag.err.Reason().(type) {
	case FailToDoSomething:
	default:
		assert.Fail(t, ag.err.Error())
	}
	assert.True(t, exec)
}

func TestAsyncGroup_asyncGroupAsync_ok(t *testing.T) {
	var ag asyncGroupAsync[string]
	assert.False(t, ag.hasErr())

	exec := false
	fn := func() errs.Err {
		time.Sleep(50)
		exec = true
		return errs.Ok()
	}

	ag.name = "foo"
	ag.Add(fn)
	assert.False(t, ag.hasErr())
	assert.False(t, exec)

	ag.wait()
	assert.False(t, ag.hasErr())
	assert.True(t, exec)

	assert.Equal(t, len(ag.makeErrs()), 0)
	assert.True(t, exec)
}

func TestAsyncGroup_asyncGroupAsync_error(t *testing.T) {
	var ag asyncGroupAsync[string]
	assert.False(t, ag.hasErr())

	type FailToDoSomething struct{}

	exec := false
	fn := func() errs.Err {
		time.Sleep(50)
		exec = true
		return errs.New(FailToDoSomething{})
	}

	ag.name = "foo"
	ag.Add(fn)
	assert.False(t, ag.hasErr())
	assert.False(t, exec)

	ag.wait()
	assert.True(t, ag.hasErr())
	assert.True(t, exec)

	m := ag.makeErrs()
	assert.Equal(t, len(m), 1)
	switch m["foo"].Reason().(type) {
	case FailToDoSomething:
	default:
		assert.Fail(t, m["foo"].Error())
	}
	assert.True(t, exec)
}

func TestAsyncGroup_asyncGroupAsync_multipleErrors(t *testing.T) {
	var ag asyncGroupAsync[string]
	assert.False(t, ag.hasErr())

	type Err0 struct{}
	type Err1 struct{}
	type Err2 struct{}

	exec0 := false
	exec1 := false
	exec2 := false

	fn0 := func() errs.Err {
		time.Sleep(200)
		exec0 = true
		return errs.New(Err0{})
	}
	fn1 := func() errs.Err {
		time.Sleep(400)
		exec1 = true
		return errs.New(Err1{})
	}
	fn2 := func() errs.Err {
		time.Sleep(800)
		exec2 = true
		return errs.New(Err2{})
	}

	ag.name = "foo0"
	ag.Add(fn0)
	ag.name = "foo1"
	ag.Add(fn1)
	ag.name = "foo2"
	ag.Add(fn2)
	assert.False(t, ag.hasErr())
	assert.False(t, exec0)
	assert.False(t, exec1)
	assert.False(t, exec2)

	ag.wait()
	assert.True(t, ag.hasErr())
	assert.True(t, exec0)
	assert.True(t, exec1)
	assert.True(t, exec2)

	m := ag.makeErrs()
	assert.Equal(t, len(m), 3)
	assert.Equal(t, m["foo0"].ReasonName(), "Err0")
	assert.Equal(t, m["foo1"].ReasonName(), "Err1")
	assert.Equal(t, m["foo2"].ReasonName(), "Err2")
	assert.True(t, exec0)
	assert.True(t, exec1)
	assert.True(t, exec2)
}
