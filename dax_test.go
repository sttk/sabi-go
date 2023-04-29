package sabi_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sttk-go/sabi"
	"testing"
)

func TestAddGlobalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 0)

	sabi.AddGlobalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 1)

	sabi.AddGlobalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 2)
}

func TestStartUpGlobalDaxSrcs(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 0)

	sabi.AddGlobalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 1)

	err := sabi.StartUpGlobalDaxSrcs()
	assert.True(t, err.IsOk())

	assert.True(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 1)

	sabi.AddGlobalDaxSrc("bar", &BarDaxSrc{})

	assert.True(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 1)

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	assert.Nil(t, log.Next())
}

func TestStartUpGlobalDaxSrcs_failToSetUpDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	WillFailToSetUpFooDaxSrc = true

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 0)

	sabi.AddGlobalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 1)

	sabi.AddGlobalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 2)

	err := sabi.StartUpGlobalDaxSrcs()
	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case sabi.FailToStartUpGlobalDaxSrcs:
		errs := err.Reason().(sabi.FailToStartUpGlobalDaxSrcs).Errors
		assert.Equal(t, len(errs), 1)
		err1 := errs["foo"]
		r := err1.Reason().(FailToDoSomething)
		assert.Equal(t, r.Text, "FailToSetUpFooDaxSrc")
	default:
		assert.Fail(t, err.Error())
	}

	log := Logs.Front()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	if log.Value == "FooDaxSrc#End" {
		assert.Equal(t, log.Value, "FooDaxSrc#End")
		log = log.Next()
		assert.Equal(t, log.Value, "BarDaxSrc#End")
	} else {
		assert.Equal(t, log.Value, "BarDaxSrc#End")
		log = log.Next()
		assert.Equal(t, log.Value, "FooDaxSrc#End")
	}
	assert.Nil(t, log.Next())
}

func TestShutdownGlobalDaxSrcs(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 0)

	sabi.AddGlobalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 1)

	sabi.AddGlobalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 2)

	err := sabi.StartUpGlobalDaxSrcs()
	assert.True(t, err.IsOk())

	assert.True(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 2)

	sabi.ShutdownGlobalDaxSrcs()

	assert.True(t, sabi.IsGlobalDaxSrcsFixed())
	assert.Equal(t, len(sabi.GlobalDaxSrcMap()), 2)

	log := Logs.Front()
	if log.Value == "FooDaxSrc#SetUp" {
		assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
		log = log.Next()
		assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	} else {
		assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
		log = log.Next()
		assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	}
	log = log.Next()
	if log.Value == "FooDaxSrc#End" {
		assert.Equal(t, log.Value, "FooDaxSrc#End")
		log = log.Next()
		assert.Equal(t, log.Value, "BarDaxSrc#End")
	} else {
		assert.Equal(t, log.Value, "BarDaxSrc#End")
		log = log.Next()
		assert.Equal(t, log.Value, "FooDaxSrc#End")
	}
	assert.Nil(t, log.Next())
}

func TestDaxBase_SetUpLocalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 2)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDaxBase_SetUpLocalDaxSrc_unableToAddLocalDaxSrcInTxn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	sabi.Begin(base)

	assert.True(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})

	assert.True(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Nil(t, log)

	sabi.End(base)

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 2)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	log = Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDaxBase_SetUpLocalDaxSrc_failToSetUpDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	WillFailToSetUpFooDaxSrc = true

	base := sabi.NewDaxBase()

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	err := base.SetUpLocalDaxSrc("bar", BarDaxSrc{})
	assert.True(t, err.IsOk())

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	err = base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.False(t, err.IsOk())

	switch err.Reason().(type) {
	case FailToDoSomething:
		r := err.Reason().(FailToDoSomething)
		assert.Equal(t, r.Text, "FailToSetUpFooDaxSrc")
	default:
		assert.Fail(t, err.Error())
	}

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.FreeAllLocalDaxSrcs()

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	log := Logs.Front()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#End")
	assert.Nil(t, log.Next())
}

func TestDaxBase_FreeLocalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.FreeLocalDaxSrc("foo")

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("bar", &BarDaxSrc{})

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 2)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.FreeLocalDaxSrc("bar")

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.FreeLocalDaxSrc("foo")

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#End")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#End")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#End")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDaxBase_FreeLocalDaxSrc_unableToFreeLocalDaxSrcInTxn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.SetUpLocalDaxSrc("foo", FooDaxSrc{})

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	sabi.Begin(base)

	assert.True(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.FreeLocalDaxSrc("foo")

	assert.True(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	sabi.End(base)

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 1)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)

	base.FreeLocalDaxSrc("foo")

	assert.False(t, sabi.IsLocalDaxSrcsFixed(base))
	assert.Equal(t, len(sabi.LocalDaxSrcMap(base)), 0)
	assert.Equal(t, len(sabi.DaxConnMap(base)), 0)
}

