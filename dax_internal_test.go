package sabi

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	Logs                     list.List
	WillFailToSetupFooDaxSrc bool
	WillFailToSetupBarDaxSrc bool
)

func Reset() {
	FixErrHandlers()

	isGlobalDaxSrcsFixed = false
	globalDaxSrcEntryList.head = nil
	globalDaxSrcEntryList.last = nil

	WillFailToSetupFooDaxSrc = false
	WillFailToSetupBarDaxSrc = false

	Logs.Init()
}

type /* error reasons */ (
	FailToSetupFooDaxSrc struct{}
	FailToSetupBarDaxSrc struct{}
)

type FooDaxSrc struct{}

func (ds FooDaxSrc) Setup(ag *AsyncGroup) Err {
	if WillFailToSetupFooDaxSrc {
		return NewErr(FailToSetupFooDaxSrc{})
	}
	Logs.PushBack("FooDaxSrc#Setup")
	return Ok()
}

func (ds FooDaxSrc) Close() {
	Logs.PushBack("FooDaxSrc#Close")
}

type BarDaxSrc struct{}

func (ds *BarDaxSrc) Setup(ag *AsyncGroup) Err {
	ag.Add(func() Err {
		if WillFailToSetupFooDaxSrc {
			return NewErr(FailToSetupBarDaxSrc{})
		}
		Logs.PushBack("FooDaxSrc#Setup")
		return Ok()
	})
	return Ok()
}

func (ds *BarDaxSrc) Close() {
	Logs.PushBack("BarDaxSrc#Close")
}

func TestUses(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		Reset()
		defer Reset()

		Uses("foo", FooDaxSrc{})

		ent0 := globalDaxSrcEntryList.head
		assert.Equal(t, ent0.name, "foo")
		assert.IsType(t, ent0.ds, FooDaxSrc{})
		assert.False(t, ent0.local)
		assert.False(t, ent0.deleted)
		assert.Nil(t, ent0.prev)
		assert.Nil(t, ent0.next)

		Uses("*foo", &FooDaxSrc{})

		ent0 = globalDaxSrcEntryList.head
		assert.Equal(t, ent0.name, "foo")
		assert.IsType(t, ent0.ds, FooDaxSrc{})
		assert.False(t, ent0.local)
		assert.False(t, ent0.deleted)
		assert.Nil(t, ent0.prev)

		ent1 := ent0.next
		assert.Equal(t, ent1.name, "*foo")
		assert.IsType(t, ent1.ds, &FooDaxSrc{})
		assert.False(t, ent1.local)
		assert.False(t, ent1.deleted)
		assert.Equal(t, ent1.prev, ent0)
		assert.Nil(t, ent1.next)

		Uses("bar", &BarDaxSrc{})

		ent0 = globalDaxSrcEntryList.head
		assert.Equal(t, ent0.name, "foo")
		assert.IsType(t, ent0.ds, FooDaxSrc{})
		assert.False(t, ent0.local)
		assert.False(t, ent0.deleted)
		assert.Nil(t, ent0.prev)

		ent1 = ent0.next
		assert.Equal(t, ent1.name, "*foo")
		assert.IsType(t, ent1.ds, &FooDaxSrc{})
		assert.False(t, ent1.local)
		assert.False(t, ent1.deleted)
		assert.Equal(t, ent1.prev, ent0)

		ent2 := ent1.next
		assert.Equal(t, ent2.name, "bar")
		assert.IsType(t, ent2.ds, &BarDaxSrc{})
		assert.False(t, ent2.local)
		assert.False(t, ent2.deleted)
		assert.Equal(t, ent2.prev, ent1)
		assert.Nil(t, ent2.next)
	})

	t.Run("A name already exists", func(t *testing.T) {
		Reset()
		defer Reset()

		Uses("foo", FooDaxSrc{})

		ent0 := globalDaxSrcEntryList.head
		assert.Equal(t, ent0.name, "foo")
		assert.IsType(t, ent0.ds, FooDaxSrc{})
		assert.False(t, ent0.local)
		assert.False(t, ent0.deleted)
		assert.Nil(t, ent0.prev)
		assert.Nil(t, ent0.next)

		Uses("foo", &BarDaxSrc{})

		ent0 = globalDaxSrcEntryList.head
		assert.Equal(t, ent0.name, "foo")
		assert.IsType(t, ent0.ds, FooDaxSrc{})
		assert.False(t, ent0.local)
		assert.False(t, ent0.deleted)
		assert.Nil(t, ent0.prev)

		ent1 := ent0.next
		assert.Equal(t, ent1.name, "foo")
		assert.IsType(t, ent1.ds, &BarDaxSrc{})
		assert.False(t, ent1.local)
		assert.False(t, ent1.deleted)
		assert.Equal(t, ent1.prev, ent0)
		assert.Nil(t, ent1.next)
	})
}
