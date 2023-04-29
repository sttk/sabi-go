package sabi_test

import (
	"container/list"
	"github.com/sttk-go/sabi"
)

var (
	Logs list.List

	WillFailToSetUpFooDaxSrc   bool
	WillFailToCommitFooDaxConn bool
	WillFailToCreateFooDaxConn bool

	WillFailToCreateBDaxConn = false
	WillFailToCommitBDaxConn = false
)

func ClearDaxBase() {
	sabi.ClearGlobalDaxSrcs()

	Logs.Init()

	WillFailToSetUpFooDaxSrc = false
	WillFailToCommitFooDaxConn = false
	WillFailToCreateFooDaxConn = false

	WillFailToCreateBDaxConn = false
	WillFailToCommitBDaxConn = false
}

type /* error reasons */ (
	FailToDoSomething struct{ Text string }

	FailToCreateBDaxConn struct{}
	FailToCommitBDaxConn struct{}
	FailToRunLogic       struct{}
)

// FooDaxConn

type FooDaxConn struct {
	Label string
	Map   map[string]string
}

func (conn FooDaxConn) Commit() sabi.Err {
	if WillFailToCommitFooDaxConn {
		return sabi.NewErr(FailToDoSomething{Text: "FailToCommitFooDaxConn"})
	}
	Logs.PushBack("FooDaxConn#Commit")
	return sabi.Ok()
}

func (conn FooDaxConn) Rollback() {
	Logs.PushBack("FooDaxConn#Rollback")
}

func (conn FooDaxConn) Close() {
	Logs.PushBack("FooDaxConn#Close")
}

// FooDaxSrc

type FooDaxSrc struct {
	Label string
}

func (ds FooDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	if WillFailToCreateFooDaxConn {
		return nil, sabi.NewErr(FailToDoSomething{Text: "FailToCreateFooDaxConn"})
	}
	Logs.PushBack("FooDaxSrc#CreateDaxConn")
	return FooDaxConn{Label: ds.Label, Map: make(map[string]string)}, sabi.Ok()
}

func (ds FooDaxSrc) SetUp() sabi.Err {
	if WillFailToSetUpFooDaxSrc {
		return sabi.NewErr(FailToDoSomething{Text: "FailToSetUpFooDaxSrc"})
	}
	Logs.PushBack("FooDaxSrc#SetUp")
	return sabi.Ok()
}

func (ds FooDaxSrc) End() {
	Logs.PushBack("FooDaxSrc#End")
}

// FooDax

type FooDax struct {
	sabi.Dax
}

func (dax FooDax) GetFooData() (string, sabi.Err) {
	conn, err := dax.GetDaxConn("foo")
	if err.IsNotOk() {
		return "", err
	}
	return conn.(FooDaxConn).Map["data"], err
}

// BarDaxConn

type BarDaxConn struct {
	Label string
	Map   map[string]string
}

func (conn BarDaxConn) Commit() sabi.Err {
	Logs.PushBack("BarDaxConn#Commit")
	return sabi.Ok()
}

func (conn BarDaxConn) Rollback() {
	Logs.PushBack("BarDaxConn#Rollback")
}

func (conn BarDaxConn) Close() {
	Logs.PushBack("BarDaxConn#Close")
}

// BarDaxSrc

type BarDaxSrc struct {
	Label string
}

func (ds BarDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	Logs.PushBack("BarDaxSrc#CreateDaxConn")
	return BarDaxConn{Label: ds.Label, Map: make(map[string]string)}, sabi.Ok()
}

func (ds BarDaxSrc) SetUp() sabi.Err {
	Logs.PushBack("BarDaxSrc#SetUp")
	return sabi.Ok()
}

func (ds BarDaxSrc) End() {
	Logs.PushBack("BarDaxSrc#End")
}

// MapDaxSrc

type MapDaxSrc struct {
	dataMap map[string]string
}

func NewMapDaxSrc() MapDaxSrc {
	return MapDaxSrc{dataMap: make(map[string]string)}
}

func (ds MapDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return MapDaxConn{dataMap: ds.dataMap}, sabi.Ok()
}

func (ds MapDaxSrc) SetUp() sabi.Err {
	return sabi.Ok()
}

