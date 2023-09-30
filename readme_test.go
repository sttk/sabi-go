package sabi_test

import (
	"fmt"
	"os"
	//"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sttk/sabi"
	"github.com/sttk/sabi/errs"
)

var origOsArgs []string = os.Args

func reset() {
	os.Args = origOsArgs
}

type ( // possible error reasons
	NoName        struct{}
	FailToGetHour struct{}
	FailToOutput  struct{ Text string }
)

type GreetDax interface {
	UserName() (string, errs.Err)
	Hour() (int, errs.Err)
	Output(text string) errs.Err
}

func GreetLogic(dax GreetDax) errs.Err {
	hour, err := dax.Hour()
	if err.IsNotOk() {
		return err
	}

	var s string
	switch {
	case 5 <= hour && hour < 12:
		s = "Good morning, "
	case 12 <= hour && hour < 16:
		s = "Good afternoon, "
	case 16 <= hour && hour < 21:
		s = "Good evening, "
	default:
		s = "Hi, "
	}

	err = dax.Output(s)
	if err.IsNotOk() {
		return err
	}

	name, err := dax.UserName()
	if err.IsNotOk() {
		return err
	}

	return dax.Output(name + ".\n")
}

type MapGreetDax struct {
	sabi.Dax
	m map[string]any
}

func (dax MapGreetDax) UserName() (string, errs.Err) {
	name, exists := dax.m["username"]
	if !exists {
		return "", errs.New(NoName{})
	}
	return name.(string), errs.Ok()
}

func (dax MapGreetDax) Hour() (int, errs.Err) {
	hour, exists := dax.m["hour"]
	if !exists {
		return 0, errs.New(FailToGetHour{})
	}
	return hour.(int), errs.Ok()
}

func (dax MapGreetDax) Output(text string) errs.Err {
	var s string
	v, exists := dax.m["greeting"]
	if exists {
		s = v.(string)
	}
	if s == "error" { // for testings the error case
		return errs.New(FailToOutput{Text: text})
	}
	dax.m["greeting"] = s + text
	return errs.Ok()
}

func NewMapGreetDaxBase(m map[string]any) sabi.DaxBase {
	base := sabi.NewDaxBase()
	return struct {
		sabi.DaxBase
		MapGreetDax
	}{
		DaxBase:     base,
		MapGreetDax: MapGreetDax{m: m},
	}
}

func TestGreetLogic_morning(t *testing.T) {
	m := make(map[string]any)
	base := NewMapGreetDaxBase(m)

	m["username"] = "everyone"
	m["hour"] = 10

	err := sabi.Txn(base, GreetLogic)
	if err.IsNotOk() {
		t.Errorf(err.Error())
	}

	if m["greeting"] != "Good morning, everyone.\n" {
		t.Errorf("Bad greeting: %v\n", m["greeting"])
	}
}

type CliArgsDax struct {
	sabi.Dax
}

func (dax CliArgsDax) UserName() (string, errs.Err) {
	if len(os.Args) <= 1 {
		return "", errs.New(NoName{})
	}
	return os.Args[1], errs.Ok()
}

//func (dax CliArgsDax) Hour() (int, errs.Err) {
//	if len(os.Args) <= 2 {
//		return 0, errs.New(FailToGetHour{})
//	}
//	n, err := strconv.Atoi(os.Args[2])
//	if err != nil {
//		return 0, errs.New(FailToGetHour{}, err)
//	}
//	return n, errs.Ok()
//}

type ConsoleDax struct {
	sabi.Dax
}

// func (dax ConsoleDax) Output(text string) errs.Err {
func (dax ConsoleDax) Print(text string) errs.Err {
	fmt.Print(text)
	return errs.Ok()
}

func NewGreetDaxBase() sabi.DaxBase {
	base := sabi.NewDaxBase()
	return struct {
		sabi.DaxBase
		CliArgsDax
		SystemClockDax
		MemoryDax
		ConsoleDax
	}{
		DaxBase:        base,
		CliArgsDax:     CliArgsDax{Dax: base},
		SystemClockDax: SystemClockDax{Dax: base},
		MemoryDax:      MemoryDax{Dax: base},
		ConsoleDax:     ConsoleDax{Dax: base},
	}
}

func app() errs.Err {
	base := NewGreetDaxBase()
	defer base.Close()

	return base.Uses("memory", &MemoryDaxSrc{}).
		IfOk(sabi.Txn_(base, GreetLogic)).
		IfOk(sabi.Txn_(base, PrintLogic))

	//return base.Uses("memory", MemoryDaxSrc{}).
	//	IfOk(sabi.Txn_(base, GreetLogic))

	//return sabi.Txn(base, GreetLogic)

}

func TestGreetApp_main(t *testing.T) {
	defer reset()

	os.Args = []string{"cmd", "foo", "10"}

	if err := sabi.StartApp(app); err.IsNotOk() {
		//fmt.Println(err.Error())
		//os.Exit(1);
		t.Errorf(err.Error())
	}
}

type SystemClockDax struct {
	sabi.Dax
}

func (dax SystemClockDax) Hour() (int, errs.Err) {
	return time.Now().Hour(), errs.Ok()
}

type MemoryDaxSrc struct {
	buf strings.Builder
}

func (ds *MemoryDaxSrc) Setup(ag sabi.AsyncGroup) errs.Err {
	return errs.Ok()
}

func (ds *MemoryDaxSrc) Close() {
	ds.buf.Reset()
}

func (ds *MemoryDaxSrc) CreateDaxConn() (sabi.DaxConn, errs.Err) {
	return MemoryDaxConn{buf: &(ds.buf)}, errs.Ok()
}

type MemoryDaxConn struct {
	buf *strings.Builder
}

func (conn MemoryDaxConn) Append(text string) {
	conn.buf.WriteString(text)
}

func (conn MemoryDaxConn) Get() string {
	return conn.buf.String()
}

func (conn MemoryDaxConn) Commit(ag sabi.AsyncGroup) errs.Err {
	return errs.Ok()
}

func (conn MemoryDaxConn) IsCommitted() bool {
	return true
}

func (conn MemoryDaxConn) Rollback(ag sabi.AsyncGroup) {
}

func (conn MemoryDaxConn) ForceBack(ag sabi.AsyncGroup) {
}

func (conn MemoryDaxConn) Close() {
}

type MemoryDax struct {
	sabi.Dax
}

func (dax MemoryDax) Output(text string) errs.Err {
	conn, err := sabi.GetDaxConn[MemoryDaxConn](dax, "memory")
	if err.IsNotOk() {
		return err
	}
	conn.Append(text)
	return err
}

func (dax MemoryDax) GetText() (string, errs.Err) {
	conn, err := sabi.GetDaxConn[MemoryDaxConn](dax, "memory")
	if err.IsNotOk() {
		return "", err
	}
	return conn.Get(), err
}

type PrintDax interface {
	GetText() (string, errs.Err)
	Print(text string) errs.Err
}

func PrintLogic(dax PrintDax) errs.Err {
	text, err := dax.GetText()
	if err.IsNotOk() {
		return err
	}
	return dax.Print(text)
}
