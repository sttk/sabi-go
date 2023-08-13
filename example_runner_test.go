package sabi_test

import (
	"github.com/sttk/sabi"
	"github.com/sttk/sabi/errs"
)

func NewMyDaxBase() sabi.DaxBase {
	return sabi.NewDaxBase()
}

func ExampleSeq() {
	type FooDax interface {
		sabi.Dax
		// ...
	}

	type BarDax interface {
		sabi.Dax
		// ...
	}

	base := NewMyDaxBase()

	txn1 := sabi.Txn_[FooDax](base, func(dax FooDax) errs.Err {
		// ...
		return errs.Ok()
	})

	txn2 := sabi.Txn_[BarDax](base, func(dax BarDax) errs.Err {
		// ...
		return errs.Ok()
	})

	err := sabi.Seq(txn1, txn2)
	if err.IsNotOk() {
		// ...
	}
}

func ExamplePara() {
	type FooDax interface {
		sabi.Dax
		// ...
	}

	type BarDax interface {
		sabi.Dax
		// ...
	}

	base := NewMyDaxBase()

	txn1 := sabi.Txn_[FooDax](base, func(dax FooDax) errs.Err {
		// ...
		return errs.Ok()
	})

	txn2 := sabi.Txn_[BarDax](base, func(dax BarDax) errs.Err {
		// ...
		return errs.Ok()
	})

	err := sabi.Para(txn1, txn2)
	if err.IsNotOk() {
		// ...
	}
}
