package errs

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
	isErrCfgFixed = false
}

func TestAddSyncHandler_oneHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddSyncHandler(func(e Err, o ErrOcc) {})

	assert.NotNil(t, syncErrHandlers.head)
	assert.NotNil(t, syncErrHandlers.last)
	assert.Equal(t, syncErrHandlers.head, syncErrHandlers.last)

	assert.Nil(t, syncErrHandlers.last.next)
	assert.Nil(t, syncErrHandlers.head.next)

	assert.NotNil(t, syncErrHandlers.head.handler)
	assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.handler).String(), "func(errs.Err, errs.ErrOcc)")
}

func TestAddSyncHandler_twoHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddSyncHandler(func(e Err, o ErrOcc) {})
	AddSyncHandler(func(e Err, o ErrOcc) {})

	assert.NotNil(t, syncErrHandlers.head)
	assert.NotNil(t, syncErrHandlers.last)
	assert.NotEqual(t, syncErrHandlers.head, syncErrHandlers.last)

	assert.Equal(t, syncErrHandlers.head.next, syncErrHandlers.last)
	assert.Nil(t, syncErrHandlers.last.next)

	assert.NotNil(t, syncErrHandlers.head.handler)
	assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.handler).String(), "func(errs.Err, errs.ErrOcc)")

	assert.NotNil(t, syncErrHandlers.head.next.handler)
	assert.Equal(t, reflect.TypeOf(syncErrHandlers.head.next.handler).String(), "func(errs.Err, errs.ErrOcc)")
}

func TestAddAsyncHandler_oneHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddAsyncHandler(func(e Err, o ErrOcc) {})

	assert.NotNil(t, asyncErrHandlers.head)
	assert.NotNil(t, asyncErrHandlers.last)
	assert.Equal(t, asyncErrHandlers.head, asyncErrHandlers.last)

	assert.Nil(t, asyncErrHandlers.last.next)
	assert.Nil(t, asyncErrHandlers.head.next)

	assert.NotNil(t, asyncErrHandlers.head.handler)
	assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.handler).String(), "func(errs.Err, errs.ErrOcc)")
}

func TestAddAsyncHandler_twoHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddAsyncHandler(func(e Err, o ErrOcc) {})
	AddAsyncHandler(func(e Err, o ErrOcc) {})

	assert.NotNil(t, asyncErrHandlers.head)
	assert.NotNil(t, asyncErrHandlers.last)
	assert.NotEqual(t, asyncErrHandlers.head, asyncErrHandlers.last)

	assert.Equal(t, asyncErrHandlers.head.next, asyncErrHandlers.last)
	assert.Nil(t, asyncErrHandlers.last.next)

	assert.NotNil(t, asyncErrHandlers.head.handler)
	assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.handler).String(), "func(errs.Err, errs.ErrOcc)")

	assert.NotNil(t, asyncErrHandlers.head.next.handler)
	assert.Equal(t, reflect.TypeOf(asyncErrHandlers.head.next.handler).String(), "func(errs.Err, errs.ErrOcc)")
}

func TestFixCfg(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	AddSyncHandler(func(err Err, occ ErrOcc) {})
	AddAsyncHandler(func(err Err, occ ErrOcc) {})

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

	assert.False(t, isErrCfgFixed)

	FixCfg()

	assert.True(t, isErrCfgFixed)

	AddSyncHandler(func(err Err, occ ErrOcc) {})
	AddAsyncHandler(func(err Err, occ ErrOcc) {})

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

	type ReasonForNotification struct{}

	New(ReasonForNotification{})

	assert.False(t, isErrCfgFixed)

	FixCfg()

	assert.True(t, isErrCfgFixed)

	New(ReasonForNotification{})
}

func TestNotifyErr_withErrHandler(t *testing.T) {
	ClearErrHandlers()
	defer ClearErrHandlers()

	syncLogs := list.New()
	asyncLogs := list.New()

	type ReasonForNotification struct{}

	AddSyncHandler(func(e Err, o ErrOcc) {
		syncLogs.PushBack(fmt.Sprintf("%s-1:%s:%d:%s",
			e.ReasonName(), o.File(), o.Line(), o.Time().String()))
	})
	AddSyncHandler(func(e Err, o ErrOcc) {
		syncLogs.PushBack(fmt.Sprintf("%s-2:%s:%d:%s",
			e.ReasonName(), o.File(), o.Line(), o.Time().String()))
	})
	AddAsyncHandler(func(e Err, o ErrOcc) {
		time.Sleep(100 * time.Millisecond)
		asyncLogs.PushBack(fmt.Sprintf("%s-3:%s:%d:%s",
			e.ReasonName(), o.File(), o.Line(), o.Time().String()))
	})
	AddAsyncHandler(func(e Err, o ErrOcc) {
		time.Sleep(10 * time.Millisecond)
		asyncLogs.PushBack(fmt.Sprintf("%s-4:%s:%d:%s",
			e.ReasonName(), o.File(), o.Line(), o.Time().String()))
	})

	assert.False(t, isErrCfgFixed)

	New(ReasonForNotification{})

	assert.Equal(t, syncLogs.Len(), 0)
	assert.Equal(t, asyncLogs.Len(), 0)

	FixCfg()

	assert.True(t, isErrCfgFixed)

	New(ReasonForNotification{})

	assert.Equal(t, syncLogs.Len(), 2)
	log := syncLogs.Front()
	assert.Contains(t, log.Value, "ReasonForNotification-1:notify_test.go:198:")
	log = log.Next()
	assert.Contains(t, log.Value, "ReasonForNotification-2:notify_test.go:198:")
	log = log.Next()
	assert.Nil(t, log)

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, asyncLogs.Len(), 2)
	log = asyncLogs.Front()
	assert.Contains(t, log.Value, "ReasonForNotification-4:notify_test.go:198:")
	log = log.Next()
	assert.Contains(t, log.Value, "ReasonForNotification-3:notify_test.go:198:")
	log = log.Next()
	assert.Nil(t, log)
}