func TestDaxBase_GetDaxConn_withLocalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	conn, err := base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case sabi.DaxSrcIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	err = sabi.StartUpGlobalDaxSrcs()

	conn, err = base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case sabi.DaxSrcIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	err = base.SetUpLocalDaxSrc("foo", FooDaxSrc{})

	conn, err = base.GetDaxConn("foo")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	var conn2 sabi.DaxConn
	conn2, err = base.GetDaxConn("foo")
	assert.Equal(t, conn2, conn)
	assert.True(t, err.IsOk())
}

func TestDaxBase_GetDaxConn_withGlobalDaxSrc(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	conn, err := base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case sabi.DaxSrcIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	sabi.AddGlobalDaxSrc("foo", FooDaxSrc{})

	err = sabi.StartUpGlobalDaxSrcs()
	conn, err = base.GetDaxConn("foo")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	var conn2 sabi.DaxConn
	conn2, err = base.GetDaxConn("foo")
	assert.Equal(t, conn2, conn)
	assert.True(t, err.IsOk())
}

func TestDaxBase_GetDaxConn_localDsIsTakenPriorityOfGlobalDs(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	conn, err := base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case sabi.DaxSrcIsNotFound:
		assert.Equal(t, err.Get("Name"), "foo")
	default:
		assert.Fail(t, err.Error())
	}

	sabi.AddGlobalDaxSrc("foo", FooDaxSrc{Label: "global"})

	err = sabi.StartUpGlobalDaxSrcs()
	assert.True(t, err.IsOk())

	err = base.SetUpLocalDaxSrc("foo", FooDaxSrc{Label: "local"})
	assert.True(t, err.IsOk())

	conn, err = base.GetDaxConn("foo")
	assert.Equal(t, conn.(FooDaxConn).Label, "local")
	assert.True(t, err.IsOk())
}

func TestDaxBase_GetDaxConn_failToCreateDaxConn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	WillFailToCreateFooDaxConn = true

	base := sabi.NewDaxBase()

	err := base.SetUpLocalDaxSrc("foo", FooDaxSrc{})

	var conn sabi.DaxConn
	conn, err = base.GetDaxConn("foo")
	assert.Nil(t, conn)
	switch err.Reason().(type) {
	case sabi.FailToCreateDaxConn:
		assert.Equal(t, err.Get("Name"), "foo")
		switch err.Cause().(sabi.Err).Reason().(type) {
		case FailToDoSomething:
		default:
			assert.Fail(t, err.Error())
		}
	default:
		assert.Fail(t, err.Error())
	}
}

func TestDaxBase_GetDaxConn_commit(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	err := base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.True(t, err.IsOk())

	err = base.SetUpLocalDaxSrc("bar", BarDaxSrc{})
	assert.True(t, err.IsOk())

	sabi.Begin(base)

	var conn sabi.DaxConn
	conn, err = base.GetDaxConn("foo")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	conn, err = base.GetDaxConn("bar")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	err = sabi.Commit(base)
	assert.True(t, err.IsOk())

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#CreateDaxConn")
	log = log.Next()
	if log.Value == "FooDaxConn#Commit" {
		assert.Equal(t, log.Value, "FooDaxConn#Commit")
		log = log.Next()
		assert.Equal(t, log.Value, "BarDaxConn#Commit")
	} else {
		assert.Equal(t, log.Value, "BarDaxConn#Commit")
		log = log.Next()
		assert.Equal(t, log.Value, "FooDaxConn#Commit")
	}
	log = log.Next()
	assert.Nil(t, log)
}

func TestDaxBase_GetDaxConn_failToCommit(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	WillFailToCommitFooDaxConn = true

	base := sabi.NewDaxBase()

	err := base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.True(t, err.IsOk())

	err = base.SetUpLocalDaxSrc("bar", BarDaxSrc{})
	assert.True(t, err.IsOk())

	sabi.Begin(base)

	var conn sabi.DaxConn
	conn, err = base.GetDaxConn("foo")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	conn, err = base.GetDaxConn("bar")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	err = sabi.Commit(base)
	assert.False(t, err.IsOk())
	switch err.Reason().(type) {
	case sabi.FailToCommitDaxConn:
		m := err.Get("Errors").(map[string]sabi.Err)
		assert.Equal(t, m["foo"].ReasonName(), "FailToDoSomething")
	default:
		assert.Fail(t, err.Error())
	}

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxConn#Commit")
	log = log.Next()
	assert.Nil(t, log)
}

