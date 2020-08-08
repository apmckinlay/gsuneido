// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/runtime/types"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuRecord(t *testing.T) {
	r := new(SuRecord)
	Assert(t).That(r.Type(), Is(types.Record))
	Assert(t).That(r.String(), Is("[]"))
	r.Set(SuStr("a"), SuInt(123))
	Assert(t).That(r.String(), Is("[a: 123]"))
}

func TestSuRecord_ReadonlyUnpack(t *testing.T) {
	b := RecordBuilder{}
	b.Add(SuInt(123))
	b.Add(SuStr("foobar"))
	rec := b.Build()
	dbrec := DbRec{Record: rec}
	row := Row{dbrec}

	hdr := &Header{Columns: []string{"num", "str"},
		Fields: [][]string{{"num", "str"}}}
	hdr.EnsureMap()

	surec := SuRecordFromRow(row, hdr, nil)

	Assert(t).That(surec.Get(nil, SuStr("str")), Is(SuStr("foobar")))
	surec.SetReadOnly()
	Assert(t).That(surec.Get(nil, SuStr("num")), Is(SuInt(123)))
}
