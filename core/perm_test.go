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

	p.AddTable("foo", "rauds")
	assert.This(p.tablePerms["foo"]).Is(PermRead | PermAdd | PermUpdate | PermDelete | PermSchema)

	p.AddTable("bar", "ru")
	assert.This(p.tablePerms["bar"]).Is(PermRead | PermUpdate)

	p.AddTable("baz", "")
	assert.This(p.tablePerms["baz"]).Is(0)

	assert.This(func() {
		p.AddTable("x", "rx")
	}).Panics("invalid table rights")
	assert.This(func() {
		p.AddTable("x", "rar")
	}).Panics("invalid table rights")
	assert.This(func() {
		p.AddTable("x", "z")
	}).Panics("invalid table rights")
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
