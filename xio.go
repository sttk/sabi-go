// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

// Xio is an interface for a set of inputs/outputs, and requires 2 methods:
// #GetConn which gets a connection to an external data sourc, and #innerMap
// which gets a map to communicate data among multiple inputs/outputs.
type Xio interface {
	GetConn(name string) (Conn, Err)
	InnerMap() map[string]any
}
