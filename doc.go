// Copyright (C) 2022-2023 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

/*
Package github.com/sttk-go/sabi is a small framework to separate logic parts and data accesses parts for Golang applications.

Logic

A logic part is implemented as a function.
This function takes a dax, which is an abbreviation of 'data access', as an argument.
A dax has all methods to be used in a logic, and each method is associated with each data access procedure to target data sources.
Since a dax conceals its data access procedures, only logical procedure appears in a logic part.
In a logic part, these are no concern where a data comes from and a data goes to.

For example, Greet is a logic and GreetDax is a dax interface:

    type GreetDax interface {
      GetName() (string, sabi.Err)
      Say(greeting string) sabi.Err
    }

    func GreetLogic(dax GreetDax) sabi.Err {
      name, err := dax.GetName()
      if !err.IsOk() {
        return err
      }
      return dax.Say("Hello, " + name)
    }

In Greet function, there are no detail codes for getting name and putting greeting.
In this logic function, it's only concern to convert a name to a greeting.

Dax for unit tests

To test a logic, the simplest dax implementation is what using a map.
The following code is an example which implements two methods: GetName and Say which are same to GreetDax interface above.

  type mapDax struct {
    m map[string]string
  }

  type (
    NoName struct {}, // An error reason when getting no name.
  )

  func (dax mapDax) GetName() (string, sabi.Err) {
    name, exists := dax.m["name"]
    if !exists {
      return "", sabi.ErrBy(NoName{})
    }
    return name, sabi.Ok()
  }

  func (dax mapDax) Say(greeting string) sabi.Err {
    dax.m["greeting"] = greeting
    return sabi.Ok()
  }

And the following code is an example of a test case.

  func TestGreetLogic(t *testing.T) {
    m := make(map[string]string)
    dax := mapDax{m: m}

    base := sabi.NewDaxBase()
    proc := sabi.NewProc[GreetDax](base, dax)

    m["name"] = "World"
    err := proc.RunTxn(GreetLogic)
    assert.True(t, err.IsOk())
    assert.Equal(t, m["greeting"], "Hello, World")
  }

Dax for real data access

An actual dax ordinarily consists of multiple sub dax by input sources and output destinations.

The following code is an example of a dax with no external data source.
This dax outputs a greeting to standard output.

  type SayConsoleDax struct {}

  type (
    FailToPrint struct {}
  )

  func (dax SayConsoleDax) Say(text string) sabi.Err {
    _, e := fmt.Println(text)
    if e != nil {
      return sabi.ErrBy(FailToPrint{}, e)
    }
    return sabi.Ok()
  }

And the following code is an example of a dax with an external data source.
This dax accesses to a database and provides an implementation of GetName method of GreetDax.

  type UserSqlDax struct {
    sqldax.SqlDax
  }

  type (
    FailToCreateStmt struct {}
    NoUser struct {}
    FailToQueryUserName struct {}
  )

  func (dax UserSqlDax) GetName() (string, sabi.Err) {
    conn, err := dax.GetSqlDaxConn("sql")
    if !err.IsOk() {
      return err
    }
    stmt, err := conn.Prepare("SELECT username FROM users LIMIT 1")
    if err != nil {
      return "", sabi.ErrBy(FailToCreateStmt{})
    }
    defer stmt.Close()

    var username string
    err = stmt.QueryRow().Scan(&username)
    switch {
    case err == sql.ErrNoRows:
      return "", sabi.ErrBy(NoUser{})
    case err != nil:
      return "", sabi.ErrBy(FailToQueryUserName{})
    default:
      return username, sabi.Ok()
    }
  }

Mapping dax interface and implementations

A dax interface can be related to multiple dax implementations.

In the following code, GetName method of GreetDax interface is corresponded to the same named method of UserSqlDax, and Say method of GreetDax interface is corresponded to the same named method of SayConsoleDax.

  func NewGreetProc() sabi.Proc[GreetDax] {
    base := sabi.NewDaxBase()

    dax := struct {
      UserSqlDax
      SayConsoleDax
    } {
      UserSqlDax: UserSqlDax{SqlDax: sqldax.NewSqlDax(base)},
      SayConsoleDax: SayConsoleDax{},
    }

    return sabi.NewProc[GreetDax](base, dax)
  }

Executing logic

The following code implements a main function which execute a GreetLogic.
GreetLogic is executed in a transaction process by GreetProc#RunTxn, so the database update can be rollbacked when an error is occured.

The init function registers a SqlDaxSrc which creates a DaxConn which connects to a database. The SqlDaxConn is registerd with a name "sql" and is obtained by GetSqlDaxConn("sql") in UserSqlDax#GetName.

  func init() {
    sabi.AddGlobalDaxSrc("sql", sqldax.SqlDaxSrc{driver: "driver-name", dsn: "ds-name"})
    sabi.FixGlobalDaxSrcs()
  }

  func main() {
    proc := NewGreetProc()
    err := proc.RunTxn(GreetLogic)
    switch err.Reason().(type) {
    default:
      os.Exit(1)
    case sabi.NoError:
    }
  }

*/
package sabi
