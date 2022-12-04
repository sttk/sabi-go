package sabi_test

import (
	"github.com/sttk-go/sabi"
)

type BazDax interface {
	GetData() string
	SetData(data string)
}

var base = sabi.NewDaxBase()

var dax = struct {
	FooGetterDax
	BarSetterDax
}{
	FooGetterDax: FooGetterDax{Dax: base},
	BarSetterDax: BarSetterDax{Dax: base},
}

func ExamplePara() {
	proc := sabi.NewProc[BazDax](base, dax)
	txn1 := proc.Txn(func(dax BazDax) sabi.Err { /* ... */ return sabi.Ok() })
	txn2 := proc.Txn(func(dax BazDax) sabi.Err { /* ... */ return sabi.Ok() })

	paraRunner := sabi.Para(txn1, txn2)

	err := sabi.RunSeq(paraRunner)

	// Output:

	unused(err)
	sabi.Clear()
}

func ExampleSeq() {
	proc := sabi.NewProc[BazDax](base, dax)
	txn1 := proc.Txn(func(dax BazDax) sabi.Err { /* ... */ return sabi.Ok() })
	txn2 := proc.Txn(func(dax BazDax) sabi.Err { /* ... */ return sabi.Ok() })

	seqRunner := sabi.Seq(txn1, txn2)

	err := sabi.RunSeq(seqRunner)

	// Output:

	unused(err)
	sabi.Clear()
}

func ExampleRunPara() {
	proc := sabi.NewProc[BazDax](base, dax)
	txn1 := proc.Txn(func(dax BazDax) sabi.Err { /* ... */ return sabi.Ok() })
	txn2 := proc.Txn(func(dax BazDax) sabi.Err { /* ... */ return sabi.Ok() })

	err := sabi.RunPara(txn1, txn2)

	// Output:

	unused(err)
	sabi.Clear()
}

func ExampleRunSeq() {
	proc := sabi.NewProc[BazDax](base, dax)
	txn1 := proc.Txn(func(dax BazDax) sabi.Err { /* ... */ return sabi.Ok() })
	txn2 := proc.Txn(func(dax BazDax) sabi.Err { /* ... */ return sabi.Ok() })

	err := sabi.RunSeq(txn1, txn2)

	// Output:

	unused(err)
	sabi.Clear()
}
