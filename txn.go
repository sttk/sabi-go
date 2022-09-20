// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

type txnRunner[D any] struct {
	logics   []func(D) Err
	connBase *ConnBase
	dax      D
}

func (txn txnRunner[D]) Run() Err {
	txn.connBase.begin()

	err := Ok()

	for _, logic := range txn.logics {
		err = logic(txn.dax)
		if !err.IsOk() {
			break
		}
	}

	if err.IsOk() {
		err = txn.connBase.commit()
	}

	if !err.IsOk() {
		txn.connBase.rollback()
	}

	txn.connBase.close()

	return err
}
