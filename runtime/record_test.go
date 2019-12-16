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
	Assert(t).That([]byte(rec), Equals([]byte{0}))
	b.AddRaw("one")
	rec = b.Build()
	Assert(t).That([]byte(rec), Equals([]byte{type8 << 6, 1, 7, 4, 'o', 'n', 'e'}))
	Assert(t).That(rec.GetRaw(0), Equals("one"))

	b = RecordBuilder{}
	b.Add(SuInt(123))
	b.Add(SuStr("foobar"))

	rec = b.Build()
	Assert(t).That(rec.mode(), Equals(type8))
	Assert(t).That(rec.Count(), Equals(2))
	Assert(t).That(rec.GetVal(0), Equals(SuInt(123)))
	Assert(t).That(rec.GetVal(1), Equals(SuStr("foobar")))

	s := strings.Repeat("helloworld", 30)
	b.AddRaw(s)
	rec = b.Build()
	Assert(t).That(rec.mode(), Equals(type16))
	Assert(t).That(rec.GetRaw(2), Equals(s))

}

func TestLength(t *testing.T) {
	Assert(t).That(tblength(0, 0), Equals(1))
	Assert(t).That(tblength(1, 1), Equals(5))
	Assert(t).That(tblength(1, 200), Equals(204))
	Assert(t).That(tblength(1, 248), Equals(252))

	Assert(t).That(tblength(1, 252), Equals(258))
	Assert(t).That(tblength(1, 300), Equals(306))

	Assert(t).That(tblength(1, 0x10000), Equals(0x1000a))
}
