package sabi_test

import (
	"os"

	"github.com/sttk/sabi"
	"github.com/sttk/sabi/errs"
)

///

type CliArgDaxSrc struct{}

func (ds CliArgDaxSrc) Setup(ag sabi.AsyncGroup) errs.Err {
	return errs.Ok()
}
func (ds CliArgDaxSrc) Close() {}
func (ds CliArgDaxSrc) CreateDaxConn() (sabi.DaxConn, errs.Err) {
	return nil, errs.Ok()
}

func NewCliArgDaxSrc(osArgs []string) CliArgDaxSrc {
	return CliArgDaxSrc{}
}

type DatabaseDaxSrc struct{}

func (ds DatabaseDaxSrc) Setup(ag sabi.AsyncGroup) errs.Err {
	return errs.Ok()
}
func (ds DatabaseDaxSrc) Close() {}
func (ds DatabaseDaxSrc) CreateDaxConn() (sabi.DaxConn, errs.Err) {
	return nil, errs.Ok()
}

func NewDatabaseDaxSrc(driverName, dataSourceName string) DatabaseDaxSrc {
	return DatabaseDaxSrc{}
}

var (
	driverName     string
	dataSourceName string
)

type HttpRequestDaxSrc struct{}

func (ds HttpRequestDaxSrc) Setup(ag sabi.AsyncGroup) errs.Err {
	return errs.Ok()
}
func (ds HttpRequestDaxSrc) Close() {}
func (ds HttpRequestDaxSrc) CreateDaxConn() (sabi.DaxConn, errs.Err) {
	return nil, errs.Ok()
}

func NewHttpRequestDaxSrc(req any) HttpRequestDaxSrc {
	return HttpRequestDaxSrc{}
}

var req any

type HttpResponseDaxSrc struct{}

func (ds HttpResponseDaxSrc) Setup(ag sabi.AsyncGroup) errs.Err {
	return errs.Ok()
}
func (ds HttpResponseDaxSrc) Close() {}
func (ds HttpResponseDaxSrc) CreateDaxConn() (sabi.DaxConn, errs.Err) {
	return nil, errs.Ok()
}

func NewHttpResponseDaxSrc(resp any) HttpResponseDaxSrc {
	return HttpResponseDaxSrc{}
}

var resp any

type CliArgOptionDax struct {
	sabi.Dax
}

type DatabaseSetDax struct {
	sabi.Dax
}

type HttpReqParamDax struct {
	sabi.Dax
}

type HttpRespOutputDax struct {
	sabi.Dax
}

///

func ExampleClose() {
	sabi.Uses("cliargs", NewCliArgDaxSrc(os.Args))

	err := sabi.Setup()
	if err.IsNotOk() {
		// ...
	}
	defer sabi.Close()

	// ...
}

func ExampleSetup() {
	sabi.Uses("cliargs", NewCliArgDaxSrc(os.Args))
	sabi.Uses("database", NewDatabaseDaxSrc(driverName, dataSourceName))

	err := sabi.Setup()
	if err.IsNotOk() {
		// ...
	}
	defer sabi.Close()

	// ...
}

func ExampleStartApp() {
	sabi.Uses("cliargs", NewCliArgDaxSrc(os.Args))
	sabi.Uses("database", NewDatabaseDaxSrc(driverName, dataSourceName))

	app := func() errs.Err {
		// ...
		return errs.Ok()
	}

	err := sabi.StartApp(app)
	if err.IsNotOk() {
		// ...
	}
}

func ExampleUses() {
	sabi.Uses("cliargs", NewCliArgDaxSrc(os.Args))
	sabi.Uses("database", NewDatabaseDaxSrc(driverName, dataSourceName))

	err := sabi.StartApp(func() errs.Err {
		// ...
		return errs.Ok()
	})
	if err.IsNotOk() {
		// ...
	}
	// ...
}

func ExampleDaxBase() {
	sabi.Uses("cliargs", NewCliArgDaxSrc(os.Args))
	sabi.Uses("database", NewDatabaseDaxSrc(driverName, dataSourceName))

	err := sabi.Setup()
	if err.IsNotOk() {
		// ...
	}
	defer sabi.Close()

	NewMyBase := func() sabi.DaxBase {
		base := sabi.NewDaxBase()
		return &struct {
			sabi.DaxBase
			CliArgOptionDax
			DatabaseSetDax
			HttpReqParamDax
			HttpRespOutputDax
		}{
			DaxBase:           base,
			CliArgOptionDax:   CliArgOptionDax{Dax: base},
			DatabaseSetDax:    DatabaseSetDax{Dax: base},
			HttpReqParamDax:   HttpReqParamDax{Dax: base},
			HttpRespOutputDax: HttpRespOutputDax{Dax: base},
		}
	}

	type GetSetDax struct {
		sabi.Dax
		// ...
	}

	GetSetLogic := func(dax GetSetDax) errs.Err {
		// ...
		return errs.Ok()
	}

	type OutputDax struct {
		sabi.Dax
		// ...
	}

	OutputLogic := func(dax OutputDax) errs.Err {
		// ...
		return errs.Ok()
	}

	base := NewMyBase()
	defer base.Close()

	err = base.Uses("httpReq", NewHttpRequestDaxSrc(req)).
		IfOk(sabi.Txn_(base, GetSetLogic)).
		IfOk(base.Disuses_("httpReq")).
		IfOk(base.Uses_("httpResp", NewHttpResponseDaxSrc(resp))).
		IfOk(sabi.Txn_(base, OutputLogic))
	if err.IsNotOk() {
		// ...
	}
}

func ExampleNewDaxBase() {
	base0 := sabi.NewDaxBase()

	base := &struct {
		sabi.DaxBase
		CliArgOptionDax
		DatabaseSetDax
		HttpReqParamDax
		HttpRespOutputDax
	}{
		DaxBase:           base0,
		CliArgOptionDax:   CliArgOptionDax{Dax: base0},
		DatabaseSetDax:    DatabaseSetDax{Dax: base0},
		HttpReqParamDax:   HttpReqParamDax{Dax: base0},
		HttpRespOutputDax: HttpRespOutputDax{Dax: base0},
	}
	// Output:
	_ = base
}

func NewMyBase() sabi.DaxBase {
	base := sabi.NewDaxBase()
	return &struct {
		sabi.DaxBase
		CliArgOptionDax
		DatabaseSetDax
		HttpReqParamDax
		HttpRespOutputDax
	}{
		DaxBase:           base,
		CliArgOptionDax:   CliArgOptionDax{Dax: base},
		DatabaseSetDax:    DatabaseSetDax{Dax: base},
		HttpReqParamDax:   HttpReqParamDax{Dax: base},
		HttpRespOutputDax: HttpRespOutputDax{Dax: base},
	}
}

func ExampleTxn() {
	base := NewMyBase()
	defer base.Close()

	type GetSetDax struct {
		sabi.Dax
		// ...
	}

	GetSetLogic := func(dax GetSetDax) errs.Err {
		// ...
		return errs.Ok()
	}

	err := sabi.Txn(base, GetSetLogic)
	if err.IsNotOk() {
		// ...
	}
}
