// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTokens(t *testing.T) {
	tzr := newInput("Now is (the) time for accumulations.")
	var toks []string
	for tok := tzr.Next(); tok != ""; tok = tzr.Next() {
		toks = append(toks, tok)
	}
	assert.T(t).This(toks).
		Is([]string{"now", "time", "accumul"})
}
