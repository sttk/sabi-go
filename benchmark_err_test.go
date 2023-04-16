package sabi_test

import (
	"github.com/sttk-go/sabi"
	"strconv"
	"testing"
)

func b_unused(v any) {}

func returnNilError() error {
	return nil
}
func BenchmarkError_nil(b *testing.B) {
	var err error
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnNilError()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func returnOkErr() sabi.Err {
	return sabi.Ok()
}
func BenchmarkErr_ok(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnOkErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_nil_isNil(b *testing.B) {
	var err error
	e := returnNilError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if e == nil {
			err = e
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkErr_ok_isOk(b *testing.B) {
	var err sabi.Err
	e := returnOkErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if e.IsOk() {
			err = e
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_nil_typeSwitch(b *testing.B) {
	var err error
	e := returnNilError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		switch e.(type) {
		case nil:
			err = e
		default:
			b.Fail()
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkErr_ok_typeSwitch(b *testing.B) {
	var err sabi.Err
	e := returnOkErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		switch e.Reason().(type) {
		case nil:
			err = e
		default:
			b.Fail()
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_nil_ErrorString(b *testing.B) {
	var str string
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := "nil"
		str = s
	}
	b.StopTimer()
	b_unused(str)
}

func BenchmarkErr_ok_ErrorString(b *testing.B) {
	var str string
	e := returnOkErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
}

type EmptyError struct {
}

func returnEmptyError() error {
	return EmptyError{}
}
func (e EmptyError) Error() string {
	return "EmptyError"
}
func BenchmarkError_empty(b *testing.B) {
	var err error
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnEmptyError()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

type EmptyReason struct {
}

func returnEmptyReasonedErr() sabi.Err {
	return sabi.NewErr(EmptyReason{})
}
func BenchmarkErr_emptyReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnEmptyReasonedErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_empty_isNotNil(b *testing.B) {
	var err error
	e := returnEmptyError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if e != nil {
			err = e
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkErr_emptyReason_isNotOk(b *testing.B) {
	var err sabi.Err
	e := returnEmptyReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !e.IsOk() {
			err = e
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_empty_typeSwitch(b *testing.B) {
	var err error
	e := returnEmptyError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		switch e.(type) {
		case EmptyError:
			err = e
		default:
			b.Fail()
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkErr_emptyReason_typeSwitch(b *testing.B) {
	var err sabi.Err
	e := returnEmptyReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		switch e.Reason().(type) {
		case EmptyReason:
			err = e
		default:
			b.Fail()
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_empty_ErrorString(b *testing.B) {
	var str string
	e := returnEmptyError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
}

func BenchmarkErr_emptyReason_ErrorString(b *testing.B) {
	var str string
	e := returnEmptyReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
	b_unused(e)
}

type OneFieldError struct {
	FieldA string
}

func (e OneFieldError) Error() string {
	return "OneFieldError{FieldA:" + e.FieldA + "}"
}
func returnOneFieldError() error {
	return OneFieldError{FieldA: "abc"}
}
func BenchmarkError_oneField(b *testing.B) {
	var err error
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnOneFieldError()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

type OneFieldReason struct {
	FieldA string
}

func returnOneFieldReasonedErr() sabi.Err {
	return sabi.NewErr(OneFieldReason{FieldA: "abc"})
}
func BenchmarkErr_oneFieldReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnOneFieldReasonedErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func returnOneFieldErrorPtr() error {
	return &OneFieldError{FieldA: "abc"}
}
func BenchmarkError_oneFieldPtr(b *testing.B) {
	var err error
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnOneFieldErrorPtr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func returnOneFieldReasonedPtrErr() sabi.Err {
	return sabi.NewErr(&OneFieldReason{FieldA: "abc"})
}
func BenchmarkErr_oneFieldReasonPtr(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnOneFieldReasonedPtrErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_oneField_isNotNil(b *testing.B) {
	var err error
	e := returnOneFieldError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if e != nil {
			err = e
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkErr_oneFieldReason_isNotOk(b *testing.B) {
	var err sabi.Err
	e := returnOneFieldReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !e.IsOk() {
			err = e
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_oneField_typeSwitch(b *testing.B) {
	var err error
	e := returnOneFieldError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		switch e.(type) {
		case OneFieldError:
			err = e
		default:
			b.Fail()
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkErr_oneFieldReason_typeSwitch(b *testing.B) {
	var err sabi.Err
	e := returnOneFieldReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		switch e.Reason().(type) {
		case OneFieldReason:
			err = e
		default:
			b.Fail()
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_oneField_ErrorString(b *testing.B) {
	var str string
	e := returnOneFieldError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
}

func BenchmarkErr_oneFieldReason_ErrorString(b *testing.B) {
	var str string
	e := returnOneFieldReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
}

type FiveFieldError struct {
	FieldA string
	FieldB int
	FieldC bool
	FieldD string
	FieldE string
}

func (e FiveFieldError) Error() string {
	return "FiveFieldError{FieldA:" + e.FieldA +
		",FieldB:" + strconv.Itoa(e.FieldB) +
		",FieldC:" + strconv.FormatBool(e.FieldC) +
		",FieldD:" + e.FieldD + ",FieldE:" + e.FieldE +
		"}"
}
func returnFiveFieldError() error {
	return FiveFieldError{
		FieldA: "abc", FieldB: 123, FieldC: true, FieldD: "def", FieldE: "ghi",
	}
}
func BenchmarkError_fiveField(b *testing.B) {
	var err error
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnFiveFieldError()
		err = e
	}
	b.StopTimer()
	b_unused(err)
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
func BenchmarkErr_fiveFieldReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnFiveFieldReasonedErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_fiveField_isNotNil(b *testing.B) {
	var err error
	e := returnFiveFieldError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if e != nil {
			err = e
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkErr_fiveFieldReason_isNotOk(b *testing.B) {
	var err sabi.Err
	e := returnFiveFieldReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !e.IsOk() {
			err = e
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_fiveField_typeSwitch(b *testing.B) {
	var err error
	e := returnFiveFieldError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		switch e.(type) {
		case FiveFieldError:
			err = e
		default:
			b.Fail()
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkErr_fiveFieldReason_typeSwitch(b *testing.B) {
	var err sabi.Err
	e := returnFiveFieldReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		switch e.Reason().(type) {
		case FiveFieldReason:
			err = e
		default:
			b.Fail()
		}
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_fiveField_ErrorString(b *testing.B) {
	var str string
	e := returnFiveFieldError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
}

func BenchmarkErr_fiveFieldReason_ErrorString(b *testing.B) {
	var str string
	e := returnFiveFieldReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
}

type HavingCauseError struct {
	Cause error
}

func (e HavingCauseError) Error() string {
	return "HavingCauseError{cause:" + e.Cause.Error() + "}"
}
func (e HavingCauseError) Unwrap() error {
	return e.Cause
}
func returnHavingCauseError() error {
	return HavingCauseError{Cause: EmptyError{}}
}
func BenchmarkError_havingCause(b *testing.B) {
	var err error
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnHavingCauseError()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

type HavingCauseReason struct {
}

func returnHavingCauseReasonedErr() sabi.Err {
	return sabi.NewErr(HavingCauseError{}, EmptyError{})
}
func BenchmarkErr_havingCauseReason(b *testing.B) {
	var err sabi.Err
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e := returnHavingCauseReasonedErr()
		err = e
	}
	b.StopTimer()
	b_unused(err)
}

func BenchmarkError_havingCause_ErrorString(b *testing.B) {
	var str string
	e := returnHavingCauseError()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
}

func BenchmarkErr_havingCauseReason_ErrorString(b *testing.B) {
	var str string
	e := returnHavingCauseReasonedErr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := e.Error()
		str = s
	}
	b.StopTimer()
	b_unused(str)
}