func (ds MapDaxSrc) End() {
}

// MapDaxConn

type MapDaxConn struct {
	dataMap map[string]string
}

func (conn MapDaxConn) Commit() sabi.Err {
	return sabi.Ok()
}

func (conn MapDaxConn) Rollback() {
}

func (conn MapDaxConn) Close() {
}

// MapDax

type MapDax struct {
	sabi.Dax
}

func NewMapDax(base sabi.DaxBase) MapDax {
	return MapDax{Dax: base}
}

func (dax MapDax) GetMapDaxConn(name string) (MapDaxConn, sabi.Err) {
	conn, err := dax.GetDaxConn(name)
	if err.IsNotOk() {
		return MapDaxConn{}, err
	}
	return conn.(MapDaxConn), err
}

// HogeFugaDax

type HogeFugaDax interface {
	GetHogeData() (string, sabi.Err)
	SetFugaData(data string) sabi.Err
}

// HogeFugaLogic

func HogeFugaLogic(dax HogeFugaDax) sabi.Err {
	data, err := dax.GetHogeData()
	if err.IsNotOk() {
		return err
	}
	err = dax.SetFugaData(data)
	return err
}

// FugaPiyoDax

type FugaPiyoDax interface {
	GetFugaData() (string, sabi.Err)
	SetPiyoData(data string) sabi.Err
}

// FugaPiyoLogic

func FugaPiyoLogic(dax FugaPiyoDax) sabi.Err {
	data, err := dax.GetFugaData()
	if err.IsNotOk() {
		return err
	}
	err = dax.SetPiyoData(data)
	return err
}

// HogeDax

type HogeDax struct {
	MapDax
}

func (dax HogeDax) GetHogeData() (string, sabi.Err) {
	conn, err := dax.GetMapDaxConn("hoge")
	if err.IsNotOk() {
		return "", err
	}
	data := conn.dataMap["hogehoge"]
	return data, err
}

func (dax HogeDax) SetHogeData(data string) sabi.Err {
	conn, err := dax.GetMapDaxConn("hoge")
	if err.IsNotOk() {
		return err
	}
	conn.dataMap["hogehoge"] = data
	return err
}

// FugaDax

type FugaDax struct {
	MapDax
}

func (dax FugaDax) GetFugaData() (string, sabi.Err) {
	conn, err := dax.GetMapDaxConn("fuga")
	if err.IsNotOk() {
		return "", err
	}
	data := conn.dataMap["fugafuga"]
	return data, err
}

func (dax FugaDax) SetFugaData(data string) sabi.Err {
	conn, err := dax.GetMapDaxConn("fuga")
	if err.IsNotOk() {
		return err
	}
	conn.dataMap["fugafuga"] = data
	return err
}

// PiyoDax

type PiyoDax struct {
	MapDax
}

func (dax PiyoDax) GetPiyoData() (string, sabi.Err) {
	conn, err := dax.GetMapDaxConn("piyo")
	if err.IsNotOk() {
		return "", err
	}
	data := conn.dataMap["piyopiyo"]
	return data, err
}

func (dax PiyoDax) SetPiyoData(data string) sabi.Err {
	conn, err := dax.GetMapDaxConn("piyo")
	if err.IsNotOk() {
		return err
	}
	conn.dataMap["piyopiyo"] = data
	return err
}

// HogeFugaPiyoDaxBase

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

// ADaxSrc

type ADaxSrc struct {
	AMap map[string]string
}

func NewADaxSrc() ADaxSrc {
	return ADaxSrc{AMap: make(map[string]string)}
}

func (ds ADaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return ADaxConn{AMap: ds.AMap}, sabi.Ok()
}

func (ds ADaxSrc) SetUp() sabi.Err {
	return sabi.Ok()
}

func (ds ADaxSrc) End() {
}

// ADaxConn

type ADaxConn struct {
	AMap map[string]string
}

func (conn ADaxConn) Commit() sabi.Err {
	Logs.PushBack("ADaxConn#Commit")
	return sabi.Ok()
}

func (conn ADaxConn) Rollback() {
	Logs.PushBack("ADaxConn#Rollback")
}

func (conn ADaxConn) Close() {
	Logs.PushBack("ADaxConn#Close")
}

