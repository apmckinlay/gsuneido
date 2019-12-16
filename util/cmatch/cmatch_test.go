// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package cmatch

import (
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"testing"
)

func TestIs(t *testing.T) {
	cm := Is('x')
	Assert(t).That(cm.Match('x'), Equals(true))
	Assert(t).That(cm.Match('y'), Equals(false))
}

func TestAnyOf(t *testing.T) {
	cm := AnyOf("abc")
	Assert(t).That(cm.Match('b'), Equals(true))
	Assert(t).That(cm.Match('x'), Equals(false))
}

func TestInRange(t *testing.T) {
	cm := InRange('a', 'c')
	for c := byte('a'); c <= 'c'; c++ {
		Assert(t).That(cm.Match(c), Equals(true))
	}
	Assert(t).That(cm.Match('x'), Equals(false))
}

func TestNegate(t *testing.T) {
	cm := Is('x').Negate()
	Assert(t).That(cm.Match('x'), Equals(false))
	Assert(t).That(cm.Match('y'), Equals(true))
}

func TestOr(t *testing.T) {
	cm := Is('x').Or(Is('y'))
	Assert(t).That(cm.Match('x'), Equals(true))
	Assert(t).That(cm.Match('y'), Equals(true))
	Assert(t).That(cm.Match('z'), Equals(false))
}

func TestCountIn(t *testing.T) {
	cm := AnyOf(" \t\r\n")
	Assert(t).That(cm.CountIn("now is the time"), Equals(3))
}

func TestIndexIn(t *testing.T) {
	cm := AnyOf("abc")
	Assert(t).That(cm.IndexIn("foobar"), Equals(3))
	Assert(t).That(cm.IndexIn("hello"), Equals(-1))
}

func TestPredefined(t *testing.T) {
	cm := InRange('a', 'z')
	Assert(t).That(cm.Match('x'), Equals(true))
	Assert(t).That(cm.Match('X'), Equals(false))
	Assert(t).That(cm.Match('5'), Equals(false))
}

func TestTrim(t *testing.T) {
	cm := Is(' ')
	Assert(t).That(cm.Trim(""), Equals(""))
	Assert(t).That(cm.Trim(" "), Equals(""))
	Assert(t).That(cm.Trim("    "), Equals(""))
	Assert(t).That(cm.Trim("hello"), Equals("hello"))
	Assert(t).That(cm.Trim("  hello"), Equals("hello"))
	Assert(t).That(cm.Trim("hello "), Equals("hello"))
	Assert(t).That(cm.Trim(" hello  "), Equals("hello"))
}

func TestTrimLeft(t *testing.T) {
	cm := Is(' ')
	Assert(t).That(cm.TrimLeft(""), Equals(""))
	Assert(t).That(cm.TrimLeft(" "), Equals(""))
	Assert(t).That(cm.TrimLeft("    "), Equals(""))
	Assert(t).That(cm.TrimLeft("hello"), Equals("hello"))
	Assert(t).That(cm.TrimLeft("  hello"), Equals("hello"))
	Assert(t).That(cm.TrimLeft("hello "), Equals("hello "))
	Assert(t).That(cm.TrimLeft(" hello  "), Equals("hello  "))
}
