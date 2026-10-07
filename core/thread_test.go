// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestThreadPerms(t *testing.T) {
	assert := assert.T(t)
	parent := &Thread{}
	p := &Perms{}
	p.AddTable("allowed", "read")
	parent.SetPerms(p)
	parent.SetNewPerms(p)
	child := parent.NewChild()
	assert.That(child.Perms() == p)
	assert.That(child.NewPerms() == nil)
	assert.True(child.Perms().TableActAllowed("allowed", PermRead))
	assert.False(child.Perms().TableActAllowed("denied", PermRead))
	parent.Reset()
	assert.That(parent.Perms() == nil)
	assert.That(parent.NewPerms() == nil)
	assert.That(child.Perms() == p)
}
