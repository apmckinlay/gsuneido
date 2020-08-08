// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestBuilder(t *testing.T) {
	var b RecordBuilder
	rec := b.Build()
	Assert(t).That([]byte(rec), Is([]byte{0}))
	b.AddRaw("one")
	rec = b.Build()
	Assert(t).That([]byte(rec), Is([]byte{type8 << 6, 1, 7, 4, 'o', 'n', 'e'}))
	Assert(t).That(rec.GetRaw(0), Is("one"))

	b = RecordBuilder{}
	b.Add(SuInt(123))
	b.Add(SuStr("foobar"))

	rec = b.Build()
	Assert(t).That(rec.mode(), Is(type8))
	Assert(t).That(rec.Count(), Is(2))
	Assert(t).That(rec.GetVal(0), Is(SuInt(123)))
	Assert(t).That(rec.GetVal(1), Is(SuStr("foobar")))

	s := strings.Repeat("helloworld", 30)
	b.AddRaw(s)
	rec = b.Build()
	Assert(t).That(rec.mode(), Is(type16))
	Assert(t).That(rec.GetRaw(2), Is(s))

}

func TestLength(t *testing.T) {
	Assert(t).That(tblength(0, 0), Is(1))
	Assert(t).That(tblength(1, 1), Is(5))
	Assert(t).That(tblength(1, 200), Is(204))
	Assert(t).That(tblength(1, 248), Is(252))

	Assert(t).That(tblength(1, 252), Is(258))
	Assert(t).That(tblength(1, 300), Is(306))

	Assert(t).That(tblength(1, 0x10000), Is(0x1000a))
}
