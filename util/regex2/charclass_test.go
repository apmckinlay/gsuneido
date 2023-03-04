// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCharClass(t *testing.T) {
	test := func(b *builder, c byte, expected bool) {
		t.Helper()
		s := string(b.data)
		if b.isSet {
			pat := Pattern(s)
			assert.T(t).This(matchFullSet(pat, c)).Is(expected)
		} else {
			assert.T(t).This(-1 != strings.IndexByte(s, c)).Is(expected)
		}
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
