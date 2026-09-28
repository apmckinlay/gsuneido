// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestAuthPerms(t *testing.T) {
	assert := assert.T(t)
	defer Global.UnloadAll()

	// success - perms are moved to serverConn
	Global.TestDef("Auth", compile.Constant("function (@args) { return true }"))
	th := &Thread{}
	result, perms := auth(th, &SuObject{})
	assert.True(result)
	assert.That(perms != nil)
	assert.That(th.Perms() == nil)

	// failure - perms are discarded
	Global.TestDef("Auth", compile.Constant("function (@args) { return false }"))
	th = &Thread{}
	result, perms = auth(th, &SuObject{})
	assert.False(result)
	assert.That(perms == nil)
	assert.That(th.Perms() == nil)

	// throw - perms are discarded and the error propagates
	Global.TestDef("Auth", compile.Constant("function (@args) { Nope() }"))
	th = &Thread{}
	assert.This(func() { auth(th, &SuObject{}) }).Panics("can't find Nope")
	assert.That(th.Perms() == nil)
}
