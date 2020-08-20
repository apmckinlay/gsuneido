// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tr

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func Test_makset(t *testing.T) {
	test := func(s, expected string) {
		t.Helper()
		assert.T(t).This(string(Set(""))).Is("")
	}
	test("foo", "foo")
	test("-foo", "-foo")
	test("foo-", "foo-")
	test("-0-9-", "-0123456789-")
}

// ptest support ---------------------------------------------------------------

func TestPtest(t *testing.T) {
	if !ptest.RunFile("tr.test") {
		t.Fail()
	}
}

// pt_tr is a ptest for matching
// usage: "string", "from", "to", "result"
func pt_tr(args []string, _ []bool) bool {
	return Replace(args[0], Set(args[1]), Set(args[2])) == args[3]
}

var _ = ptest.Add("tr", pt_tr)
