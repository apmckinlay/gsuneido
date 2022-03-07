// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestJoin(t *testing.T) {
	assert := assert.T(t).This
	assert(Join("", nil)).Is("")
	assert(Join("", []string{})).Is("")
	assert(Join("", []string{"one", "two", "three"})).Is("onetwothree")
	assert(Join(",", []string{"one", "two", "three"})).Is("one,two,three")
	assert(Join(", ", []string{"one", "two", "three"})).Is("one, two, three")
	assert(Join("()", []string{"one", "two", "three"})).Is("(onetwothree)")
	assert(Join("[::]", []string{"one", "two", "three"})).Is("[one::two::three]")
}
