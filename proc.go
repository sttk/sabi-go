// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

// Proc is a structure type which represents a procedure.
type Proc[D any] struct {
	connBase *ConnBase
	dax      D
}

// NewProc is a function which create a new Proc.
func NewProc[D any](connBase *ConnBase, dax D) Proc[D] {
	return Proc[D]{connBase: connBase, dax: dax}
}

// AddLocalConnCfg is a method which registers a procedure-local ConnCfg
// with a specified name.
func (proc Proc[D]) AddLocalConnCfg(name string, cfg ConnCfg) {
	proc.connBase.AddLocalConnCfg(name, cfg)
}

// RunTxn is a method which runs logic functions specified as arguments in a
// transaction.
func (proc Proc[D]) RunTxn(logics ...func(dax D) Err) Err {
	proc.connBase.begin()

	err := Ok()

	for _, logic := range logics {
		err = logic(proc.dax)
		if !err.IsOk() {
			break
		}
	}

	if err.IsOk() {
		err = proc.connBase.commit()
	}

	if !err.IsOk() {
		proc.connBase.rollback()
	}

	proc.connBase.close()

	return err
}

// NewTxn is a method which creates a transaction having specified logic
// functions.
func (proc Proc[D]) Txn(logics ...func(dax D) Err) Runner {
	return txnRunner[D]{
		logics:   logics,
		connBase: proc.connBase,
		dax:      proc.dax,
	}
}
