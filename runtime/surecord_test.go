// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSuRecord(t *testing.T) {
	r := new(SuRecord)
	assert.T(t).This(r.Type()).Is(types.Record)
	assert.T(t).This(r.String()).Is("[]")
	r.Set(SuStr("a"), SuInt(123))
	assert.T(t).This(r.String()).Is("[a: 123]")
}

func TestSuRecord_ReadonlyUnpack(t *testing.T) {
	b := RecordBuilder{}
	b.Add(SuInt(123))
	b.Add(SuStr("foobar"))
	rec := b.Build()
	dbrec := DbRec{Record: rec}
	row := Row{dbrec}

	hdr := NewHeader([][]string{{"num", "str"}}, []string{"num", "str"})
	surec := SuRecordFromRow(row, hdr, nil)

	assert.T(t).This(surec.Get(nil, SuStr("str"))).Is(SuStr("foobar"))
	surec.SetReadOnly()
	assert.T(t).This(surec.Get(nil, SuStr("num"))).Is(SuInt(123))
}
