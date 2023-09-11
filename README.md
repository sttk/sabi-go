# [Sabi][repo-url] [![Go Reference][pkg-dev-img]][pkg-dev-url] [![CI Status][ci-img]][ci-url] [![MIT License][mit-img]][mit-url]

A small framework to separate logics and data accesses for Golang application.

## Concepts

The concept of this framework is separation and reintegration of necessary and
redundant parts based on the perspectives of the whole and the
parts.
The separation of logics and data accesses is the most prominent and
fundamental part of this concept.

### Separation of logics and data accesses

In general, a program consists of procedures and data. And procedures include data accesses for operating data, and the rest of procedures are logics. So we can say that a program consists of logics, data accesses and data.

We often think to separate an application to multiple layers, for example, controller layer, business logic layer, and data access layer.
The logics and data accesses mentioned in this framework may appear to follow such layering.
However, the controller layer also has data accesses such as transforming user requests and responses for the business logic layer.
Generally, such layers of an application is established as vertical stages of data processing within a data flow.

In this framework, the relationship between logics and data accesses is not
defined by layers but by lanes.
Their relationship is vertical in terms of invocation, but is horizontal conceptually.
DaxBase serves as an intermediary that connects both of them.

### Separation of data accesses for each logic

A logic is a function that takes dax interface as its only one argument.
The type of this dax is a type parameter of the logic function, and also a type
parameter of the transaction function, Txn, that executes logics.

Therefore, since the type of dax can be changed for each logic or transaction,
it is possible to limit data accesses used by the logic, by declaring only
necessary data access methods from among ones defined in DaxBase instance..

At the same time, since all data accesses of a logic is done through this sole
dax interface, this dax interface serves as a list of data access methods used
by a logic.

### Separation of data accesses by data sources and reintegration of them

Data access methods are implemented as methods of some Dax structs that
embedding a DaxBase.
Furthermore these Dax structs can be integrated into a single new DaxBase.

A Dax struct can be created in any unit, but it is clearer to create it in the
unit of the data source.
By doing so, the definition of a new DaxBase also serves as a list of the data
sources being used.

## Import declaration

To use this package in your code, the following import declaration is necessary.

```
import (
    "github.com/sttk/sabi"
    "github.com/sttk/sabi/errs"
)
```

## Usage

### Logic and an interface for its data access

A logic is implemented as a function.
This function takes only an argument, dax, which is an interface that gathers
only the data access methods needed by this logic function.

Since dax conceals details of data access procedures, this function only
includes logical procedures.
In this logical part, there is no concern about where the data is input from or where it is output to.

For example, in the following code, GreetLogic is a logic function and GreetDax
is a dax interface for GreetLogic.

```
type ( // possible error reasons
    NoName struct {}
    FailToGetHour struct {}
    FailToOutput struct { Text string }
)

type GreetDax interface {
    UserName() (string, errs.Err)
    Hour() int
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

    return dax.Output(name)
}
```

In GreetLogic, there are no codes for inputting the hour, inputting a user name,
and outputing a greeting.
This function has only concern to create a greeting text.

### Data accesses for unit testing

To test a logic function, the simplest dax struct is what using a map.
The following code is an example of a dax struct using a map and having
three methods that are same to GreetDax interface methods above.

```
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
    if m["greeting"] == "error" { // for testing the error case.
      return errs.New(FailToOutput{Text: text})
    }
    dax.m["greeting"] = text
    return errs.Ok()
}

func NewMapGreetDaxBase(m map[string]any) sabi.DaxBase {
    base := sabi.NewDaxBase()
    return struct {
      sabi.DaxBase
      MapGreetDax
    } {
      DaxBase: base,
      MapGreetDax: MapGreetDax{m: m},
    }
}
```

And the following code is an example of a test case.

```
func TestGreetLogic_morning(t *testing.T) {
    m := make(map[string]any)
    base := NewGreetDaxBase(m)

    m["username"] = "everyone"
    m["hour"] = 10
    err := sabi.Txn[GreetDax](base, GreetLogic)
    if err.IsNotOk() {
      t.Errorf(err.Error())
    }
    if m["greeting"] == "Good morning, everyone" {
      t.Errorf("Bad greeting: %v\n", m["greeting"])
    }
}
```

### Data accesses for actual use

In actual use, multiple data sources are often used.
In this example, an user name and the hour are input as command line argument,
and greeting is output to console.
Therefore, two dax struct are created and they are integrated into a new struct
based on DaxBase.
Since Golang is structural typing language, this new DaxBase can be casted to
GreetDax.

The following code is an example of a dax struct which inputs an user name and
the hour from command line argument.

```
type CliArgsDax struct {
    sabi.Dax
}

func (dax CliArgsDax) UserName() (string, errs.Err) {
    if len(os.Args) <= 1 {
      return "", errs.New(NoName{})
    }
    return os.Args[1], errs.Ok()
}

func (dax CliArgsDax) Hour() (string, errs.Err) {
    if len(os.Args) <= 2 {
      return 0, errs.New(FailToGetHour{})
    }
    return os.Args[2], errs.Ok()
}
```

The following code is an example of a dax struct which output a text to
console.

```
type ConsoleDax struct {
    sabi.Dax
}

func (dax ConsoleDax) Output(text string) errs.Err {
    fmt.Println(text)
    return errs.Ok()
}
```

And the following code is an example of a constructor function of a struct
based on DaxBase into which the above two dax are integrated.
This implementation also serves as a list of the external data sources being
used.

```
func NewGreetDaxBase() sabi.DaxBase {
    base := sabi.NewDaxBase()
    return struct {
        sabi.DaxBase
        CliArgsDax
        ConsoleDax
    } {
        DaxBase: base,
        CliArgsDax: CliArgsDax{Dax: base},
        ConsoleDax: ConsoleDax{Dax: base},
    }
}
```

