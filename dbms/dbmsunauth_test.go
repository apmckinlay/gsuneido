// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNewPerms(t *testing.T) {
	assert := assert.T(t)
	defer Global.UnloadAll()

	// failure - perms are discarded
	Global.TestDef("Auth", compile.Constant("function (@args) { return false }"))
	th := &Thread{}
	dbms := &DbmsUnauth{dbms: &DbmsLocal{}}
	result := dbms.Auth(th, &SuObject{})
	assert.False(result)
	assert.That(th.Perms() == nil)
	assert.That(th.NewPerms() == nil)

	// throw - perms are discarded and the error propagates
	Global.TestDef("Auth", compile.Constant("function (@args) { Nope() }"))
	th = &Thread{}
	assert.This(func() { dbms.Auth(th, &SuObject{}) }).Panics("can't find Nope")
	assert.That(th.Perms() == nil)
	assert.That(th.NewPerms() == nil)

	// success - perms are kept on the thread
	Global.TestDef("Auth", compile.Constant("function (@args) { return true }"))
	result = dbms.Auth(th, &SuObject{})
	assert.True(result)
	assert.That(th.Perms() != nil)
	assert.That(th.NewPerms() == nil)
}

func TestStandalonePerms(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	defer db.Close()
	local := NewDbmsLocal(db)
	unauth := Unauth(local)
	defer Global.UnloadAll()
	th := &Thread{}
	th.SetDbms(unauth)
	assert.That(th.Dbms() == unauth)
	assert.That(th.Perms() == nil)
	Global.TestDef("Auth", &SuBuiltinRaw{
		Fn: func(th *Thread, _ *ArgSpec, _ []Value) Value {
			assert.That(th.Perms() == th.NewPerms())
			args := SuObjectOf(SuStr("tables sort table"))
			assert.This(func() { th.Dbms().Get(th, args, Next) }).Panics("not authorized: tables")
			th.NewPerms().AddTable("tables", "read")
			th.Dbms().Get(th, args, Next)
			return True
		}, ParamSpec: ParamSpecAt})
	assert.True(th.Dbms().Auth(th, &SuObject{}))
	assert.That(th.Dbms() == local)
	assert.True(th.Perms().TableActAllowed("tables", PermRead))
	assert.False(th.Perms().TableActAllowed("private", PermRead))
	th.Reset()
	assert.That(th.NewPerms() == nil)
}
