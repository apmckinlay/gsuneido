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
		assert.T(t).This(string(New(s))).Is(expected)
	}
	test("", "")
	test("foo", "foo")
	test("^foo", "^foo")
	test("-foo", "-foo")
	test("foo-", "foo-")
	test("m-p", "mnop")
	test("-0-9-", "-0123456789-")
	test("\xfa-\xff", "\xfa\xfb\xfc\xfd\xfe\xff")
	test("z-a", "")
}

func Fuzz_makset(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		New(s)
	})
}

// to run: go test -fuzz=Fuzz_makset -run=Fuzz_makset

func FuzzReplace(f *testing.F) {
	f.Fuzz(func(t *testing.T, s1, s2, s3 string) {
		Replace(s1, New(s2), New(s3))
	})
}

// to run: go test -fuzz=FuzzReplace -run=FuzzReplace

// ptest support ---------------------------------------------------------------

func TestPtest(t *testing.T) {
	if !ptest.RunFile("tr.test") {
		t.Fail()
	}
}

// pt_tr is a ptest for matching
// usage: "string", "from", "to", "result"
func ptTr(args []string, _ []bool) bool {
	return Replace(args[0], New(args[1]), New(args[2])) == args[3]
}

var _ = ptest.Add("tr", ptTr)