func TestDaxBase_GetDaxConn_rollback(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	err := base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.True(t, err.IsOk())
	err = base.SetUpLocalDaxSrc("bar", BarDaxSrc{})
	assert.True(t, err.IsOk())

	sabi.Begin(base)

	var conn sabi.DaxConn
	conn, err = base.GetDaxConn("foo")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	conn, err = base.GetDaxConn("bar")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	sabi.Rollback(base)

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#CreateDaxConn")
	log = log.Next()
	if log.Value == "FooDaxConn#Rollback" {
		assert.Equal(t, log.Value, "FooDaxConn#Rollback")
		log = log.Next()
		assert.Equal(t, log.Value, "BarDaxConn#Rollback")
	} else {
		assert.Equal(t, log.Value, "BarDaxConn#Rollback")
		log = log.Next()
		assert.Equal(t, log.Value, "FooDaxConn#Rollback")
	}
	log = log.Next()
	assert.Nil(t, log)
}

func TestDaxBase_GetDaxConn_close(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	base := sabi.NewDaxBase()

	err := base.SetUpLocalDaxSrc("foo", FooDaxSrc{})
	assert.True(t, err.IsOk())

	err = base.SetUpLocalDaxSrc("bar", BarDaxSrc{})
	assert.True(t, err.IsOk())

	sabi.Begin(base)

	var conn sabi.DaxConn
	conn, err = base.GetDaxConn("foo")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	conn, err = base.GetDaxConn("bar")
	assert.NotNil(t, conn)
	assert.True(t, err.IsOk())

	err = sabi.Commit(base)
	assert.True(t, err.IsOk())

	sabi.End(base)
	assert.True(t, err.IsOk())

	base.FreeAllLocalDaxSrcs()

	log := Logs.Front()
	assert.Equal(t, log.Value, "FooDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#SetUp")
	log = log.Next()
	assert.Equal(t, log.Value, "FooDaxSrc#CreateDaxConn")
	log = log.Next()
	assert.Equal(t, log.Value, "BarDaxSrc#CreateDaxConn")
	log = log.Next()
	if log.Value == "FooDaxConn#Commit" {
		assert.Equal(t, log.Value, "FooDaxConn#Commit")
		log = log.Next()
		assert.Equal(t, log.Value, "BarDaxConn#Commit")
	} else {
		assert.Equal(t, log.Value, "BarDaxConn#Commit")
		log = log.Next()
		assert.Equal(t, log.Value, "FooDaxConn#Commit")
	}
	log = log.Next()
	if log.Value == "FooDaxConn#Close" {
		assert.Equal(t, log.Value, "FooDaxConn#Close")
		log = log.Next()
		assert.Equal(t, log.Value, "BarDaxConn#Close")
	} else {
		assert.Equal(t, log.Value, "BarDaxConn#Close")
		log = log.Next()
		assert.Equal(t, log.Value, "FooDaxConn#Close")
	}
	log = log.Next()
	if log.Value == "FooDaxSrc#End" {
		assert.Equal(t, log.Value, "FooDaxSrc#End")
		log = log.Next()
		assert.Equal(t, log.Value, "BarDaxSrc#End")
	} else {
		assert.Equal(t, log.Value, "BarDaxSrc#End")
		log = log.Next()
		assert.Equal(t, log.Value, "FooDaxSrc#End")
	}
	log = log.Next()
	assert.Nil(t, log)
}

func TestDax_runTxn(t *testing.T) {
	ClearDaxBase()
	defer ClearDaxBase()

	hogeDs := NewMapDaxSrc()
	fugaDs := NewMapDaxSrc()
	piyoDs := NewMapDaxSrc()

	base := NewHogeFugaPiyoDaxBase()

	var err sabi.Err
	if err = base.SetUpLocalDaxSrc("hoge", hogeDs); err.IsNotOk() {
		assert.Fail(t, err.Error())
		return
	}
	if err = base.SetUpLocalDaxSrc("fuga", fugaDs); err.IsNotOk() {
		assert.Fail(t, err.Error())
		return
	}
	if err = base.SetUpLocalDaxSrc("piyo", piyoDs); err.IsNotOk() {
		assert.Fail(t, err.Error())
		return
	}

	hogeDs.dataMap["hogehoge"] = "Hello, world"

	if err = sabi.RunTxn(base, HogeFugaLogic); err.IsNotOk() {
		assert.Fail(t, err.Error())
		return
	}
	if err = sabi.RunTxn(base, FugaPiyoLogic); err.IsNotOk() {
		assert.Fail(t, err.Error())
		return
	}

	assert.Equal(t, piyoDs.dataMap["piyopiyo"], "Hello, world")
}
