package sabi_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"testing"
)

//// MapDaxSrc

type MapDaxSrc struct {
	dataMap map[string]string
}

func NewMapDaxSrc() MapDaxSrc {
	return MapDaxSrc{dataMap: make(map[string]string)}
}

func (ds MapDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return &MapDaxConn{dataMap: ds.dataMap}, sabi.Ok()
}

func (ds MapDaxSrc) SetUp() sabi.Err {
	return sabi.Ok()
}

func (ds MapDaxSrc) End() {
}

//// MapDaxConn

type MapDaxConn struct {
	dataMap map[string]string
}

func (conn *MapDaxConn) Commit() sabi.Err {
	return sabi.Ok()
}

func (conn *MapDaxConn) Rollback() {
}

func (conn *MapDaxConn) Close() {
}

//// MapDax

type MapDax struct {
	sabi.Dax
}

func NewMapDax(base sabi.DaxBase) MapDax {
	return MapDax{Dax: base}
}

func (dax MapDax) GetMapDaxConn(name string) (*MapDaxConn, sabi.Err) {
	conn, err := dax.GetDaxConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*MapDaxConn), err
}

//// HogeFugaDax & HogeFugaLogic

type HogeFugaDax interface {
	GetHogeData() (string, sabi.Err)
	SetFugaData(data string) sabi.Err
}

func HogeFugaLogic(dax HogeFugaDax) sabi.Err {
	data, err := dax.GetHogeData()
	if !err.IsOk() {
		return err
	}
	err = dax.SetFugaData(data)
	return err
}

//// FugaPiyoDax & FugaPiyoLogic

type FugaPiyoDax interface {
	GetFugaData() (string, sabi.Err)
	SetPiyoData(data string) sabi.Err
}

type PiyoHogeraDax interface {
	GetPiyoData() (string, sabi.Err)
	SetHogeraData(data string) sabi.Err
}

func FugaPiyoLogic(dax FugaPiyoDax) sabi.Err {
	data, err := dax.GetFugaData()
	if !err.IsOk() {
		return err
	}
	err = dax.SetPiyoData(data)
	return err
}

//// HogeDax

type HogeDax struct {
	MapDax
}

func (dax HogeDax) GetHogeData() (string, sabi.Err) {
	conn, err := dax.GetMapDaxConn("hoge")
	if !err.IsOk() {
		return "", err
	}
	data := conn.dataMap["hogehoge"]
	return data, err
}

//// FugaDax

type FugaDax struct {
	MapDax
}

func (dax FugaDax) SetFugaData(data string) sabi.Err {
	conn, err := dax.GetMapDaxConn("fuga")
	if !err.IsOk() {
		return err
	}
	conn.dataMap["fugafuga"] = data
	return err
}

func (dax FugaDax) GetFugaData() (string, sabi.Err) {
	conn, err := dax.GetMapDaxConn("fuga")
	if !err.IsOk() {
		return "", err
	}
	data := conn.dataMap["fugafuga"]
	return data, err
}

//// PiyoDax

type PiyoDax struct {
	MapDax
}

func (dax PiyoDax) SetPiyoData(data string) sabi.Err {
	conn, err := dax.GetMapDaxConn("piyo")
	if !err.IsOk() {
		return err
	}
	conn.dataMap["piyopiyo"] = data
	return err
}

func (dax PiyoDax) GetPiyoData() (string, sabi.Err) {
	conn, err := dax.GetMapDaxConn("piyo")
	if !err.IsOk() {
		return "", err
	}
	data := conn.dataMap["piyopiyo"]
	return data, err
}

//// HogeFugaPiyoDaxBase

func NewHogeFugaPiyoDaxBase() sabi.DaxBase {
	base := sabi.NewDaxBase()

	return struct {
		sabi.DaxBase
		HogeDax
		FugaDax
		PiyoDax
	}{
		DaxBase: base,
		HogeDax: HogeDax{MapDax: NewMapDax(base)},
		FugaDax: FugaDax{MapDax: NewMapDax(base)},
		PiyoDax: PiyoDax{MapDax: NewMapDax(base)},
	}
}

//// Test cases

func TestDax_runTxn(t *testing.T) {
	sabi.ClearDaxBase()
	defer sabi.ClearDaxBase()

	hogeDs := NewMapDaxSrc()
	fugaDs := NewMapDaxSrc()
	piyoDs := NewMapDaxSrc()

	base := NewHogeFugaPiyoDaxBase()
	base.SetUpLocalDaxSrc("hoge", hogeDs)
	base.SetUpLocalDaxSrc("fuga", fugaDs)
	base.SetUpLocalDaxSrc("piyo", piyoDs)

	hogeDs.dataMap["hogehoge"] = "Hello, world"

	err := sabi.RunTxn(base, HogeFugaLogic)
	assert.True(t, err.IsOk())
	err = sabi.RunTxn(base, FugaPiyoLogic)
	assert.True(t, err.IsOk())

	assert.Equal(t, piyoDs.dataMap["piyopiyo"], "Hello, world")

}
