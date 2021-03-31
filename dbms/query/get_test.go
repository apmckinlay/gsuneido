// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTableLookup(t *testing.T) {
	pack := func(n int) string {
		return rt.Pack(rt.IntVal(n).(rt.Packable))
	}
	key := func(vals ...int) string {
		if len(vals) == 1 {
			return pack(vals[0])
		}
		var enc ixkey.Encoder
		for _, v := range vals {
			enc.Add(pack(v))
		}
		return enc.String()
	}
	test := func(query, key, expected string) {
		t.Helper()
		q := ParseQuery(query)
		Setup(q, readMode, testTran{})
		row := q.(*Table).Lookup(key)
		assert.T(t).This(fmt.Sprint(row)).Is(expected)
	}
	test("tables", key(123), "[<123>]")
	test("columns", key(12, 34), "[<12, 34>]")
}
