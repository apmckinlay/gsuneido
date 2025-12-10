// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build gui

package dbms

import "github.com/apmckinlay/gsuneido/util/assert"

func Server(dbms *DbmsLocal) {
	assert.ShouldNotReachHere()
}

func StopServer() {
	assert.ShouldNotReachHere()
}

func Conns() string {
	panic(assert.ShouldNotReachHere())
}
