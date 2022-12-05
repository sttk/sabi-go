// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

// Proc is a structure type which represents a procedure.
type Proc[D any] struct {
	daxBase *DaxBase
	dax     D
}

// NewProc is a function which create a new Proc.
func NewProc[D any](daxBase *DaxBase, dax D) Proc[D] {
	return Proc[D]{daxBase: daxBase, dax: dax}
}

// AddLocalDaxSrc is a method which registers a procedure-local DaxSrc
// with a specified name.
func (proc Proc[D]) AddLocalDaxSrc(name string, ds DaxSrc) {
	proc.daxBase.AddLocalDaxSrc(name, ds)
}

// RunTxn is a method which runs logic functions specified as arguments in a
// transaction.
func (proc Proc[D]) RunTxn(logics ...func(dax D) Err) Err {
	proc.daxBase.begin()

	err := Ok()

	for _, logic := range logics {
		err = logic(proc.dax)
		if !err.IsOk() {
			break
		}
	}

	if err.IsOk() {
		err = proc.daxBase.commit()
	}

	if !err.IsOk() {
		proc.daxBase.rollback()
	}

	proc.daxBase.close()

	return err
}

// Txn is a method which creates a transaction having specified logic
// functions.
func (proc Proc[D]) Txn(logics ...func(dax D) Err) Runner {
	return txnRunner[D]{
		logics:  logics,
		daxBase: proc.daxBase,
		dax:     proc.dax,
	}
}

type txnRunner[D any] struct {
	logics  []func(D) Err
	daxBase *DaxBase
	dax     D
}

func (txn txnRunner[D]) Run() Err {
	txn.daxBase.begin()

	err := Ok()

	for _, logic := range txn.logics {
		err = logic(txn.dax)
		if !err.IsOk() {
			break
		}
	}

	if err.IsOk() {
		err = txn.daxBase.commit()
	}

	if !err.IsOk() {
		txn.daxBase.rollback()
	}

	txn.daxBase.close()

	return err
}
