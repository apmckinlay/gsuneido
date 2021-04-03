// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package assert_test

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestAssert(t *testing.T) {
	assert := assert.T(t)
	assert.That(true)
	assert.True(true)
	assert.False(false)
	assert.Nil(nil)
	assert.This(123).Is(123)
	assert.This(byte(123)).Is(int64(123))
	assert.This(nil).Is(nil)
	assert.This([]byte(nil)).Is(nil)
	assert.This(func() { panic("a test err") }).Panics("test")
	assert.This(" one\t\ntwo ").Like("one\ntwo")
}
