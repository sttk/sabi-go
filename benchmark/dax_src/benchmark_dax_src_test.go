package sabi_test

import (
	"github.com/sttk-go/sabi"
	"testing"
)

type FooDaxConn struct {}

func (conn FooDaxConn) Commit() sabi.Err { return sabi.Ok() }
func (conn FooDaxConn) Rollback() {}
func (conn FooDaxConn) Close() {}

type FooDaxSrc struct {}
func (ds FooDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return FooDaxConn{}, sabi.Ok()
}
func (ds FooDaxSrc) SetUp() sabi.Err { return sabi.Ok() }
func (ds FooDaxSrc) End() {}

func BenchmarkDaxSrc_AddGlobalDaxSrc(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sabi.AddGlobalDaxSrc("foo", FooDaxSrc{})
	}
	b.StopTimer()
}

func BenchmarkDaxSrc_AddGlobalDaxSrcPointer(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sabi.AddGlobalDaxSrc("foo", &FooDaxSrc{})
	}
	b.StopTimer()
}

func NewFooDaxSrc() FooDaxSrc {
	return FooDaxSrc{}
}

func BenchmarkDaxSrc_AddGlobalDaxSrc_withNewFunction(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sabi.AddGlobalDaxSrc("foo", NewFooDaxSrc())
	}
	b.StopTimer()
}

func BenchmarkDaxSrc_StartUpGlobalDaxSrcs(b *testing.B) {
	sabi.AddGlobalDaxSrc("foo", FooDaxSrc{})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sabi.StartUpGlobalDaxSrcs()
	}
	b.StopTimer()
}
