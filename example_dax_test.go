package sabi_test

import (
	"github.com/sttk-go/sabi"
)

func ExampleDax() {
	base := sabi.NewDaxBase()

	type MyDax interface {
		GetData() string
		SetData(data string)
	}

	dax := struct {
		FooGetterDax
		BarSetterDax
	}{
		FooGetterDax: FooGetterDax{Dax: base},
		BarSetterDax: BarSetterDax{Dax: base},
	}

	proc := sabi.NewProc[MyDax](base, dax)
	proc.RunTxn(func(dax MyDax) sabi.Err {
		data := dax.GetData()
		dax.SetData(data)
		return sabi.Ok()
	})

	// Output:

	sabi.Clear()
}
