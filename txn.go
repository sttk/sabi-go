// Copyright (C) 2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

// RunTxn is a function which runs logic functions specified as arguments in a
// transaction.
func RunTxn[D any](base DaxBase, logics ...func(dax D) Err) Err {
	base.begin()

	dax := base.(D)
	err := Ok()

	for _, logic := range logics {
		err = logic(dax)
		if !err.IsOk() {
			break
		}
	}

	if err.IsOk() {
		err = base.commit()
	}

	if !err.IsOk() {
		base.rollback()
	}

	base.end()

	return err
}

// Txn is a function which creates a transaction having specified logic
// functions.
func Txn[D any](base DaxBase, logics ...func(dax D) Err) Runner {
	return txnRunner[D]{
		base:   base,
		logics: logics,
	}
}

type txnRunner[D any] struct {
	base   DaxBase
	logics []func(D) Err
}

func (txn txnRunner[D]) Run() Err {
	return RunTxn(txn.base, txn.logics...)
}
