package sabi

type FooXio struct {
	Xio
}

func NewFooXio(xio Xio) FooXio {
	return FooXio{Xio: xio}
}

func (xio FooXio) GetFooConn(name string) (*FooConn, Err) {
	conn, err := xio.GetConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*FooConn), Ok()
}

type BarXio struct {
	Xio
}

func NewBarXio(xio Xio) BarXio {
	return BarXio{Xio: xio}
}

func (xio BarXio) GetBarConn(name string) (*BarConn, Err) {
	conn, err := xio.GetConn(name)
	if !err.IsOk() {
		return nil, err
	}
	return conn.(*BarConn), Ok()
}

/*
func TestXioBase_GetConn_ForEachDataSrc(t *testing.T) {
	Clear()
	defer Clear()

	base := NewXioBase()
	base.AddLocalConnCfg("foo", FooConnCfg{})
	base.AddLocalConnCfg("bar", &BarConnCfg{})
	base.SealLocalConnCfgs()

	fooXio := NewFooXio(base)
	fooConn, err0 := fooXio.GetFooConn("foo")
	assert.True(t, err0.IsOk())
	assert.Equal(t, reflect.TypeOf(fooConn).String(), "*sabi.FooConn")

	barXio := NewBarXio(base)
	barConn, err1 := barXio.GetBarConn("bar")
	assert.True(t, err1.IsOk())
	assert.Equal(t, reflect.TypeOf(barConn).String(), "*sabi.BarConn")
}
*/
