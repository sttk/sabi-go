package sabi_test

import (
	"github.com/sttk/sabi"
	"github.com/sttk/sabi/errs"
)

type SyncDaxSrc struct{}

func (ds SyncDaxSrc) Setup(ag sabi.AsyncGroup) errs.Err {
	// ...
	return errs.Ok()
}

func (ds SyncDaxSrc) Close() {}
func (ds SyncDaxSrc) CreateDaxConn() (sabi.DaxConn, errs.Err) {
	return nil, errs.Ok()
}

type AsyncDaxSrc struct{}

func (ds AsyncDaxSrc) Setup(ag sabi.AsyncGroup) errs.Err {
	ag.Add(func() errs.Err {
		// ...
		return errs.Ok()
	})
	return errs.Ok()
}

func (ds AsyncDaxSrc) Close() {}
func (ds AsyncDaxSrc) CreateDaxConn() (sabi.DaxConn, errs.Err) {
	return nil, errs.Ok()
}

func ExampleAsyncGroup() {
	sabi.Uses("sync", SyncDaxSrc{})
	sabi.Uses("async", AsyncDaxSrc{})

	err := sabi.Setup()
	if err.IsNotOk() {
		return
	}
}
