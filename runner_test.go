package sabi_test

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"testing"
	"time"
)

type (
	FailToRun struct {
		Name string
	}
)

var logs list.List
var errorRunnerName string

func ClearLogs() {
	logs.Init()
	errorRunnerName = ""
}

type MyRunner struct {
	Name string
	Wait time.Duration
}

func (r MyRunner) Run() sabi.Err {
	time.Sleep(r.Wait)
	if r.Name == errorRunnerName {
		return sabi.ErrBy(FailToRun{Name: r.Name})
	}
	logs.PushBack(r.Name)
	return sabi.Ok()
}

func TestSeq(t *testing.T) {
	ClearLogs()
	defer ClearLogs()

	r0 := MyRunner{Name: "r-0", Wait: 50 * time.Millisecond}
	r1 := MyRunner{Name: "r-1", Wait: 10 * time.Millisecond}
	r2 := MyRunner{Name: "r-2", Wait: 20 * time.Millisecond}
	r3 := sabi.Seq(r0, r1, r2)

	err := r3.Run()
	assert.True(t, err.IsOk())

	assert.Equal(t, logs.Len(), 3)
	assert.Equal(t, logs.Front().Value, "r-0")
	assert.Equal(t, logs.Front().Next().Value, "r-1")
	assert.Equal(t, logs.Front().Next().Next().Value, "r-2")
}

func TestSeq_FailToRun(t *testing.T) {
	ClearLogs()
	defer ClearLogs()

	errorRunnerName = "r-1"

	r0 := MyRunner{Name: "r-0", Wait: 50 * time.Millisecond}
	r1 := MyRunner{Name: "r-1", Wait: 10 * time.Millisecond}
	r2 := MyRunner{Name: "r-2", Wait: 20 * time.Millisecond}
	r3 := sabi.Seq(r0, r1, r2)

	err := r3.Run()
	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case FailToRun:
		assert.Equal(t, err.Get("Name"), "r-1")
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, logs.Len(), 1)
	assert.Equal(t, logs.Front().Value, "r-0")
}

func TestPara(t *testing.T) {
	ClearLogs()
	defer ClearLogs()

	r0 := MyRunner{Name: "r-0", Wait: 50 * time.Millisecond}
	r1 := MyRunner{Name: "r-1", Wait: 10 * time.Millisecond}
	r2 := MyRunner{Name: "r-2", Wait: 20 * time.Millisecond}
	r3 := sabi.Para(r0, r1, r2)

	err := r3.Run()
	assert.True(t, err.IsOk())

	assert.Equal(t, logs.Len(), 3)
	assert.Equal(t, logs.Front().Value, "r-1")
	assert.Equal(t, logs.Front().Next().Value, "r-2")
	assert.Equal(t, logs.Front().Next().Next().Value, "r-0")
}

func TestPara_FailToRun(t *testing.T) {
	ClearLogs()
	defer ClearLogs()

	errorRunnerName = "r-1"

	r0 := MyRunner{Name: "r-0", Wait: 50 * time.Millisecond}
	r1 := MyRunner{Name: "r-1", Wait: 10 * time.Millisecond}
	r2 := MyRunner{Name: "r-2", Wait: 20 * time.Millisecond}
	r3 := sabi.Para(r0, r1, r2)

	err := r3.Run()
	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case sabi.FailToRunInParallel:
		errs := err.Get("Errors").(map[int]sabi.Err)
		assert.Equal(t, len(errs), 1)
		assert.Equal(t, errs[0].Error(), "{reason=FailToRun, Name=r-1}")
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, logs.Len(), 2)
	assert.Equal(t, logs.Front().Value, "r-2")
	assert.Equal(t, logs.Front().Next().Value, "r-0")
}

func TestRunSeq(t *testing.T) {
	ClearLogs()
	defer ClearLogs()

	r0 := MyRunner{Name: "r-0", Wait: 50 * time.Millisecond}
	r1 := MyRunner{Name: "r-1", Wait: 10 * time.Millisecond}
	r2 := MyRunner{Name: "r-2", Wait: 20 * time.Millisecond}

	err := sabi.RunSeq(r0, r1, r2)
	assert.True(t, err.IsOk())

	assert.Equal(t, logs.Len(), 3)
	assert.Equal(t, logs.Front().Value, "r-0")
	assert.Equal(t, logs.Front().Next().Value, "r-1")
	assert.Equal(t, logs.Front().Next().Next().Value, "r-2")
}

func TestRunSeq_FailToRun(t *testing.T) {
	ClearLogs()
	defer ClearLogs()

	errorRunnerName = "r-1"

	r0 := MyRunner{Name: "r-0", Wait: 50 * time.Millisecond}
	r1 := MyRunner{Name: "r-1", Wait: 10 * time.Millisecond}
	r2 := MyRunner{Name: "r-2", Wait: 20 * time.Millisecond}

	err := sabi.RunSeq(r0, r1, r2)

	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case FailToRun:
		assert.Equal(t, err.Get("Name"), "r-1")
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, logs.Len(), 1)
	assert.Equal(t, logs.Front().Value, "r-0")
}

func TestRunPara(t *testing.T) {
	ClearLogs()
	defer ClearLogs()

	r0 := MyRunner{Name: "r-0", Wait: 50 * time.Millisecond}
	r1 := MyRunner{Name: "r-1", Wait: 10 * time.Millisecond}
	r2 := MyRunner{Name: "r-2", Wait: 20 * time.Millisecond}

	err := sabi.RunPara(r0, r1, r2)
	assert.True(t, err.IsOk())

	assert.Equal(t, logs.Len(), 3)
	assert.Equal(t, logs.Front().Value, "r-1")
	assert.Equal(t, logs.Front().Next().Value, "r-2")
	assert.Equal(t, logs.Front().Next().Next().Value, "r-0")
}

func TestRunPara_FailToRun(t *testing.T) {
	ClearLogs()
	defer ClearLogs()

	errorRunnerName = "r-1"

	r0 := MyRunner{Name: "r-0", Wait: 50 * time.Millisecond}
	r1 := MyRunner{Name: "r-1", Wait: 10 * time.Millisecond}
	r2 := MyRunner{Name: "r-2", Wait: 20 * time.Millisecond}

	err := sabi.RunPara(r0, r1, r2)

	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case sabi.FailToRunInParallel:
		errs := err.Get("Errors").(map[int]sabi.Err)
		assert.Equal(t, len(errs), 1)
		assert.Equal(t, errs[0].Error(), "{reason=FailToRun, Name=r-1}")
	default:
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, logs.Len(), 2)
	assert.Equal(t, logs.Front().Value, "r-2")
	assert.Equal(t, logs.Front().Next().Value, "r-0")
}
