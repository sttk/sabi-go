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

func (dax FooDax) GetFooDaxConn(name string) (*FooDaxConn, Err) {
	conn, err := dax.GetDaxConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*FooDaxConn), Ok()
}

type BarDax struct {
	Dax
}

func NewBarDax(dax Dax) BarDax {
	return BarDax{Dax: dax}
}

func (dax BarDax) GetBarDaxConn(name string) (*BarDaxConn, Err) {
	conn, err := dax.GetDaxConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*BarDaxConn), Ok()
}

func TestDax_GetXxxConn(t *testing.T) {
	Clear()
	defer Clear()

	base := NewDaxBase()
	base.AddLocalDaxSrc("foo", FooDaxSrc{})
	base.AddLocalDaxSrc("bar", &BarDaxSrc{})

	base.begin()

	fooDax := NewFooDax(base)
	fooConn, fooErr := fooDax.GetFooDaxConn("foo")
	assert.True(t, fooErr.IsOk())
	assert.Equal(t, reflect.TypeOf(fooConn).String(), "*sabi.FooDaxConn")

	barDax := NewBarDax(base)
	barConn, barErr := barDax.GetBarDaxConn("bar")
	assert.True(t, barErr.IsOk())
	assert.Equal(t, reflect.TypeOf(barConn).String(), "*sabi.BarDaxConn")
}
