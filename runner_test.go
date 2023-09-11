package sabi_test

import (
	"container/list"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sttk/sabi"
	"github.com/sttk/sabi/errs"
)

var runnerLogs list.List

func clearRunnerLogs() {
	runnerLogs.Init()
}

func TestSeq(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	slowerRunner := func() errs.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return errs.Ok()
	}

	fasterRunner := func() errs.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return errs.Ok()
	}

	err := sabi.Seq(slowerRunner, fasterRunner)
	assert.True(t, err.IsOk())

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestSeq_failToRunFormerRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() errs.Err {
		time.Sleep(50 * time.Millisecond)
		return errs.New(FailToRun{})
	}

	fasterRunner := func() errs.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return errs.Ok()
	}

	err := sabi.Seq(slowerRunner, fasterRunner)
	assert.True(t, err.IsNotOk())
	switch err.Reason().(type) {
	case FailToRun:
	default:
		assert.Fail(t, err.Error())
	}

	log := runnerLogs.Front()
	assert.Nil(t, log)
}

func TestRunSeq_failToRunLatterRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() errs.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return errs.Ok()
	}

	fasterRunner := func() errs.Err {
		time.Sleep(10 * time.Millisecond)
		return errs.New(FailToRun{})
	}

	err := sabi.Seq(slowerRunner, fasterRunner)
	assert.True(t, err.IsNotOk())
	switch err.Reason().(type) {
	case FailToRun:
	default:
		assert.Fail(t, err.Error())
	}

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestSeq_runner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	slowerRunner := func() errs.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return errs.Ok()
	}

	fasterRunner := func() errs.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return errs.Ok()
	}

	runner := sabi.Seq_(slowerRunner, fasterRunner)
	err := runner()
	assert.True(t, err.IsOk())

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestPara(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	slowerRunner := func() errs.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return errs.Ok()
	}

	fasterRunner := func() errs.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return errs.Ok()
	}

	err := sabi.Para(slowerRunner, fasterRunner)
	assert.True(t, err.IsOk())

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestPara_failToRunFormerRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() errs.Err {
		time.Sleep(50 * time.Millisecond)
		return errs.New(FailToRun{})
	}

	fasterRunner := func() errs.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return errs.Ok()
	}

	err := sabi.Para(slowerRunner, fasterRunner)
	assert.True(t, err.IsNotOk())
	switch err.Reason().(type) {
	case sabi.FailToRunInParallel:
		errs := err.Reason().(sabi.FailToRunInParallel).Errors
		assert.Equal(t, len(errs), 1)
		switch errs[0].Reason().(type) {
		case FailToRun:
		default:
			assert.Fail(t, errs[0].Error())
		}
	default:
		assert.Fail(t, err.Error())
	}

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestPara_failToRunLatterRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() errs.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return errs.Ok()
	}

	fasterRunner := func() errs.Err {
		time.Sleep(10 * time.Millisecond)
		return errs.New(FailToRun{})
	}

	err := sabi.Para(slowerRunner, fasterRunner)
	assert.True(t, err.IsNotOk())
	switch err.Reason().(type) {
	case sabi.FailToRunInParallel:
		errs := err.Reason().(sabi.FailToRunInParallel).Errors
		assert.Equal(t, len(errs), 1)
		switch errs[1].Reason().(type) {
		case FailToRun:
		default:
			assert.Fail(t, errs[1].Error())
		}
	default:
		assert.Fail(t, err.Error())
	}

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestPara_runner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	slowerRunner := func() errs.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return errs.Ok()
	}

	fasterRunner := func() errs.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return errs.Ok()
	}

	runner := sabi.Para_(slowerRunner, fasterRunner)
	err := runner()
	assert.True(t, err.IsOk())

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Nil(t, log)
}
