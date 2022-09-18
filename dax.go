// Copyright (C) 2022 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

package sabi

// Dax is an interface for a set of data accesses, and requires two method:
// #GetConn which gets a connection to an external data access, and #InnerMap
// which gets a map to communicate data among multiple data accesses.
type Dax interface {
	GetConn(name string) (Conn, Err)
	InnerMap() map[string]any
}
