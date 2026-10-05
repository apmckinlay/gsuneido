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

	// success - perms are moved to serverConn
	Global.TestDef("Auth", compile.Constant("function (@args) { return true }"))
	th := &Thread{}
	result, perms := auth(th, &SuObject{})
	assert.True(result)
	assert.That(perms != nil)
	assert.That(th.Perms() == nil)
	assert.That(th.NewPerms() == nil)

	// failure - perms are discarded
	Global.TestDef("Auth", compile.Constant("function (@args) { return false }"))
	th = &Thread{}
	result, perms = auth(th, &SuObject{})
	assert.False(result)
	assert.That(perms == nil)
	assert.That(th.Perms() == nil)
	assert.That(th.NewPerms() == nil)

	// throw - perms are discarded and the error propagates
	Global.TestDef("Auth", compile.Constant("function (@args) { Nope() }"))
	th = &Thread{}
	assert.This(func() { auth(th, &SuObject{}) }).Panics("can't find Nope")
	assert.That(th.Perms() == nil)
	assert.That(th.NewPerms() == nil)
}

func TestNewPermsRestore(t *testing.T) {
	assert := assert.T(t)
	defer Global.UnloadAll()
	th := &Thread{}
	test := func(fn Value, success bool) {
		Global.TestDef("Auth", fn)
		result, perms := auth(th, &SuObject{})
		assert.This(result).Is(success)
		assert.That((perms != nil) == success)
		assert.That(th.Perms() == nil)
		assert.That(th.NewPerms() == nil)
	}
	test(compile.Constant("function (@args) { return true }"), true)
	test(compile.Constant("function (@args) { return false }"), false)

	Global.TestDef("Auth", compile.Constant("function (@args) { Nope() }"))
	assert.This(func() { auth(th, &SuObject{}) }).Panics("can't find Nope")
	assert.That(th.Perms() == nil)
	assert.That(th.NewPerms() == nil)
}

func TestStandalonePerms(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	defer db.Close()
	local := NewDbmsLocal(db)
	unauth := Unauth(local)
	prevDbms := StandaloneDbms.Load()
	prevGetDbms := GetDbms
	defer func() {
		StandaloneDbms.Store(prevDbms)
		GetDbms = prevGetDbms
		Global.UnloadAll()
	}()
	StandaloneDbms.Store(unauth)
	GetDbms = func() IDbms { return StandaloneDbms.Load() }
	th := NewThread(nil)
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
