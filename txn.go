// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

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
