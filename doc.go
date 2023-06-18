// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

/*
Package github.com/sttk/sabi is a small framework to separate logic parts and data accesses parts for Golang applications.

# Logic

A logic is implemented as a function.
This function takes only an argument, dax which is an interface and collects data access methods used in this function.
Also, this function returns only a sabi.Err value which indicates that this function succeeds or not.
Since the dax hides details of data access procedures, only logical procedure appears in this function.
In this logic part, it's no concern where a data comes from or goes to.

For example, in the following code, GreetLogic is a logic function and GreetDax is a dax interface.

	import "github.com/sttk/sabi"

	type GreetDax interface {
		UserName() (string, sabi.Err)
		Say(greeting string) sabi.Err
	}

	type ( // possible error reasons
		NoName struct {}
		FailToOutput struct {Text string}
	)

	func GreetLogic(dax GreetDax) sabi.Err {
		name, err := dax.UserName()
		if !err.IsOk() {
			return err
		}
		return dax.Say("Hello, " + name)
	}

In GreetLogic function, there are no codes for getting a user name and output a greeting text.
In this logic function, it's only concern to create a greeting text from a user name.

# Dax for unit tests

To test a logic function, the simplest dax implementation is what using a map.
The following code is an example of a dax implementation using a map and having two methods: UserName and Say which are same to GreetDax interface above.

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
		} {
			DaxBase: base,
			mapGreetDax: mapGreetDax{m: m},
		}
	}

And the following code is an example of a test case.

	import (
		"github.com/stretchr/testify/assert"
		"testing"
	)

	func TestGreetLogic_normal(t *testing.T) {
		m := make(map[string]any)
		base := NewMapGreetDaxBase(m)

		m["username"] = "World"
		err := sabi.RunTxn(base, GreetLogic)
		assert.Equal(t, m["greeting"], "Hello, World")
	}

# Dax for real data accesses

In actual case, multiple data sources are often used.
In this example, an user name is input as command line argument, and greeting is output to standard output (console output).
Therefore, two dax implementations are attached to the single GreetDax interface.

The following code is an example of a dax implementation which inputs an user name from command line argument.

	import "os"

	type CliArgsUserDax struct {
	}

	func (dax CliArgsUserDax) UserName() (string, sabi.Err) {
		if len(os.Args) <= 1 {
			return "", sabi.NewErr(NoName{})
		}
		return os.Args[1], sabi.Ok()
	}

In addition, the following code is an example of a dax implementation which outputs a greeting test to console.

	import "fmt"

	type ConsoleOutputDax struct {
	}

	func (dax ConsoleOutputDax) Say(text string) sabi.Err {
		_, e := fmt.Println(text)
		if e != nil {
			return sabi.NewErr(FailToOutput{Text: text}, e)
		}
		return sabi.Ok()
	}

And these dax implementations are combined to a DaxBase as follows:

	func NewGreetDaxBase() sabi.DaxBase {
		base := sabi.NewDaxBase()
		return struct {
			sabi.DaxBase
			CliArgsUserDax
			ConsoleOutputDax
		} {
			DaxBase: base,
			CliArgsUserDax: CliArgsUserDax{},
			ConsoleOutputDax: ConsoleOutputDax{},
		}
	}

# Executing logic

The following code implements a main function which execute a GreetLogic.
sabi.RunTxn executes the GreetLogic function in a transaction process.

	import "log"

	func main() {
		base := NewGreetDaxBase()
		err := sabi.RunTxn(base, GreetLogic)
		if !err.IsOk() {
			log.Fatalln(err.Reason())
		}
	}

# Moving outputs to another transaction process

sabi.RunTxn executes logic functions in a transaction. If a logic function updates database and causes an error in the transaction, its update is rollbacked.
If console output is executed in the same transaction with database update, the rollbacked result is possible to be output to console.
Therefore, console output is wanted to execute after the transaction of database update is successfully completed.

What should be done to achieve it are to add a dax interface for next transaction, to change ConsoleOutputDax to hold a greeting text in Say method, to add a new method to output it in next transaction, and to execute the next transaction in the main function.

	type PrintDax interface {  // Added.
		Print() sabi.Err
	}

	type ConsoleOutputDax struct {
		text string  // Added
	}
	func (dax *ConsoleOutputDax) Say(text string) sabi.Err { // Changed to pointer
		dax.text = text  // Changed
		return sabi.Ok()
	}
	func (dax *ConsoleOutputDax) Print() sabi.Err {  // Added
		_, e := fmt.Println(dax.text)
		if e != nil {
			return sabi.NewErr(FailToOutput{Text: dax.text}, e)
		}
		return sabi.Ok()
	}

	func NewGreetDaxBase() sabi.DaxBase {
		base := sabi.NewDaxBase()
		return struct {
			sabi.DaxBase
			CliArgsUserDax
			*ConsoleOutputDax    // Changed
		}{
			DaxBase:           base,
			CliArgsUserDax:    CliArgsUserDax{},
			ConsoleOutputDax2: &ConsoleOutputDax2{}, // Changed
		}
	}

And the main function is modified as follows:

	func main() {
		err := sabi.RunTxn(base, GreetLogic)
		if !err.IsOk() {
			log.Fatalln(err.Reason())
		}
		err = sabi.RunTxn(base, func(dax PrintDax) sabi.Err {
			return dax.Print()
		})
		if !err.IsOk() {
			log.Fatalln(err.Reason())
		}
	}

Or, the main function is able to rewrite as follows:

	func main() {
		txn0 := sabi.Txn(base, GreetLogic)
		txn1 := sabi.Txn(base, func(dax PrintDax) sabi.Err {
			return dax.Print()
		})
		err := sabi.RunSeq(txn0, txn1)
		if !err.IsOk() {
			log.Fatalln(err.Reason())
		}
	}

The important point is that the GreetLogic function is not changed.
Since this change is not related to the application logic, it is confined to the data access part only.
*/
package sabi
