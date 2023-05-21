package sabi_test

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"testing"
	"time"
)

var runnerLogs list.List

func clearRunnerLogs() {
	runnerLogs.Init()
}

func TestRunSeq(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	err := sabi.RunSeq(slowerRunner, fasterRunner)
	assert.True(t, err.IsOk())

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestRunSeq_failToRunFormerRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		return sabi.NewErr(FailToRun{})
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	err := sabi.RunSeq(slowerRunner, fasterRunner)
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

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		return sabi.NewErr(FailToRun{})
	}

	err := sabi.RunSeq(slowerRunner, fasterRunner)
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

func TestSeq(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	runner := sabi.Seq(slowerRunner, fasterRunner)
	err := runner()
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

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		return sabi.NewErr(FailToRun{})
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	runner := sabi.Seq(slowerRunner, fasterRunner)
	err := runner()
	assert.True(t, err.IsNotOk())
	switch err.Reason().(type) {
	case FailToRun:
	default:
		assert.Fail(t, err.Error())
	}

	log := runnerLogs.Front()
	assert.Nil(t, log)
}

func TestSeq_failToRunLatterRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		return sabi.NewErr(FailToRun{})
	}

	runner := sabi.Seq(slowerRunner, fasterRunner)
	err := runner()
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

func TestRunPara(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	err := sabi.RunPara(slowerRunner, fasterRunner)
	assert.True(t, err.IsOk())

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestRunPara_failToRunFormerRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		return sabi.NewErr(FailToRun{})
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	err := sabi.RunPara(slowerRunner, fasterRunner)
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

func TestRunPara_failToRunLatterRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		return sabi.NewErr(FailToRun{})
	}

	err := sabi.RunPara(slowerRunner, fasterRunner)
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

func TestPara(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	runner := sabi.Para(slowerRunner, fasterRunner)
	err := runner()
	assert.True(t, err.IsOk())

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestPara_failToRunSlowerRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		return sabi.NewErr(FailToRun{})
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	runner := sabi.Para(slowerRunner, fasterRunner)
	err := runner()
	assert.True(t, err.IsNotOk())
	switch err.Reason().(type) {
	case sabi.FailToRunInParallel:
		errs := err.Reason().(sabi.FailToRunInParallel).Errors
		assert.Equal(t, len(errs), 1)
		switch errs[0].Reason().(type) {
		case FailToRun:
		default:
			assert.Fail(t, errs[1].Error())
		}
	default:
		assert.Fail(t, err.Error())
	}

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Nil(t, log)
}

func TestPara_failToRunFasterRunner(t *testing.T) {
	clearRunnerLogs()
	defer clearRunnerLogs()

	type FailToRun struct{}

	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		return sabi.NewErr(FailToRun{})
	}

	runner := sabi.Para(slowerRunner, fasterRunner)
	err := runner()
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

func TestRunner_ifOk(t *testing.T) {
	slowerRunner := func() sabi.Err {
		time.Sleep(50 * time.Millisecond)
		runnerLogs.PushBack("slower runner.")
		return sabi.Ok()
	}

	fasterRunner := func() sabi.Err {
		time.Sleep(10 * time.Millisecond)
		runnerLogs.PushBack("faster runner.")
		return sabi.Ok()
	}

	seq0 := sabi.Seq(slowerRunner)
	seq1 := sabi.Seq(fasterRunner)

	err := seq0().IfOk(seq1)
	assert.True(t, err.IsOk())

	log := runnerLogs.Front()
	assert.Equal(t, log.Value, "slower runner.")
	log = log.Next()
	assert.Equal(t, log.Value, "faster runner.")
	log = log.Next()
	assert.Nil(t, log)
}
