package sabi

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type ReasonForNotification struct{}

func ClearErrHandlers() {
	syncErrHandlers.head = nil
	syncErrHandlers.last = nil
	asyncErrHandlers.head = nil
	asyncErrHandlers.last = nil
	isErrCfgsFixed = false
}

func TestAddErrSyncHandler_oneHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddSyncErrHandler(func(err Err, occ ErrOccasion) {})

	assert.NotNil(t, syncErrHandlers.head)
	assert.NotNil(t, syncErrHandlers.last)
	assert.Equal(t, syncErrHandlers.head, syncErrHandlers.last)

	assert.Nil(t, syncErrHandlers.last.next)
	assert.Nil(t, syncErrHandlers.head.next)

	assert.NotNil(t, syncErrHandlers.head.handler)
	assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.handler).String(), "func(sabi.Err, sabi.ErrOccasion)")
}

func TestAddErrSyncHandler_twoHandlers(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddSyncErrHandler(func(err Err, occ ErrOccasion) {})
	AddSyncErrHandler(func(err Err, occ ErrOccasion) {})

	assert.NotNil(t, syncErrHandlers.head)
	assert.NotNil(t, syncErrHandlers.last)
	assert.NotEqual(t, syncErrHandlers.head, syncErrHandlers.last)

	assert.Equal(t, syncErrHandlers.head.next, syncErrHandlers.last)
	assert.Nil(t, syncErrHandlers.last.next)

	assert.NotNil(t, syncErrHandlers.head.handler)
	assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.handler).String(), "func(sabi.Err, sabi.ErrOccasion)")

	assert.NotNil(t, syncErrHandlers.head.next.handler)
	assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.next.handler).String(), "func(sabi.Err, sabi.ErrOccasion)")
}

func TestAddErrAsyncHandler_zeroHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	assert.Nil(t, asyncErrHandlers.head)
	assert.Nil(t, asyncErrHandlers.last)
}

func TestAddErrAsyncHandler_oneHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddAsyncErrHandler(func(err Err, occ ErrOccasion) {})

	assert.NotNil(t, asyncErrHandlers.head)
	assert.NotNil(t, asyncErrHandlers.last)
	assert.Equal(t, asyncErrHandlers.head, asyncErrHandlers.last)

	assert.Nil(t, asyncErrHandlers.last.next)
	assert.Nil(t, asyncErrHandlers.head.next)

	assert.NotNil(t, asyncErrHandlers.head.handler)
	assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.handler).String(), "func(sabi.Err, sabi.ErrOccasion)")
}

func TestAddErrAsyncHandler_twoHandlers(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddAsyncErrHandler(func(err Err, occ ErrOccasion) {})
	AddAsyncErrHandler(func(err Err, occ ErrOccasion) {})

	assert.NotNil(t, asyncErrHandlers.head)
	assert.NotNil(t, asyncErrHandlers.last)
	assert.NotEqual(t, asyncErrHandlers.head, asyncErrHandlers.last)

	assert.Equal(t, asyncErrHandlers.head.next, asyncErrHandlers.last)
	assert.Nil(t, asyncErrHandlers.last.next)

	assert.NotNil(t, asyncErrHandlers.head.handler)
	assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.handler).String(), "func(sabi.Err, sabi.ErrOccasion)")

	assert.NotNil(t, asyncErrHandlers.head.next.handler)
	assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.next.handler).String(), "func(sabi.Err, sabi.ErrOccasion)")
}

func TestFixErrCfgs(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddSyncErrHandler(func(err Err, occ ErrOccasion) {})
	AddAsyncErrHandler(func(err Err, occ ErrOccasion) {})

	assert.NotNil(t, syncErrHandlers.head)
	assert.NotNil(t, syncErrHandlers.last)
	assert.Equal(t, syncErrHandlers.head, syncErrHandlers.last)
	assert.NotNil(t, syncErrHandlers.head.handler)
	assert.Nil(t, syncErrHandlers.head.next)
	assert.Nil(t, syncErrHandlers.last.next)

	assert.NotNil(t, asyncErrHandlers.head)
	assert.NotNil(t, asyncErrHandlers.last)
	assert.Equal(t, asyncErrHandlers.head, asyncErrHandlers.last)
	assert.NotNil(t, asyncErrHandlers.head.handler)
	assert.Nil(t, asyncErrHandlers.head.next)
	assert.Nil(t, asyncErrHandlers.last.next)

	assert.False(t, isErrCfgsFixed)

	FixErrCfgs()

	assert.True(t, isErrCfgsFixed)

	AddSyncErrHandler(func(err Err, occ ErrOccasion) {})
	AddAsyncErrHandler(func(err Err, occ ErrOccasion) {})

	assert.NotNil(t, syncErrHandlers.head)
	assert.NotNil(t, syncErrHandlers.last)
	assert.Equal(t, syncErrHandlers.head, syncErrHandlers.last)
	assert.NotNil(t, syncErrHandlers.head.handler)
	assert.Nil(t, syncErrHandlers.head.next)
	assert.Nil(t, syncErrHandlers.last.next)

	assert.NotNil(t, asyncErrHandlers.head)
	assert.NotNil(t, asyncErrHandlers.last)
	assert.Equal(t, asyncErrHandlers.head, asyncErrHandlers.last)
	assert.NotNil(t, asyncErrHandlers.head.handler)
	assert.Nil(t, asyncErrHandlers.head.next)
	assert.Nil(t, asyncErrHandlers.last.next)
}

func TestNotifyErr_withNoErrHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	NewErr(ReasonForNotification{})

	assert.False(t, isErrCfgsFixed)

	FixErrCfgs()

	assert.True(t, isErrCfgsFixed)

	NewErr(ReasonForNotification{})
}

func TestNotifyErr_withHandlers(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	syncLogs := list.New()
	asyncLogs := list.New()

	AddSyncErrHandler(func(err Err, occ ErrOccasion) {
		syncLogs.PushBack(
			err.ReasonName() + "-1:" + occ.File() + ":" + strconv.Itoa(occ.Line()))
	})
	AddSyncErrHandler(func(err Err, occ ErrOccasion) {
		syncLogs.PushBack(
			err.ReasonName() + "-2:" + occ.File() + ":" + strconv.Itoa(occ.Line()))
	})
	AddAsyncErrHandler(func(err Err, occ ErrOccasion) {
		asyncLogs.PushBack(
			err.ReasonName() + "-3:" + occ.File() + ":" + strconv.Itoa(occ.Line()))
	})

	NewErr(ReasonForNotification{})

	assert.False(t, isErrCfgsFixed)

	assert.Equal(t, syncLogs.Len(), 0)
	assert.Equal(t, asyncLogs.Len(), 0)

	FixErrCfgs()

	NewErr(ReasonForNotification{})

	assert.True(t, isErrCfgsFixed)

	assert.Equal(t, syncLogs.Len(), 2)
	assert.Equal(t, syncLogs.Front().Value,
		"ReasonForNotification-1:notify_test.go:195")
	assert.Equal(t, syncLogs.Front().Next().Value,
		"ReasonForNotification-2:notify_test.go:195")

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, asyncLogs.Len(), 1)
	assert.Equal(t, asyncLogs.Front().Value,
		"ReasonForNotification-3:notify_test.go:195")
}
