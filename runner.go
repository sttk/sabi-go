// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

type /* error reasons */ (
	// FailToRunInParallel is an error reason which indicates some runner
	// functions which run in parallel failed.
	FailToRunInParallel struct {
		Errors map[int]Err
	}
)

// RunSeq is a function which runs specified runner functions sequencially.
func RunSeq(runners ...func() Err) Err {
	for _, runner := range runners {
		err := runner()
		if err.IsNotOk() {
			return err
		}
	}
	return Ok()
}

// Seq is a function which creates a runner function which runs multiple
// runner functions specified as arguments sequencially.
func Seq(runners ...func() Err) func() Err {
	return func() Err {
		return RunSeq(runners...)
	}
}

type indexedErr struct {
	index int
	err   Err
}

func RunPara(runners ...func() Err) Err {
	ch := make(chan indexedErr)

	for i, r := range runners {
		go func(index int, runner func() Err, ch chan indexedErr) {
			err := runner()
			ie := indexedErr{index: index, err: err}
			ch <- ie
		}(i, r, ch)
	}

	errs := make(map[int]Err)
	n := len(runners)
	for i := 0; i < n; i++ {
		select {
		case ie := <-ch:
			if ie.err.IsNotOk() {
				errs[ie.index] = ie.err
			}
		}
	}

	if len(errs) > 0 {
		return NewErr(FailToRunInParallel{Errors: errs})
	}

	return Ok()
}

// Para is a function which creates a runner function which runs multiple
// runner functions specified as arguments in parallel.
func Para(runners ...func() Err) func() Err {
	return func() Err {
		return RunPara(runners...)
	}
}
