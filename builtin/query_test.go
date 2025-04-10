// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestQueryWhere(t *testing.T) {
	names := []Value{SuStr("master_key"), SuStr("flag")}
	as := &ArgSpec{Nargs: 2, Spec: []byte{0, 1}, Names: names}
	args := []Value{DateFromLiteral("19010101"), False}
	assert.T(t).This(queryWhere(as, args)).
		Is("\nwhere master_key is #19010101\nand flag is false")

	names = []Value{SuStr("master_key"), SuStr("flag"), SuStr("sort")}
	as = &ArgSpec{Nargs: 3, Spec: []byte{0, 1, 2}, Names: names}
	args = []Value{DateFromLiteral("19010101"), False, SuStr("sortfield")}
	assert.T(t).This(queryWhere(as, args)).
		Is("\nwhere master_key is #19010101\nand flag is false sort sortfield")
}
