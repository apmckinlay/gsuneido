// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCache(*testing.T) {
	var cache Cache
	data := strings.Fields("a b c d e f g h")
	for _, d := range data {
		p := cache.Get(d)
		assert.This(p.String()).Is(Compile(d).String())
	}
	const n = 100
	for i := 0; i < n; i++ {
		j := rand.Intn(len(data))
		p := cache.Get(data[j])
		assert.This(p.String()).Is(Compile(data[j]).String())
	}
	// assert.T(t).This(CacheGet).Is(n + len(data))
	// assert.T(t).This(CacheHit).Is(n)

	// CacheGet, CacheHit = 0, 0
	data = strings.Fields("a b c d e f g h i j k l m n o")
	for i := 0; i < n; i++ {
		j := rand.Intn(len(data))
		p := cache.Get(data[j])
		assert.This(p.String()).Is(Compile(data[j]).String())
	}
	// assert.T(t).Msg(CacheHit).That(CacheHit > n/2)
}
