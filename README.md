# [Sabi][repo-url] [![Go Reference][pkg-dev-img]][pkg-dev-url] [![CI Status][ci-img]][ci-url] [![MIT License][mit-img]][mit-url]

A small framework to separate logics and data accesses for Golang application.

## Concept

Sabi is a small framework for Golang applications.
This framework separates an application to logic parts and data access parts, and enables to implement each of them independently, then to combine them.

### Separation of logics and data accesses

In general, a program consists of procedures and data.
And procedures include data accesses for operating data, and the rest of procedures are logics.
So we can say that a program consists of logics, data accesses and data.

Furthermore, we often think to separate an application to multiple layers, for example, controller layer, application logic layer, and data access layer.
The logic and data access mentioned in this framework are partially matched those layers, but are not matched in another part.
For example, in the controller layer, there are input data and output data. (In a web application there are request data and response data, and in a command line application there are console input and output.)
Even though all logical processes are moved into the application logic layer, it is remained to transform input data of the controller layer into input data of the application logic layer, and to transform output data of the application logic layer into the output data of the controller layer.
The data accesses mentioned in this framework also includes those data accesses.

### Changes composition of data access methods by concerns

Dax is a collection of data access methods. These methods will be collected/divided by data source from an implementation perspective. On the other hand, they will be collected/divided by logic from a usage perspective

In general programming, a developer chooses the necessary methods for their logic from among all available methods. And after programming, those methods will be buried in the program code of the logic, and it will become unclear which methods are used without tracing the logic.

In applications using the Sabi framework, a logic is implemented as a function that takes only one argument, a dax interface. And this interface can define only the methods required by the logic.
Therefore, a dax interface can make clear which methods are used in a logic. And also, a dax interface can constraint methods available for a logic.


## Usage

The usage of this framework is described on the overview in the go package document.

See https://pkg.go.dev/github.com/sttk-go/sabi#pkg-overview.


## Supporting Go versions

This framework supports Go 1.18 or later.

### Actual test results for each Go version:

```
% gvm-fav
Now using version go1.18.10
go version go1.18.10 darwin/amd64
ok  	github.com/sttk-go/sabi	0.834s	coverage: 99.6% of statements

Now using version go1.19.5
go version go1.19.5 darwin/amd64
ok  	github.com/sttk-go/sabi	0.836s	coverage: 99.6% of statements

Now using version go1.20
go version go1.20 darwin/amd64
ok  	github.com/sttk-go/sabi	0.843s	coverage: 99.6% of statements

Back to go1.20
Now using version go1.20
%
```


## License

Copyright (C) 2022-2023 Takayuki Sato

This program is free software under MIT License.<br>
See the file LICENSE in this distribution for more details.


[repo-url]: https://github.com/sttk-go/sabi
[pkg-dev-img]: https://pkg.go.dev/badge/github.com/sttk-go/sabi.svg
[pkg-dev-url]: https://pkg.go.dev/github.com/sttk-go/sabi
[ci-img]: https://github.com/sttk-go/sabi/actions/workflows/go.yml/badge.svg?branch=main
[ci-url]: https://github.com/sttk-go/sabi/actions
[mit-img]: https://img.shields.io/badge/license-MIT-green.svg
[mit-url]: https://opensource.org/licenses/MIT
