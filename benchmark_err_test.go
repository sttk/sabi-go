package sabi_test

import (
	"errors"
	"github.com/sttk-go/sabi"
	"testing"
)

type /* error reason */ (
	ReasonForBenchWithNoParam    struct{}
	ReasonForBenchWithManyParams struct {
		Param1, Param2, Param3, Param4, Param5  string
		Param6, Param7, Param8, Param9, Param10 int
	}
)

func _err_unused(v interface{}) {
}

func BenchmarkErr_NewErr_empty(b *testing.B) {
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		err := sabi.NewErr(ReasonForBenchWithNoParam{})
		_err_unused(err)
	}

	b.StopTimer()
}

func BenchmarkErr_NewErr_manyParams(b *testing.B) {
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		err := sabi.NewErr(ReasonForBenchWithManyParams{
			Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			Param2:  "abcdefghijklmnopqrstuvwxyz",
			Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
			Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
			Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
			Param6:  1234567890,
			Param7:  9876543210,
			Param8:  1111111111,
			Param9:  2222222222,
			Param10: 3333333333,
		})
		_err_unused(err)
	}

	b.StopTimer()
}

func ProcForBenchByVal() sabi.Err {
	err := sabi.NewErr(ReasonForBenchWithManyParams{
		Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Param2:  "abcdefghijklmnopqrstuvwxyz",
		Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
		Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
		Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
		Param6:  1234567890,
		Param7:  9876543210,
		Param8:  1111111111,
		Param9:  2222222222,
		Param10: 3333333333,
	})
	return err
}

func BenchmarkErr_NewErr_byValue(b *testing.B) {
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		err := ProcForBenchByVal()
		_err_unused(err)
	}

	b.StopTimer()
}

func ProcForBenchByPtr() *sabi.Err {
	err := sabi.NewErr(ReasonForBenchWithManyParams{
		Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Param2:  "abcdefghijklmnopqrstuvwxyz",
		Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
		Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
		Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
		Param6:  1234567890,
		Param7:  9876543210,
		Param8:  1111111111,
		Param9:  2222222222,
		Param10: 3333333333,
	})
	return &err
}

func BenchmarkErr_NewErr_byPtr(b *testing.B) {
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		err := ProcForBenchByPtr()
		_err_unused(err)
	}

	b.StopTimer()
}

func BenchmarkErr_Reason_emtpy(b *testing.B) {
	err := sabi.NewErr(ReasonForBenchWithNoParam{})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		reason := err.Reason()
		_err_unused(reason)
	}

	b.StopTimer()
}

func BenchmarkErr_Reason_manyParams(b *testing.B) {
	err := sabi.NewErr(ReasonForBenchWithManyParams{
		Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Param2:  "abcdefghijklmnopqrstuvwxyz",
		Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
		Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
		Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
		Param6:  1234567890,
		Param7:  9876543210,
		Param8:  1111111111,
		Param9:  2222222222,
		Param10: 3333333333,
	})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		reason := err.Reason()
		_err_unused(reason)
	}

	b.StopTimer()
}

func BenchmarkErr_Reason_type_emtpy(b *testing.B) {
	err := sabi.NewErr(ReasonForBenchWithNoParam{})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		switch err.Reason().(type) {
		case ReasonForBenchWithNoParam, *ReasonForBenchWithNoParam:
		}
	}

	b.StopTimer()
}

func BenchmarkErr_Reason_type_manyParams(b *testing.B) {
	err := sabi.NewErr(ReasonForBenchWithManyParams{
		Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Param2:  "abcdefghijklmnopqrstuvwxyz",
		Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
		Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
		Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
		Param6:  1234567890,
		Param7:  9876543210,
		Param8:  1111111111,
		Param9:  2222222222,
		Param10: 3333333333,
	})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		switch err.Reason().(type) {
		case ReasonForBenchWithManyParams, *ReasonForBenchWithManyParams:
		}
	}

	b.StopTimer()
}

func BenchmarkErr_ReasonName(b *testing.B) {
	err := sabi.NewErr(ReasonForBenchWithManyParams{
		Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Param2:  "abcdefghijklmnopqrstuvwxyz",
		Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
		Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
		Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
		Param6:  1234567890,
		Param7:  9876543210,
		Param8:  1111111111,
		Param9:  2222222222,
		Param10: 3333333333,
	})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		name := err.ReasonName()
		_err_unused(name)
	}

	b.StopTimer()
}

func BenchmarkErr_Cause(b *testing.B) {
	cause := errors.New("Causal error")
	err := sabi.NewErr(ReasonForBenchWithManyParams{
		Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Param2:  "abcdefghijklmnopqrstuvwxyz",
		Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
		Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
		Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
		Param6:  1234567890,
		Param7:  9876543210,
		Param8:  1111111111,
		Param9:  2222222222,
		Param10: 3333333333,
	}, cause)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		cause := err.Cause()
		_err_unused(cause)
	}

	b.StopTimer()
}

func BenchmarkErr_Situation(b *testing.B) {
	err := sabi.NewErr(ReasonForBenchWithManyParams{
		Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Param2:  "abcdefghijklmnopqrstuvwxyz",
		Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
		Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
		Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
		Param6:  1234567890,
		Param7:  9876543210,
		Param8:  1111111111,
		Param9:  2222222222,
		Param10: 3333333333,
	})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		m := err.Situation()
		_err_unused(m)
	}

	b.StopTimer()
}

func BenchmarkErr_Get(b *testing.B) {
	err := sabi.NewErr(ReasonForBenchWithManyParams{
		Param1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Param2:  "abcdefghijklmnopqrstuvwxyz",
		Param3:  "8b91114c-f620-46bc-a991-25c7ac0f7935",
		Param4:  "f59de3cb-b91c-4c1d-a152-787173d1ab9b",
		Param5:  "c50c24cd-fe6f-4de7-803e-193d705376b7",
		Param6:  1234567890,
		Param7:  9876543210,
		Param8:  1111111111,
		Param9:  2222222222,
		Param10: 3333333333,
	})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		s1 := err.Get("Param1").(string)
		s2 := err.Get("Param2").(string)
		s3 := err.Get("Param3").(string)
		s4 := err.Get("Param4").(string)
		s5 := err.Get("Param5").(string)
		n1 := err.Get("Param6").(int)
		n2 := err.Get("Param7").(int)
		n3 := err.Get("Param8").(int)
		n4 := err.Get("Param9").(int)
		n5 := err.Get("Param10").(int)
		_err_unused(s1)
		_err_unused(s2)
		_err_unused(s3)
		_err_unused(s4)
		_err_unused(s5)
		_err_unused(n1)
		_err_unused(n2)
		_err_unused(n3)
		_err_unused(n4)
		_err_unused(n5)
	}

	b.StopTimer()
}
