package sabi_test

import (
	"fmt"
	"github.com/sttk-go/sabi"
	"reflect"
)

func unused(v interface{}) {}

type FooDaxConn struct {
}

func (conn *FooDaxConn) Commit() sabi.Err {
	return sabi.Ok()
}
func (conn *FooDaxConn) Rollback() {
}
func (conn *FooDaxConn) Close() {
}

type FooDaxSrc struct {
}

func (ds FooDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return &FooDaxConn{}, sabi.Ok()
}

type BarDaxConn struct {
}

func (conn *BarDaxConn) Commit() sabi.Err {
	return sabi.Ok()
}
func (conn *BarDaxConn) Rollback() {
}
func (conn *BarDaxConn) Close() {
}

type BarDaxSrc struct {
}

func (ds BarDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return &BarDaxConn{}, sabi.Ok()
}

type FooGetterDax struct {
	sabi.Dax
}

func (dax FooGetterDax) GetData() string {
	return "hello"
}

type BarSetterDax struct {
	sabi.Dax
}

func (dax BarSetterDax) SetData(data string) {
}

func ExampleAddGlobalDaxSrc() {
	sabi.AddGlobalDaxSrc("hoge", FooDaxSrc{})
	sabi.AddGlobalDaxSrc("fuga", BarDaxSrc{})

	base := sabi.NewDaxBase()

	type FooBarDax struct {
		sabi.Dax
	}

	dax := FooBarDax{Dax: base}

	conn, err := dax.GetDaxConn("hoge")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	conn, err = dax.GetDaxConn("fuga")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	// Output:
	// conn = *sabi_test.FooDaxConn
	// err.IsOk() = true
	// conn = *sabi_test.BarDaxConn
	// err.IsOk() = true

	sabi.Clear()
}

func ExampleFixGlobalDaxSrcs() {
	sabi.AddGlobalDaxSrc("hoge", FooDaxSrc{})
	sabi.FixGlobalDaxSrcs()
	sabi.AddGlobalDaxSrc("fuga", BarDaxSrc{})

	base := sabi.NewDaxBase()

	type FooBarDax struct {
		sabi.Dax
	}

	dax := FooBarDax{Dax: base}

	conn, err := dax.GetDaxConn("hoge")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	conn, err = dax.GetDaxConn("fuga")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())
	fmt.Printf("err.Error() = %v\n", err.Error())

	// Output:
	// conn = *sabi_test.FooDaxConn
	// err.IsOk() = true
	// conn = <nil>
	// err.IsOk() = false
	// err.Error() = {reason=DaxSrcIsNotFound, Name=fuga}

	sabi.Clear()
}

func ExampleNewDaxBase() {
	base := sabi.NewDaxBase()

	// Output:
	unused(base)
}

func ExampleDaxBase_AddLocalDaxSrc() {
	base := sabi.NewDaxBase()
	base.AddLocalDaxSrc("hoge", FooDaxSrc{})

	type FooBarDax struct {
		sabi.Dax
	}

	dax := FooBarDax{Dax: base}

	conn, err := dax.GetDaxConn("hoge")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	// Output:
	// conn = *sabi_test.FooDaxConn
	// err.IsOk() = true

	sabi.Clear()
}

func ExampleDax() {
	base0 := sabi.NewDaxBase()

	type MyDax interface {
		GetData() string
		SetData(data string)
	}

	base := struct {
		sabi.DaxBase
		FooGetterDax
		BarSetterDax
	}{
		DaxBase:      base0,
		FooGetterDax: FooGetterDax{Dax: base0},
		BarSetterDax: BarSetterDax{Dax: base0},
	}

	sabi.RunTxn(base, func(dax MyDax) sabi.Err {
		data := dax.GetData()
		dax.SetData(data)
		return sabi.Ok()
	})

	// Output:

	sabi.Clear()
}
