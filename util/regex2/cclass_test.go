// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCharClass(t *testing.T) {
	test := func(b *cclass, c byte, expected bool) {
		t.Helper()
		pat := Pattern(b[:])
		assert.T(t).This(matchFullSet(pat, c)).Is(expected)
	}
	test(digit, 'x', false)
	test(digit, '0', true)
	test(digit, '5', true)
	test(digit, '9', true)
	test(notWord, 'x', false)
	test(notWord, '_', false)
	test(notWord, '5', false)
	test(notWord, ' ', true)
	test(notWord, '+', true)
	test(space, ' ', true)
	test(space, '\t', true)
	test(space, '\r', true)
	test(space, '\n', true)
	test(space, 'x', false)
	test(space, '0', false)
}

func TestCharClass2(t *testing.T) {
	assert.T(t).This(word.setLen()).Is(16)
	assert.T(t).This(word.listLen()).Is(63)
	assert.T(t).This(digit.listLen()).Is(10)
	assert.T(t).This(digit.list()).
		Is([]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'})
}
