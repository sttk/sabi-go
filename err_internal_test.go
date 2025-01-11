package sabi

import (
	"container/list"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ClearErrHandlers() {
	syncErrHandlers.head = nil
	syncErrHandlers.last = nil
	asyncErrHandlers.head = nil
	asyncErrHandlers.last = nil
	isErrHandlersFixed = false
}

func TestAddErrSyncHandler(t *testing.T) {
	const fn_sig string = "func(sabi.Err, sabi.ErrOccasion)"

	t.Run("zero handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		assert.Nil(t, syncErrHandlers.head)
		assert.Nil(t, syncErrHandlers.last)
	})

	t.Run("one handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddSyncErrHandler(func(e Err, o ErrOccasion) {})

		assert.NotNil(t, syncErrHandlers.head)
		assert.NotNil(t, syncErrHandlers.last)
		assert.Equal(t, syncErrHandlers.head, syncErrHandlers.last)

		assert.Nil(t, syncErrHandlers.last.next)
		assert.Nil(t, syncErrHandlers.head.next)

		assert.NotNil(t, syncErrHandlers.head.handler)
		assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.handler).String(), fn_sig)
	})

	t.Run("two handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddSyncErrHandler(func(e Err, o ErrOccasion) {})
		AddSyncErrHandler(func(e Err, o ErrOccasion) {})

		assert.NotNil(t, syncErrHandlers.head)
		assert.NotNil(t, syncErrHandlers.last)
		assert.NotEqual(t, syncErrHandlers.head, syncErrHandlers.last)

		assert.Equal(t, syncErrHandlers.head.next, syncErrHandlers.last)
		assert.Nil(t, syncErrHandlers.last.next)

		assert.NotNil(t, syncErrHandlers.head.handler)
		assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.handler).String(), fn_sig)

		assert.NotNil(t, syncErrHandlers.head.next.handler)
		assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.next.handler).String(), fn_sig)
	})
}

func TestAddErrAsyncHandler(t *testing.T) {
	const fn_sig string = "func(sabi.Err, sabi.ErrOccasion)"

	t.Run("zero handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		assert.Nil(t, asyncErrHandlers.head)
		assert.Nil(t, asyncErrHandlers.last)
	})

	t.Run("one handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddAsyncErrHandler(func(e Err, o ErrOccasion) {})

		assert.NotNil(t, asyncErrHandlers.head)
		assert.NotNil(t, asyncErrHandlers.last)
		assert.Equal(t, asyncErrHandlers.head, asyncErrHandlers.last)

		assert.Nil(t, asyncErrHandlers.last.next)
		assert.Nil(t, asyncErrHandlers.head.next)

		assert.NotNil(t, asyncErrHandlers.head.handler)
		assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.handler).String(), fn_sig)
	})

	t.Run("two handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddAsyncErrHandler(func(e Err, o ErrOccasion) {})
		AddAsyncErrHandler(func(e Err, o ErrOccasion) {})

		assert.NotNil(t, asyncErrHandlers.head)
		assert.NotNil(t, asyncErrHandlers.last)
		assert.NotEqual(t, asyncErrHandlers.head, asyncErrHandlers.last)

		assert.Equal(t, asyncErrHandlers.head.next, asyncErrHandlers.last)
		assert.Nil(t, asyncErrHandlers.last.next)

		assert.NotNil(t, asyncErrHandlers.head.handler)
		assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.handler).String(), fn_sig)

		assert.NotNil(t, asyncErrHandlers.head.next.handler)
		assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.next.handler).String(), fn_sig)
	})
}

