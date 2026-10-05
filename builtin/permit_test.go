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
