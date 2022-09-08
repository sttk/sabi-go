// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

import (
	"sync"
)

// Conn is an interface which represents a connection to an external data
// source and requires a methods: #Commit, #Roolback and #Close to work in a
// transaction process.
type Conn interface {
	Commit() Err
	Rollback()
	Close()
}

// ConnCfg is an interface which requires a method: #CreateConn which creates
// a connection to a data source with configuration parameters.
type ConnCfg interface {
	CreateConn() (Conn, Err)
}

var (
	isGlobalConnCfgSealed bool               = false
	globalConnCfgMap      map[string]ConnCfg = make(map[string]ConnCfg)
	globalConnCfgMutex    sync.Mutex
)

// AddGlobalConnCfg registers a global ConnCfg with its name to make enable
// to use ConnCfg in all processes..
func AddGlobalConnCfg(name string, cfg ConnCfg) {
	globalConnCfgMutex.Lock()
	defer globalConnCfgMutex.Unlock()

	if !isGlobalConnCfgSealed {
		globalConnCfgMap[name] = cfg
	}
}

// SealGlobalConnCfgs makes unable to register any further global ConnCfg.
func SealGlobalConnCfgs() {
	isGlobalConnCfgSealed = true
}
