package sabi_test

import (
	"fmt"
	"github.com/sttk-go/sabi"
	"reflect"
)

func ExampleAddGlobalDaxSrc() {
	sabi.AddGlobalDaxSrc("hoge", NewMapDaxSrc())

	if sabi.StartUpGlobalDaxSrcs().IsNotOk() {
		return
	}
	defer sabi.ShutdownGlobalDaxSrcs()

	base := sabi.NewDaxBase()

	type MyDax struct {
		sabi.Dax
	}

	dax := MyDax{Dax: base}

	conn, err := sabi.GetDaxConn[MapDaxConn](dax, "hoge")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %t\n", err.IsOk())

	// Output:
	// conn = sabi_test.MapDaxConn
	// err.IsOk() = true

	ClearDaxBase()
}

func ExampleStartUpGlobalDaxSrcs() {
	sabi.AddGlobalDaxSrc("hoge", NewMapDaxSrc())

	if err := sabi.StartUpGlobalDaxSrcs(); err.IsNotOk() {
		return
	}
	defer sabi.ShutdownGlobalDaxSrcs()

	sabi.AddGlobalDaxSrc("fuga", NewMapDaxSrc())

	base := sabi.NewDaxBase()

	type MyDax struct {
		sabi.Dax
	}

	dax := MyDax{Dax: base}

	conn, err := sabi.GetDaxConn[MapDaxConn](dax, "hoge")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	conn, err = sabi.GetDaxConn[MapDaxConn](dax, "fuga")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %t\n", err.IsOk())
	fmt.Printf("err.Error() = %s\n", err.Error())

	// Output:
	// conn = sabi_test.MapDaxConn
	// err.IsOk() = true
	// conn = sabi_test.MapDaxConn
	// err.IsOk() = false
	// err.Error() = {reason=DaxSrcIsNotFound, Name=fuga}

	ClearDaxBase()
}

func ExampleDaxBase_SetUpLocalDaxSrc() {
	base := sabi.NewDaxBase()
	defer base.FreeAllLocalDaxSrcs()

	base.SetUpLocalDaxSrc("hoge", NewMapDaxSrc())

	type MyDax struct {
		sabi.Dax
	}

	dax := MyDax{Dax: base}

	conn, err := sabi.GetDaxConn[MapDaxConn](dax, "hoge")
	fmt.Printf("conn = %v\n", reflect.TypeOf(conn))
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())

	// Output:
	// conn = sabi_test.MapDaxConn
	// err.IsOk() = true

	ClearDaxBase()
}

type GettingDax struct {
	sabi.Dax
}

func (dax GettingDax) GetData() (string, sabi.Err) {
	conn, err := sabi.GetDaxConn[MapDaxConn](dax, "hoge")
	if !err.IsOk() {
		return "", err
	}
	data := conn.dataMap["hogehoge"]
	return data, err
}

type SettingDax struct {
	sabi.Dax
}

func (dax SettingDax) SetData(data string) sabi.Err {
	conn, err := sabi.GetDaxConn[MapDaxConn](dax, "fuga")
	if !err.IsOk() {
		return err
	}
	conn.dataMap["fugafuga"] = data
	return err
}

func ExampleDax() {
	// type GettingDax struct {
	//   sabi.Dax
	// }
	// func (dax GettingDax) GetData() (string, sabi.Err) {
	//   conn, err := dax.getDaxConn("hoge")
	//   if !err.IsOk() {
	//     return nil, err
	//   }
	//   return conn.(MapDaxConn).dataMap["hogehoge"], err
	// }
	//
	// type SettingDax struct {
	//   sabi.Dax
	// }
	// func (dax SettingDax) SetData(data string) sabi.Err {
	//   conn, err := dax.getDaxConn("fuga")
	//   if !err.IsOk() {
	//     return nil, err
	//   }
	//   conn.(MapDaxConn).dataMap["fugafuga"] = data
	//   return err
	// }

	hogeDs := NewMapDaxSrc()
	fugaDs := NewMapDaxSrc()

	base := sabi.NewDaxBase()
	defer base.FreeAllLocalDaxSrcs()

	if err := base.SetUpLocalDaxSrc("hoge", hogeDs); err.IsNotOk() {
		return
	}
	if err := base.SetUpLocalDaxSrc("fuga", fugaDs); err.IsNotOk() {
		return
	}

	base = struct {
		sabi.DaxBase
		GettingDax
		SettingDax
	}{
		DaxBase:    base,
		GettingDax: GettingDax{Dax: base},
		SettingDax: SettingDax{Dax: base},
	}

	type DaxForLogic interface {
		GetData() (string, sabi.Err)
		SetData(data string) sabi.Err
	}

	hogeDs.dataMap["hogehoge"] = "hello"
	err := sabi.RunTxn(base, func(dax DaxForLogic) sabi.Err {
		data, err := dax.GetData()
		if !err.IsOk() {
			return err
		}
		err = dax.SetData(data)
		return err
	})
	fmt.Printf("%t\n", err.IsOk())
	fmt.Printf("%s\n", fugaDs.dataMap["fugafuga"])

	// Output:
	// true
	// hello
}
