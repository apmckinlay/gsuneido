// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPermitAuthOnly(t *testing.T) {
	th := &Thread{}
	effective := &Perms{}
	test := func(source, message string) {
		t.Helper()
		fn := compile.Constant("function () { " + source + " }")
		assert.T(t).This(func() { th.Call(fn) }).Panics(message)
		th.Reset()
	}
	th.SetPerms(effective)
	test("Permit.Table('private', 'write')", "Permit.Table can only be called from Auth")
	test("Permit.Schema('update')", "Permit.Schema can only be called from Auth")
	test("Permit.ServerEval('Private')", "Permit.ServerEval can only be called from Auth")
	th.SetPerms(AllPerms)
	test("Permit.Table('private', 'write')", "Permit.Table can only be called from Auth")
	test("Permit.Schema('update')", "Permit.Schema can only be called from Auth")
	test("Permit.ServerEval('Private')", "Permit.ServerEval can only be called from Auth")

	th.SetPerms(effective)
	builder := &Perms{}
	th.SetNewPerms(builder)
	th.Call(compile.Constant(`function () {
		Permit.Table('private', 'read')
		Permit.Schema('create')
		Permit.ServerEval('Private')
		}`))
	assert := assert.T(t)
	assert.True(builder.TableActAllowed("private", PermRead))
	assert.True(builder.SchemaActAllowed(PermCreate))
	assert.True(builder.ServerEvalAllowed("Private"))
	assert.False(effective.TableActAllowed("private", PermRead))
	assert.False(effective.SchemaActAllowed(PermCreate))
	assert.False(effective.ServerEvalAllowed("Private"))
}

func TestPermitGetSet(t *testing.T) {
	assert := assert.T(t)

	th := &Thread{}
	effective := &Perms{}
	th.SetPerms(effective)

	// Get returns the current permissions
	get := compile.Constant(`function () { return Permit.Get() }`)
	assert.That(th.Call(get).(suPerms).perms == effective)

	// Set can only be called from Auth
	set := compile.Constant(`function (p) { Permit.Set(p) }`)
	th.SetNewPerms(nil)
	assert.This(func() { th.Call(set, suPerms{perms: effective}) }).
		Panics("Permit.Set can only be called from Auth")
	th.Reset()

	// Set from Auth replaces the current and new permissions
	target := &Perms{}
	target.AddTable("private", "write")
	th.SetPerms(effective)
	th.SetNewPerms(&Perms{})
	th.Call(set, suPerms{perms: target})
	assert.That(th.Perms() == target)
	assert.That(th.NewPerms() == target)
	assert.True(th.Perms().TableActAllowed("private", PermWrite))
}
