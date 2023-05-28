package sabi_test

import (
	"github.com/sttk-go/sabi"
	"testing"
)

func unused(v any) {}

func returnOkErr() sabi.Err {
	return sabi.Ok()
}

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
	unused(err)
}

type EmptyReason struct {
}

func returnEmptyReasonedErr() sabi.Err {
	return sabi.NewErr(EmptyReason{})
}

func BenchmarkNotify_emptyReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnEmptyReasonedErr()
		err = e
	}
	b.StopTimer()
	unused(err)
}

type OneFieldReason struct {
	FieldA string
}

func returnOneFieldReasonedErr() sabi.Err {
	return sabi.NewErr(OneFieldReason{FieldA: "abc"})
}

func BenchmarkNotify_oneFieldReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnOneFieldReasonedErr()
		err = e
	}
	b.StopTimer()
	unused(err)
}

type FiveFieldReason struct {
	FieldA string
	FieldB int
	FieldC bool
	FieldD string
	FieldE string
}

func returnFiveFieldReasonedErr() sabi.Err {
	return sabi.NewErr(FiveFieldReason{
		FieldA: "abc", FieldB: 123, FieldC: true, FieldD: "def", FieldE: "ghi",
	})
}

func BenchmarkNotify_fiveFieldReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnFiveFieldReasonedErr()
		err = e
	}
	b.StopTimer()
	unused(err)
}

type EmptyError struct {
}

func (e EmptyError) Error() string {
	return "EmptyError"
}

type HavingCauseError struct {
	Cause error
}

func (e HavingCauseError) Error() string {
	return "HavingCauseError{cause:" + e.Cause.Error() + "}"
}

func returnHavingCauseReasonedErr() sabi.Err {
	return sabi.NewErr(HavingCauseError{}, EmptyError{})
}

func BenchmarkNotify_havingCauseReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnHavingCauseReasonedErr()
		err = e
	}
	b.StopTimer()
	unused(err)
}
