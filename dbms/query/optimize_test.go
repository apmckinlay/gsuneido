// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestOptimize(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query)
		Setup(q, true, testTran{})
		assert.T(t).This(q.String()).Is(expected)
	}
	test("tables", "tables^(table)")
}

var testInfo = map[string]*meta.Info{
	"tables": {Nrows: 99, Size: 9999},
}

func (testTran) GetInfo(table string) *meta.Info {
	return testInfo[table]
}
