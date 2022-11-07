# [Sabi][repo-url] [![Go Reference][pkg-dev-img]][pkg-dev-url] [![CI Status][ci-img]][ci-url] [![MIT License][mit-img]][mit-url]

A small framework for Golang applications.

- [What is this?](#what-is-this)
- [Usage](#usage)
- [Supporting Go versions](#support-go-versions)
- [License](#license)

<a name="what-is-this"></a>
## What is this?

Sabi is a small framework to separate logics and data accesses for Golang applications.

A program consists of procedures and data.
And to operate data, procedures includes data accesses, then the rest of procedures except data accesses are logics.
Therefore, a program consists of logics, data accesses and data.

This package is an application framework which explicitly separates procedures into logics and data accesses as layers.
By using this framework, we can remove codes for data accesses from logic parts, and write only specific codes for each data source (e.g. database, messaging services files, and so on)  in data access  parts. 
Moreover, by separating as layers, applications using this framework can change data sources easily by switching data access parts.

<a name="usage"></a>
## Usage

1. [Write logics with a dax interface](#write_logic)
2. [Write dax implementations](#write_dax)
3. [Write mapping of a dax interface and dax implementations](#write_mapping)
4. [Write execution of a procedure](#write_procedure)

<a name="write_logic"></a>
### Write logics with a dax interface

```
  type MyDax interface {
    GetData() (Data, sabi.Err)
    SetData(data Data) sabi.Err
  }

  func MyLogic(dax Dax) sabi.Err {
    data, err := dax.GetData()
    if !err.IsOk() {
      return err
    }
    return dax.SetData(data)
  }
```

<a name="write_dax"></a>
### Write dax implementations

<a name="write_mapping"></a>
### Write mapping of a dax interface and dax implementations

<a name="write_procedure"></a>
### Write execution of a procedure


<a name="support-go-versions"></a>
## Supporting Go versions

This framework supports Go 1.18 or later.

### Actually checked Go versions:

- 1.19.3
- 1.18.8

<a name="license"></a>
## License

Copyright (C) 2022 Takayuki Sato

This program is free software under MIT License.<br>
See the file LICENSE in this distribution for more details.


[repo-url]: https://github.com/sttk-go/sabi
[pkg-dev-img]: https://pkg.go.dev/badge/github.com/sttk-go/sabi.svg
[pkg-dev-url]: https://pkg.go.dev/github.com/sttk-go/sabi
[ci-img]: https://github.com/sttk-go/sabi/actions/workflows/go.yml/badge.svg?branch=main
[ci-url]: https://github.com/sttk-go/sabi/actions
[mit-img]: https://img.shields.io/badge/license-MIT-green.svg
[mit-url]: https://opensource.org/licenses/MIT