### Executing a logic

The following code executes the above GreetLogic in a transaction process.

```
func main() {
    if err := sabi.StartApp(app); err.IsNotOk() {
        fmt.Println(err.Error())
        os.Exit(1)
    }
}

func app() errs.Err {
    base := NewBase()
    defer base.Close()

    return sabi.Txn(base, GreetLogic))
}
```

### Changing to a dax of another data source

In the above codes, the hour is obtained from command line arguments.
Here, assume that the specification has been changed to retrieve it
from system clock instread.

In this case, we can solve this by removing the Hour method from CliArgsDax
and creating a new Dax, SystemClockDax, which has Hour method to retrieve
a hour from system clock.

```
// func (dax CliArgsDax) Hour() (string, errs.Err) {  // Removed
//     if len(os.Args) <= 2 {
//       return 0, errs.New(FailToGetHour{})
//     }
//     return os.Args[2], errs.Ok()
// }

type SystemClockDax struct {  // Added
    sabi.Dax
}

func (dax SystemClockTimeDax) Hour() (string, errs.Err) {  // Added
    return time.Now().Hour(), errs.Ok()
}
```

And the DaxBase struct, into which multiple dax structs have been integrated,
is modified as follows.

```
func NewGreetDaxBase() sabi.DaxBase {
    base := sabi.NewDaxBase()
    return struct {
        sabi.DaxBase
        CliArgsDax
        SystemClockDax  // Added
        ConsoleDax
    } {
        DaxBase: base,
        CliArgsDax: CliArgsDax{Dax: base},
        SystemClockDax: SystemClockDax{Dax: base},  // Added
        ConsoleDax: ConsoleDax{Dax: base},
    }
}
```

### Moving outputs to next transaction process.

The above codes works normally if no error occurs.
But if an error occurs at getting user name, a incomplete string is being
output to console.
Such behavior is not appropriate for transaction processing.

So we should change the above codes to store in memory temporarily in the
existing transaction process, and output to console in the next transaction.

The following code is the implementation of MemoryDax which is memory store
dax and the DaxBase struct after replacing ConsoleDax to MemoryDax.

```
type MemoryDax struct {  // Added
    sabi.Dax
    text string
}

func (dax *MemoryDax) Output(text string) errs.Err {  // Added
    dax.text = text
    return errs.Ok()
}

func (dax *MemoryDax) GetText() string {  // Added
    return dax.text
}

func (dax ConsoleDax) Print(text string) errs.Err {  // Changed from Output
    fmt.Println(text)
    return errs.Ok()
}

func NewGreetDaxBase() sabi.DaxBase {
    base := sabi.NewDaxBase()
    return struct {
        sabi.DaxBase
        CliArgsDax
        SystemClockDax
        ConsoleDax
        MemoryDax  // Added
    } {
        DaxBase: base,
        CliArgsDax: CliArgsDax{Dax: base},
        SystemClockDax: SystemClockDax{Dax: base},
        ConsoleDax: ConsoleDax{Dax: base},
        MemoryDax: MemoryDax{Dax: base},  // Added
    }
}
```

The following code is the logic to output text to console in next transaction
process, the dax interface for the logic, and the execution of logics after
being changed.

```
type PrintDax interface {  // Added
    GetText() string
    Print(text string) errs.Err
}

func PrintLogic(dax PrintDax) errs.Err {  // Added
    text := dax.GetText()
    return dax.Print(text)
}

func app() errs.Err {
    base := NewBase()
    defer base.Close()

    return sabi.Txn(base, GreetLogic)).    // Changed
        IfOk(sabi.Txn_(base, PrintLogic))  // Added
}
```

The important point is that the GreetLogic function is not changed.
Since these changes are not related to the existing application logic, it is
limited to the data access part (and the part around the newly added logic)
only.


## Supporting Go versions

This framework supports Go 1.18 or later.

### Actual test results for each Go version:

```
% gvm-fav
Now using version go1.18.10
go version go1.18.10 darwin/amd64
ok  	github.com/sttk/sabi	0.614s	coverage: 100.0% of statements
ok  	github.com/sttk/sabi/errs	0.793s	coverage: 100.0% of statements

Now using version go1.19.13
go version go1.19.13 darwin/amd64
ok  	github.com/sttk/sabi	0.559s	coverage: 100.0% of statements
ok  	github.com/sttk/sabi/errs	0.755s	coverage: 100.0% of statements

Now using version go1.20.8
go version go1.20.8 darwin/amd64
ok  	github.com/sttk/sabi	0.566s	coverage: 100.0% of statements
ok  	github.com/sttk/sabi/errs	0.833s	coverage: 100.0% of statements

Now using version go1.21.1
go version go1.21.1 darwin/amd64
ok  	github.com/sttk/sabi	0.572s	coverage: 100.0% of statements
ok  	github.com/sttk/sabi/errs	0.835s	coverage: 100.0% of statements

Back to go1.21.1
Now using version go1.21.1
```

## License

Copyright (C) 2022-2023 Takayuki Sato

This program is free software under MIT License.<br>
See the file LICENSE in this distribution for more details.


[repo-url]: https://github.com/sttk/sabi-go
[pkg-dev-img]: https://pkg.go.dev/badge/github.com/sttk/sabi.svg
[pkg-dev-url]: https://pkg.go.dev/github.com/sttk/sabi
[ci-img]: https://github.com/sttk/sabi-go/actions/workflows/go.yml/badge.svg?branch=main
[ci-url]: https://github.com/sttk/sabi-go/actions
[mit-img]: https://img.shields.io/badge/license-MIT-green.svg
[mit-url]: https://opensource.org/licenses/MIT
