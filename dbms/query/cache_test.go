// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCache(*testing.T) {
	var c cache
	assert.This(int(c.cacheGet(nil))).Is(-1)
	ix1 := []string{"one"}
	assert.This(int(c.cacheGet(ix1))).Is(-1)
	c.cacheAdd(ix1, 123)
	assert.This(int(c.cacheGet(ix1))).Is(123)
	ix2 := []string{"one", "two"}
	assert.This(int(c.cacheGet(ix2))).Is(-1)
	c.cacheAdd(ix2, 456)
	assert.This(int(c.cacheGet(ix1))).Is(123)
	assert.This(int(c.cacheGet(ix2))).Is(456)
}
