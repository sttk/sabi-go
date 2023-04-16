package sabi_test

import (
	"github.com/sttk-go/sabi"
	"testing"
)

func BenchmarkNotify_addErrHandler(b *testing.B) {
	b.StartTimer()
	//sabi.AddSyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {})
	sabi.AddAsyncErrHandler(func(err sabi.Err, occ sabi.ErrOccasion) {})
	sabi.FixErrCfgs()
	b.StopTimer()
}

func BenchmarkNotify_ok(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnOkErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkNotify_emptyReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnEmptyReasonedErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkNotify_oneFieldReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnOneFieldReasonedErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkNotify_fiveFieldReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnFiveFieldReasonedErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkNotify_havingCauseReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnHavingCauseReasonedErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}
