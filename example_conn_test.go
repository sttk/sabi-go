package sabi_test

import (
	"fmt"
	"github.com/sttk-go/sabi"
	"reflect"
)

func unused(v interface{}) {}

type FooConn struct {
}

func (conn *FooConn) Commit() sabi.Err {
	return sabi.Ok()
}

func (conn *FooConn) Rollback() {
}

func (conn *FooConn) Close() {
}

type FooConnCfg struct {
}

func (cfg FooConnCfg) CreateConn() (sabi.Conn, sabi.Err) {
	return &FooConn{}, sabi.Ok()
}

type BarConn struct {
}

func (conn *BarConn) Commit() sabi.Err {
	return sabi.Ok()
}

func (conn *BarConn) Rollback() {
}

func (conn *BarConn) Close() {
}

type BarConnCfg struct {
}

func (cfg BarConnCfg) CreateConn() (sabi.Conn, sabi.Err) {
	return &BarConn{}, sabi.Ok()
}

func ExampleAddGlobalConnCfg() {
	sabi.AddGlobalConnCfg("foo", FooConnCfg{})
	sabi.AddGlobalConnCfg("bar", BarConnCfg{})

	base := sabi.NewConnBase()

	type FooBarDax struct {
		sabi.Dax
	}

	dax := FooBarDax{Dax: base}

	conn, err := dax.GetConn("foo")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	conn, err = dax.GetConn("bar")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	// Output:
	// conn = *sabi_test.FooConn
	// err.IsOk() = true
	// conn = *sabi_test.BarConn
	// err.IsOk() = true

	sabi.Clear()
}

func ExampleFixGlobalConnCfgs() {
	sabi.AddGlobalConnCfg("foo", FooConnCfg{})
	sabi.FixGlobalConnCfgs()
	sabi.AddGlobalConnCfg("bar", BarConnCfg{}) // Bad example

	base := sabi.NewConnBase()

	type FooBarDax struct {
		sabi.Dax
	}

	dax := FooBarDax{Dax: base}

	conn, err := dax.GetConn("foo")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	conn, err = dax.GetConn("bar")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())
	fmt.Printf("err.Error() = %v\n", err.Error())

	// Output:
	// conn = *sabi_test.FooConn
	// err.IsOk() = true
	// conn = <nil>
	// err.IsOk() = false
	// err.Error() = {reason=ConnCfgIsNotFound, Name=bar}

	sabi.Clear()
}

func ExampleNewConnBase() {
	base := sabi.NewConnBase()

	// Output:
	unused(base)
}

func ExampleConnBase_AddLocalConnCfg() {
	base := sabi.NewConnBase()
	base.AddLocalConnCfg("foo", FooConnCfg{})

	type FooBarDax struct {
		sabi.Dax
	}

	dax := FooBarDax{Dax: base}

	conn, err := dax.GetConn("foo")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	// Output:
	// conn = *sabi_test.FooConn
	// err.IsOk() = true

	sabi.Clear()
}
