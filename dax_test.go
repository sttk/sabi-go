package sabi

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type FooDax struct {
	Dax
}

func NewFooDax(dax Dax) FooDax {
	return FooDax{Dax: dax}
}

func (dax FooDax) GetFooConn(name string) (*FooConn, Err) {
	conn, err := dax.GetConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*FooConn), Ok()
}

type BarDax struct {
	Dax
}

func NewBarDax(dax Dax) BarDax {
	return BarDax{Dax: dax}
}

func (dax BarDax) GetBarConn(name string) (*BarConn, Err) {
	conn, err := dax.GetConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*BarConn), Ok()
}

func TestDax_GetXxxConn(t *testing.T) {
	Clear()
	defer Clear()

	base := NewConnBase()
	base.addLocalConnCfg("foo", FooConnCfg{})
	base.addLocalConnCfg("bar", &BarConnCfg{})
	base.begin()

	fooDax := NewFooDax(base)
	fooConn, fooErr := fooDax.GetFooConn("foo")
	assert.True(t, fooErr.IsOk())
	assert.Equal(t, reflect.TypeOf(fooConn).String(), "*sabi.FooConn")

	barDax := NewBarDax(base)
	barConn, barErr := barDax.GetBarConn("bar")
	assert.True(t, barErr.IsOk())
	assert.Equal(t, reflect.TypeOf(barConn).String(), "*sabi.BarConn")
}
