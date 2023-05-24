// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package opt

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBool(t *testing.T) {
	assert := assert.T(t)
	var b Bool
	assert.That(b.NotSet())
	assert.This(func() { b.Get() }).Panics("opt.Bool")
	assert.This(b.GetOr(true)).Is(true)
	assert.This(b.GetOr(false)).Is(false)
	b.Set(true)
	assert.That(b.IsSet())
	assert.This(b.Get()).Is(true)
	assert.This(b.GetOr(true)).Is(true)
	assert.This(b.GetOr(false)).Is(true)
	b.Set(false)
	assert.This(b.Get()).Is(false)
	assert.This(b.GetOr(true)).Is(false)
	assert.This(b.GetOr(false)).Is(false)
}

func TestInt(t *testing.T) {
	assert := assert.T(t)
	var i Int
	assert.That(i.NotSet())
	assert.This(func() { i.Get() }).Panics("opt.Int")
	assert.This(i.GetOr(123)).Is(123)
	i.Set(123)
	assert.That(i.IsSet())
	assert.This(i.Get()).Is(123)
	assert.This(i.GetOr(123)).Is(123)
	assert.This(i.GetOr(456)).Is(123)
}
