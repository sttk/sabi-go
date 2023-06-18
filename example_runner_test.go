package sabi_test

import (
	"github.com/sttk/sabi"
)

func unused(v any) {}

type BazDax interface {
	GetData() string
	SetData(data string)
}

type bazDaxImpl struct {
}

func (dax bazDaxImpl) GetData() string {
	return ""
}
func (dax bazDaxImpl) SetData(data string) {
}

var base0 = sabi.NewDaxBase()

var base = struct {
	sabi.DaxBase
	bazDaxImpl
}{
	DaxBase:    base0,
	bazDaxImpl: bazDaxImpl{},
}

func ExamplePara() {
	txn1 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })
	txn2 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })

	paraRunner := sabi.Para(txn1, txn2)

	err := paraRunner()

	// Output:

	unused(err)
}

func ExampleSeq() {
	txn1 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })
	txn2 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })

	seqRunner := sabi.Seq(txn1, txn2)

	err := seqRunner()

	// Output:

	unused(err)
}

func ExampleRunPara() {
	txn1 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })
	txn2 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })

	err := sabi.RunPara(txn1, txn2)

	// Output:

	unused(err)
}

func ExampleRunSeq() {
	txn1 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })
	txn2 := sabi.Txn(base, func(dax BazDax) sabi.Err { return sabi.Ok() })

	err := sabi.RunSeq(txn1, txn2)

	// Output:

	unused(err)
}
