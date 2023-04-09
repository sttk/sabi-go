package sabi_test

import (
	"github.com/sttk-go/sabi"
)

type BazDax interface {
	GetData() string
	SetData(data string)
}

var base0 = sabi.NewDaxBase()

var base = struct {
	sabi.DaxBase
	FooGetterDax
	BarSetterDax
}{
	DaxBase:      base0,
	FooGetterDax: FooGetterDax{Dax: base0},
	BarSetterDax: BarSetterDax{Dax: base0},
}

func ExamplePara() {
	txn1 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })
	txn2 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })

	paraRunner := sabi.Para(txn1, txn2)

	err := sabi.RunSeq(paraRunner)

	// Output:

	unused(err)
	sabi.Clear()
}

func ExampleSeq() {
	txn1 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })
	txn2 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })

	seqRunner := sabi.Seq(txn1, txn2)

	err := sabi.RunSeq(seqRunner)

	// Output:

	unused(err)
	sabi.Clear()
}

func ExampleRunPara() {
	txn1 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })
	txn2 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })

	err := sabi.RunPara(txn1, txn2)

	// Output:

	unused(err)
	sabi.Clear()
}

func ExampleRunSeq() {
	txn1 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })
	txn2 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })

	err := sabi.RunSeq(txn1, txn2)

	// Output:

	unused(err)
	sabi.Clear()
}
