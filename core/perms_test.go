// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPermsTable(t *testing.T) {
	assert := assert.T(t)
	var p Perms

	p.AddTable("foo", "write")
	assert.This(p.table["foo"]).Is(PermWrite)

	p.AddTable("bar", "read")
	assert.This(p.table["bar"]).Is(PermRead)

	p.AddTable("baz", "")
	assert.This(p.table["baz"]).Is(0)

	assert.This(func() { p.AddTable("x", "update") }).Panics("invalid table rights")
}

func TestPermsServerEval(t *testing.T) {
	assert := assert.T(t)
	var p Perms
	// with no whitelist, nothing is allowed
	assert.False(p.ServerEvalAllowed("F1"))

	p.AddServerEval("Date")
	p.AddServerEval("F1")
	assert.True(p.ServerEvalAllowed("Date"))
	assert.True(p.ServerEvalAllowed("F1"))
	assert.False(p.ServerEvalAllowed("F2"))
	assert.False(p.ServerEvalAllowed("F1x"))

	// allow all with "*"
	var p2 Perms
	p2.AddServerEval("*")
	assert.True(p2.ServerEvalAllowed("anything"))
	assert.True(p2.ServerEvalAllowed("F1"))
}

func TestPermsNil(t *testing.T) {
	assert := assert.T(t)
	test := func(p *Perms) {
		assert.False(p.TableActAllowed("foo", PermRead))
		assert.False(p.TableActAllowed("foo", PermWrite))
		assert.False(p.SchemaActAllowed(PermCreate))
		assert.False(p.SchemaActAllowed(PermUpdate))
		assert.False(p.ServerEvalAllowed("F1"))
	}
	test(nil)
	test(&Perms{})
}

func TestAllPerms(t *testing.T) {
	assert := assert.T(t)
	assert.True(AllPerms.TableActAllowed("foo", PermRead))
	assert.True(AllPerms.TableActAllowed("bar", PermWrite))
	assert.True(AllPerms.SchemaActAllowed(PermCreate))
	assert.True(AllPerms.SchemaActAllowed(PermUpdate))
	assert.True(AllPerms.ServerEvalAllowed("F1"))
	AllPerms.AddTable("foo", "")
	AllPerms.SetSchema("")
	AllPerms.AddServerEval("F2")
	assert.This(len(AllPerms.table)).Is(1)
	assert.This(len(AllPerms.serverEval)).Is(1)
}
