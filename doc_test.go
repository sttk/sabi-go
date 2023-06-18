package sabi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sttk/sabi"
)

type GreetDax interface {
	UserName() (string, sabi.Err)
	Say(greeting string) sabi.Err
}

type ( // possible error reasons
	NoName       struct{}
	FailToOutput struct{ Text string }
)

func GreetLogic(dax GreetDax) sabi.Err {
	name, err := dax.UserName()
	if !err.IsOk() {
		return err
	}
	return dax.Say("Hello, " + name)
}

type mapGreetDax struct {
	m map[string]any
}

func (dax mapGreetDax) UserName() (string, sabi.Err) {
	username, exists := dax.m["username"]
	if !exists {
		return "", sabi.NewErr(NoName{})
	}
	return username.(string), sabi.Ok()
}

func (dax mapGreetDax) Say(greeting string) sabi.Err {
	dax.m["greeting"] = greeting
	return sabi.Ok()
}

func NewMapGreetDaxBase(m map[string]any) sabi.DaxBase {
	base := sabi.NewDaxBase()
	return struct {
		sabi.DaxBase
		mapGreetDax
	}{
		DaxBase:     base,
		mapGreetDax: mapGreetDax{m: m},
	}
}

func TestGreetLogic_unitTest(t *testing.T) {
	m := make(map[string]any)
	base := NewMapGreetDaxBase(m)

	m["username"] = "World"
	err := sabi.RunTxn(base, GreetLogic)
	assert.True(t, err.IsOk())
	assert.Equal(t, m["greeting"], "Hello, World")
}

var osArgs []string

type CliArgsUserDax struct {
}

func (dax CliArgsUserDax) UserName() (string, sabi.Err) {
	if len(osArgs) <= 1 {
		return "", sabi.NewErr(NoName{})
	}
	return osArgs[1], sabi.Ok()
}

type ConsoleOutputDax struct {
}

func (dax ConsoleOutputDax) Say(text string) sabi.Err {
	_, e := fmt.Println(text)
	if e != nil {
		return sabi.NewErr(FailToOutput{Text: text}, e)
	}
	return sabi.Ok()
}

func NewGreetDaxBase() sabi.DaxBase {
	base := sabi.NewDaxBase()
	return struct {
		sabi.DaxBase
		CliArgsUserDax
		ConsoleOutputDax
	}{
		DaxBase:          base,
		CliArgsUserDax:   CliArgsUserDax{},
		ConsoleOutputDax: ConsoleOutputDax{},
	}
}

func TestGreetLogic_executingLogic(t *testing.T) {
	osArgs = []string{"cmd", "Tom"}
	base := NewGreetDaxBase()
	err := sabi.RunTxn(base, GreetLogic)
	assert.True(t, err.IsOk())
}

type PrintDax interface {
	Print() sabi.Err
}

type ConsoleOutputDax2 struct {
	text string
}

func (dax *ConsoleOutputDax2) Say(text string) sabi.Err {
	dax.text = text
	return sabi.Ok()
}
func (dax *ConsoleOutputDax2) Print() sabi.Err {
	_, e := fmt.Println(dax.text)
	if e != nil {
		return sabi.NewErr(FailToOutput{Text: dax.text}, e)
	}
	return sabi.Ok()
}

func NewGreetDaxBase2() sabi.DaxBase {
	base := sabi.NewDaxBase()
	return struct {
		sabi.DaxBase
		CliArgsUserDax
		*ConsoleOutputDax2
	}{
		DaxBase:           base,
		CliArgsUserDax:    CliArgsUserDax{},
		ConsoleOutputDax2: &ConsoleOutputDax2{},
	}
}

func TestGreetLogic_MovingOutputs(t *testing.T) {
	base := NewGreetDaxBase2()
	err := sabi.RunTxn(base, GreetLogic)
	assert.True(t, err.IsOk())
	err = sabi.RunTxn(base, func(dax PrintDax) sabi.Err {
		return dax.Print()
	})
	assert.True(t, err.IsOk())
}

func TestGreetLogic_MovingOutputs2(t *testing.T) {
	base := NewGreetDaxBase2()
	txn0 := sabi.Txn(base, GreetLogic)
	txn1 := sabi.Txn(base, func(dax PrintDax) sabi.Err {
		return dax.Print()
	})
	err := sabi.RunSeq(txn0, txn1)
	assert.True(t, err.IsOk())
}
