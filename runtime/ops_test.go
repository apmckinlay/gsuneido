// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestDiv(t *testing.T) {
	q := OpDiv(SuInt(999), SuInt(3))
	xi, xok := SuIntToInt(q)
	Assert(t).That(xok, Is(true))
	Assert(t).That(xi, Is(333))
	q = OpDiv(SuInt(1), SuInt(3))
	_ = q.(SuDnum)
}

func TestBool(t *testing.T) {
	Assert(t).True(SuBool(true) == True)
	Assert(t).True(SuBool(false) == False)
}
func TestIndex(t *testing.T) {
	Assert(t).That(ToIndex(SuInt(123)), Is(123))
	Assert(t).That(ToIndex(SuDnum{Dnum: dnum.FromInt(123)}), Is(123))
}

var abc = SuStr("abc")
var G Value

func BenchmarkCat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := EmptyStr
		for j := 0; j < 10000; j++ {
			s = OpCat(nil, s, abc)
		}
		G = s
	}
}

func BenchmarkJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ob := NewSuObject()
		for j := 0; j < 10000; j++ {
			ob.Add(abc)
		}
		G = join(ob, EmptyStr)
	}
}

func join(this Value, arg Value) Value {
	ob := ToContainer(this)
	separator := AsStr(arg)
	sb := strings.Builder{}
	sep := ""
	iter := ob.ArgsIter()
	for {
		k, v := iter()
		if k != nil || v == nil {
			break
		}
		sb.WriteString(sep)
		sep = separator
		sb.WriteString(ToStrOrString(v))
	}
	return SuStr(sb.String())
}
