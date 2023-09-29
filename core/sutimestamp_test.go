// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTimestamp_Compare(t *testing.T) {
	d := DateFromLiteral("#20010203").(SuDate)
	ts := SuTimestamp{SuDate: d, extra: 1}
	assert.T(t).True(ts.Compare(d) > 0)
	assert.T(t).True(d.Compare(ts) < 0)
	assert.T(t).True(d.Equal(d))
	assert.T(t).True(ts.Equal(ts))
	assert.T(t).False(ts.Equal(d))
	assert.T(t).False(d.Equal(ts))
	d2 := d.AddMs(1)
	assert.T(t).True(d2.Compare(d) > 0)
	assert.T(t).True(d2.Compare(ts) > 0)
	assert.T(t).True(ts.Compare(d2) < 0)
}
