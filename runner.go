// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"strconv"
)

type /* error reasons */ (
	// FailToRunInParallel is a error reason which indicates some runner which
	// is runned in parallel failed.
	FailToRunInParallel struct {
		Errors map[string]Err
	}
)

// Runner is an interface which has #Run method and is runned by RunSeq or
// RunPara functions.
type Runner interface {
	Run() Err
}

type seqRunner struct {
	runners []Runner
}

func (r seqRunner) Run() Err {
	for _, runner := range r.runners {
		err := runner.Run()
		if !err.IsOk() {
			return err
		}
	}
	return Ok()
}

// Seq is a function which creates a runner which runs multiple runners
// specified as arguments sequencially.
func Seq(runners ...Runner) Runner {
	return seqRunner{runners: runners[:]}
}

type paraRunner struct {
	runners []Runner
}

func (r paraRunner) Run() Err {
	ch := make(chan Err)

	for _, runner := range r.runners {
		go func(runner Runner, ch chan Err) {
			err := runner.Run()
			ch <- err
		}(runner, ch)
	}

	errs := make(map[string]Err)
	n := len(r.runners)
	for i := 0; i < n; i++ {
		select {
		case err := <-ch:
			if !err.IsOk() {
				errs[strconv.Itoa(i)] = err
			}
		}
	}

	if len(errs) > 0 {
		return ErrBy(FailToRunInParallel{Errors: errs})
	}

	return Ok()
}

// Para is a function which creates a runner which runs multiple runners
// specified as arguments in parallel.
func Para(runners ...Runner) Runner {
	return paraRunner{runners: runners[:]}
}

// RunSeq is a function which runs specified runners sequencially.
func RunSeq(runners ...Runner) Err {
	for _, runner := range runners {
		err := runner.Run()
		if !err.IsOk() {
			return err
		}
	}
	return Ok()
}

// RunPara is a function which runs specified runners in parallel.
func RunPara(runners ...Runner) Err {
	ch := make(chan Err)

	for _, runner := range runners {
		go func(runner Runner, ch chan Err) {
			err := runner.Run()
			ch <- err
		}(runner, ch)
	}

	errs := make(map[string]Err)
	n := len(runners)
	for i := 0; i < n; i++ {
		select {
		case err := <-ch:
			if !err.IsOk() {
				errs[strconv.Itoa(i)] = err
			}
		}
	}

	if len(errs) > 0 {
		return ErrBy(FailToRunInParallel{Errors: errs})
	}

	return Ok()
}
