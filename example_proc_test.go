package sabi_test

import (
	"github.com/sttk-go/sabi"
)

type FooGetterDax struct {
	sabi.Dax
}

func (dax FooGetterDax) GetData() string {
	return "hello"
}

type BarSetterDax struct {
	sabi.Dax
}

func (dax BarSetterDax) SetData(data string) {
}

func ExampleNewProc() {
	sabi.AddGlobalConnCfg("foo", FooConnCfg{})
	sabi.AddGlobalConnCfg("bar", BarConnCfg{})
	sabi.FixGlobalConnCfgs()

	base := sabi.NewConnBase()

	type MyDax interface {
		GetData() string
		SetData(data string)
	}

	dax := struct {
		FooGetterDax
		BarSetterDax
	}{
		FooGetterDax: FooGetterDax{Dax: base},
		BarSetterDax: BarSetterDax{Dax: base},
	}

	proc := sabi.NewProc[MyDax](base, dax)
	proc.RunTxn(func(dax MyDax) sabi.Err {
		data := dax.GetData()
		dax.SetData(data)
		return sabi.Ok()
	})

	// Output:

	sabi.Clear()
}

func ExampleProc_AddLocalConnCfg() {
	base := sabi.NewConnBase()

	type MyDax interface {
		GetData() string
		SetData(data string)
	}

	dax := struct {
		FooGetterDax
		BarSetterDax
	}{
		FooGetterDax: FooGetterDax{Dax: base},
		BarSetterDax: BarSetterDax{Dax: base},
	}

	proc := sabi.NewProc[MyDax](base, dax)

	proc.AddLocalConnCfg("foo", FooConnCfg{})
	proc.AddLocalConnCfg("bar", BarConnCfg{})

	proc.RunTxn(func(dax MyDax) sabi.Err {
		data := dax.GetData()
		dax.SetData(data)
		return sabi.Ok()
	})

	// Output:

	sabi.Clear()
}

func ExampleProc_RunTxn() {
	base := sabi.NewConnBase()

	type MyDax interface {
		GetData() string
		SetData(data string)
	}

	dax := struct {
		FooGetterDax
		BarSetterDax
	}{
		FooGetterDax: FooGetterDax{Dax: base},
		BarSetterDax: BarSetterDax{Dax: base},
	}

	proc := sabi.NewProc[MyDax](base, dax)

	err := proc.RunTxn(func(dax MyDax) sabi.Err {
		data := dax.GetData()
		dax.SetData(data)
		return sabi.Ok()
	})

	// Output:

	unused(err)
	sabi.Clear()
}

func ExampleProc_Txn() {
	base := sabi.NewConnBase()

	type MyDax interface {
		GetData() string
		SetData(data string)
	}

	dax := struct {
		FooGetterDax
		BarSetterDax
	}{
		FooGetterDax: FooGetterDax{Dax: base},
		BarSetterDax: BarSetterDax{Dax: base},
	}

	proc := sabi.NewProc[MyDax](base, dax)

	txn := proc.Txn(func(dax MyDax) sabi.Err {
		data := dax.GetData()
		dax.SetData(data)
		return sabi.Ok()
	})

	err := sabi.RunSeq(txn)

	// Output:

	unused(err)
	sabi.Clear()
}
