// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSuConcat(t *testing.T) {
	assert := assert.T(t)
	a := NewSuConcat()
	b := a
	assert.That(a.buf == b.buf) // same shared buffer
	a = a.Add("hello")
	a = a.Add("world")
	assert.That(a.buf == b.buf) // still same shared buffer
	assert.This(ToStr(b)).Is("")
	b = b.Add("foo")
	assert.That(a.buf != b.buf) // NOT the same shared buffer
}

func TestSuConcat_SetConcurrent(t *testing.T) {
	a := NewSuConcat()
	b := a
	a.SetConcurrent()
	assert.T(t).That(a.buf == b.buf) // same shared buffer
	a = a.Add("hello")
	assert.T(t).That(a.buf != b.buf) // NOT the same shared buffer
}

func TestSuConcat_Equals(t *testing.T) {
	data := []string{"", "a", "ab", "aba", "abc"}
	var str, cat [5]Value
	for i, s := range data {
		str[i] = SuStr(s)
		cat[i] = NewSuConcat().Add(s)
	}
	for i, s := range str {
		for j, c := range cat {
			expected := i == j
			if s.Equal(c) != expected || c.Equal(s) != expected {
				t.Error(s, "vs", c)
			}
		}
	}
}
