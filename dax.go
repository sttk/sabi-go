// Copyright (C) 2023-2025 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

type DaxSrc interface {
	Setup(ag *AsyncGroup) Err
	Close()
}

type daxSrcEntry struct {
	name    string
	ds      DaxSrc
	prev    *daxSrcEntry
	next    *daxSrcEntry
	local   bool
	deleted bool
}

type daxSrcEntryList struct {
	head *daxSrcEntry
	last *daxSrcEntry
}

var (
	isGlobalDaxSrcsFixed  bool = false
	globalDaxSrcEntryList daxSrcEntryList
)

func Uses(name string, ds DaxSrc) {
	if isGlobalDaxSrcsFixed {
		return
	}

	ent := &daxSrcEntry{name: name, ds: ds}

	if globalDaxSrcEntryList.head == nil {
		globalDaxSrcEntryList.head = ent
		globalDaxSrcEntryList.last = ent
	} else {
		ent.prev = globalDaxSrcEntryList.last
		globalDaxSrcEntryList.last.next = ent
		globalDaxSrcEntryList.last = ent
	}
}
