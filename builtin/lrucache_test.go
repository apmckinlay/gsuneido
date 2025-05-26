// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLruCache_concurrent(t *testing.T) {
	const size = 20
	lc := newSuLruCache(size, nil, false)
	a := &SuObject{}
	lc.Insert(Zero, a)
	assert.T(t).This(a.IsConcurrent()).Is(False)
	lc.SetConcurrent()
	assert.T(t).This(a.IsConcurrent()).Is(True)
	b := &SuObject{}
	lc.Insert(One, b)
	assert.T(t).This(b.IsConcurrent()).Is(True)
}
