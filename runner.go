// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"github.com/sttk/sabi/errs"
)

type /* error reasons */ (
	// FailToRunInParallel is an error reason which indicates that some of runner
	// functions running in parallel failed.
	FailToRunInParallel struct {
		Errors map[int]errs.Err
	}
)

// Seq is the function which runs argument functions sequencially.
func Seq(runners ...func() errs.Err) errs.Err {
	for _, runner := range runners {
		err := runner()
		if err.IsNotOk() {
			return err
		}
	}

	return errs.Ok()
}

// Seq_ is the function which creates a runner function which runs Seq
// function.
func Seq_(runners ...func() errs.Err) func() errs.Err {
	return func() errs.Err {
		return Seq(runners...)
	}
}

// Para is the function which runs argument functions in parallel.
func Para(runners ...func() errs.Err) errs.Err {
	var ag asyncGroupAsync[int]

	for i, runner := range runners {
		ag.name = i
		ag.Add(runner)
	}

	ag.wait()

	if ag.hasErr() {
		return errs.New(FailToRunInParallel{Errors: ag.makeErrs()})
	}

	return errs.Ok()
}

// Para_ is the function which creates a runner function which runs Para
// function.
func Para_(runners ...func() errs.Err) func() errs.Err {
	return func() errs.Err {
		return Para(runners...)
	}
}
