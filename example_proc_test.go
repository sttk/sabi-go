package sabi_test

import (
	"github.com/sttk-go/sabi"
)

func ExampleNewProc() {
	sabi.AddGlobalDaxSrc("foo", FooDaxSrc{})
	sabi.AddGlobalDaxSrc("bar", BarDaxSrc{})
	sabi.FixGlobalDaxSrcs()

	base := sabi.NewDaxBase()

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

func ExampleProc_AddLocalDaxSrc() {
	base := sabi.NewDaxBase()

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

	proc.AddLocalDaxSrc("foo", FooDaxSrc{})
	proc.AddLocalDaxSrc("bar", BarDaxSrc{})

	proc.RunTxn(func(dax MyDax) sabi.Err {
		data := dax.GetData()
		dax.SetData(data)
		return sabi.Ok()
	})

	// Output:

	sabi.Clear()
}

func ExampleProc_RunTxn() {
	base := sabi.NewDaxBase()

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
	base := sabi.NewDaxBase()

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
