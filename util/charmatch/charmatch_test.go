package charmatch

import (
	. "gsuneido/util/hamcrest"
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
	for c := 'a'; c <= 'c'; c++ {
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
	cm := SPACE
	Assert(t).That(cm.CountIn("now is the time"), Equals(3))
}

func TestIndexIn(t *testing.T) {
	cm := AnyOf("abc")
	Assert(t).That(cm.IndexIn("foobar"), Equals(3))
	Assert(t).That(cm.IndexIn("hello"), Equals(-1))
}

func TestPredefined(t *testing.T) {
	cm := LOWER
	Assert(t).That(cm.Match('x'), Equals(true))
	Assert(t).That(cm.Match('X'), Equals(false))
	Assert(t).That(cm.Match('5'), Equals(false))
}