func TestFixErrHandlers(t *testing.T) {
	t.Run("cannot add any more handlers after fixed", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddSyncErrHandler(func(e Err, o ErrOccasion) {})
		AddAsyncErrHandler(func(e Err, o ErrOccasion) {})

		assert.NotNil(t, syncErrHandlers.head)
		assert.NotNil(t, syncErrHandlers.last)
		assert.Equal(t, syncErrHandlers.head, syncErrHandlers.last)
		assert.Nil(t, syncErrHandlers.last.next)
		assert.Nil(t, syncErrHandlers.head.next)

		assert.NotNil(t, asyncErrHandlers.head)
		assert.NotNil(t, asyncErrHandlers.last)
		assert.Equal(t, asyncErrHandlers.head, asyncErrHandlers.last)
		assert.Nil(t, asyncErrHandlers.last.next)
		assert.Nil(t, asyncErrHandlers.head.next)

		assert.False(t, isErrHandlersFixed)

		FixErrHandlers()

		assert.True(t, isErrHandlersFixed)

		AddSyncErrHandler(func(e Err, o ErrOccasion) {})
		AddAsyncErrHandler(func(e Err, o ErrOccasion) {})

		assert.NotNil(t, syncErrHandlers.head)
		assert.NotNil(t, syncErrHandlers.last)
		assert.Equal(t, syncErrHandlers.head, syncErrHandlers.last)
		assert.Nil(t, syncErrHandlers.last.next)
		assert.Nil(t, syncErrHandlers.head.next)

		assert.NotNil(t, asyncErrHandlers.head)
		assert.NotNil(t, asyncErrHandlers.last)
		assert.Equal(t, asyncErrHandlers.head, asyncErrHandlers.last)
		assert.Nil(t, asyncErrHandlers.last.next)
		assert.Nil(t, asyncErrHandlers.head.next)
	})
}

func TestNotifyErr(t *testing.T) {
	t.Run("when there is no error handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		type FailToDoSomething struct{}

		assert.False(t, isErrHandlersFixed)
		NewErr(FailToDoSomething{})

		FixErrHandlers()
		assert.True(t, isErrHandlersFixed)
		NewErr(FailToDoSomething{})
	})

	t.Run("when there is error handlers", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		syncLogs := list.New()
		asyncLogs := list.New()

		type FailToDoSomething struct{}

		AddSyncErrHandler(func(e Err, o ErrOccasion) {
			syncLogs.PushBack(fmt.Sprintf("%s-1:%s:%d:%s",
				e.Error(), o.File(), o.Line(), o.Time().String()))
		})
		AddSyncErrHandler(func(e Err, o ErrOccasion) {
			syncLogs.PushBack(fmt.Sprintf("%s-2:%s:%d:%s",
				e.Error(), o.File(), o.Line(), o.Time().String()))
		})
		AddAsyncErrHandler(func(e Err, o ErrOccasion) {
			time.Sleep(100 * time.Millisecond)
			asyncLogs.PushBack(fmt.Sprintf("%s-3:%s:%d:%s",
				e.Error(), o.File(), o.Line(), o.Time().String()))
		})
		AddAsyncErrHandler(func(e Err, o ErrOccasion) {
			time.Sleep(10 * time.Millisecond)
			asyncLogs.PushBack(fmt.Sprintf("%s-4:%s:%d:%s",
				e.Error(), o.File(), o.Line(), o.Time().String()))
		})

		assert.False(t, isErrHandlersFixed)

		NewErr(FailToDoSomething{})

		assert.Equal(t, syncLogs.Len(), 0)
		assert.Equal(t, asyncLogs.Len(), 0)

		FixErrHandlers()

		assert.True(t, isErrHandlersFixed)

		NewErr(FailToDoSomething{})

		assert.Equal(t, syncLogs.Len(), 2)
		log := syncLogs.Front()
		assert.Contains(t, log.Value, "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi.FailToDoSomething }-1:err_internal_test.go:218:")
		log = log.Next()
		assert.Contains(t, log.Value, "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi.FailToDoSomething }-2:err_internal_test.go:218:")
		log = log.Next()
		assert.Nil(t, log)

		time.Sleep(500 * time.Millisecond)

		assert.Equal(t, asyncLogs.Len(), 2)
		log = asyncLogs.Front()
		assert.Contains(t, log.Value, "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi.FailToDoSomething }-4:err_internal_test.go:218:")
		log = log.Next()
		assert.Contains(t, log.Value, "github.com/sttk/sabi.Err { reason = github.com/sttk/sabi.FailToDoSomething }-3:err_internal_test.go:218:")
		log = log.Next()
		assert.Nil(t, log)
	})
}