// ADax

type ADax struct {
	sabi.Dax
}

func (dax ADax) GetADaxConn(name string) (ADaxConn, sabi.Err) {
	conn, err := dax.GetDaxConn(name)
	if err.IsNotOk() {
		return ADaxConn{}, err
	}
	return conn.(ADaxConn), err
}

// BDaxSrc

type BDaxSrc struct {
	BMap map[string]string
}

func NewBDaxSrc() BDaxSrc {
	return BDaxSrc{BMap: make(map[string]string)}
}

func (ds BDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	if WillFailToCreateBDaxConn {
		return nil, sabi.NewErr(FailToCreateBDaxConn{})
	}
	return BDaxConn{BMap: ds.BMap}, sabi.Ok()
}

func (ds BDaxSrc) SetUp() sabi.Err {
	return sabi.Ok()
}

func (ds BDaxSrc) End() {
}

// BDaxConn

type BDaxConn struct {
	BMap map[string]string
}

func (conn BDaxConn) Commit() sabi.Err {
	if WillFailToCommitBDaxConn {
		return sabi.NewErr(FailToCommitBDaxConn{})
	}
	Logs.PushBack("BDaxConn#Commit")
	return sabi.Ok()
}

func (conn BDaxConn) Rollback() {
	Logs.PushBack("BDaxConn#Rollback")
}

func (conn BDaxConn) Close() {
	Logs.PushBack("BDaxConn#Close")
}

// BDax

type BDax struct {
	sabi.Dax
}

func (dax BDax) GetBDaxConn(name string) (BDaxConn, sabi.Err) {
	conn, err := dax.GetDaxConn(name)
	if err.IsNotOk() {
		return BDaxConn{}, err
	}
	return conn.(BDaxConn), err
}

// CDaxSrc

type CDaxSrc struct {
	CMap map[string]string
}

func NewCDaxSrc() CDaxSrc {
	return CDaxSrc{CMap: make(map[string]string)}
}

func (ds CDaxSrc) CreateDaxConn() (sabi.DaxConn, sabi.Err) {
	return CDaxConn{CMap: ds.CMap}, sabi.Ok()
}

func (ds CDaxSrc) SetUp() sabi.Err {
	return sabi.Ok()
}

func (ds CDaxSrc) End() {
}

// CDaxConn

type CDaxConn struct {
	CMap map[string]string
}

func (conn CDaxConn) Commit() sabi.Err {
	Logs.PushBack("CDaxConn#Commit")
	return sabi.Ok()
}

func (conn CDaxConn) Rollback() {
	Logs.PushBack("CDaxConn#Rollback")
}

func (conn CDaxConn) Close() {
	Logs.PushBack("CDaxConn#Close")
}

// CDax

type CDax struct {
	sabi.Dax
}

func (dax CDax) GetCDaxConn(name string) (CDaxConn, sabi.Err) {
	conn, err := dax.GetDaxConn(name)
	if err.IsNotOk() {
		return CDaxConn{}, err
	}
	return conn.(CDaxConn), err
}

// AGetDax

type AGetDax struct {
	ADax
}

func (dax AGetDax) GetAData() (string, sabi.Err) {
	conn, err := dax.GetADaxConn("aaa")
	if !err.IsOk() {
		return "", err
	}
	data := conn.AMap["a"]
	return data, sabi.Ok()
}

// BGetSetDax

type BGetSetDax struct {
	BDax
}

func (dax BGetSetDax) GetBData() (string, sabi.Err) {
	conn, err := dax.GetBDaxConn("bbb")
	if !err.IsOk() {
		return "", err
	}
	data := conn.BMap["b"]
	return data, sabi.Ok()
}

func (dax BGetSetDax) SetBData(data string) sabi.Err {
	conn, err := dax.GetBDaxConn("bbb")
	if !err.IsOk() {
		return err
	}
	conn.BMap["b"] = data
	return sabi.Ok()
}

// CSetDax

type CSetDax struct {
	CDax
}

func (dax CSetDax) SetCData(data string) sabi.Err {
	conn, err := dax.GetCDaxConn("ccc")
	if !err.IsOk() {
		return err
	}
	conn.CMap["c"] = data
	return sabi.Ok()
}
