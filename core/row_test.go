// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestGetRaw(*testing.T) {
	hdr := SimpleHeader([]string{"one", "two", "two_lower!"})
	var rb RecordBuilder
	rb.Add(SuStr("hello"))
	rb.Add(SuStr("Hello World"))
	rec := rb.Build()
	row := Row{DbRec{Record: rec}}
	assert.This(row.GetRaw(hdr, "one")).Is(Pack(SuStr("hello")))
	assert.This(row.GetRaw(hdr, "two")).Is(Pack(SuStr("Hello World")))
	assert.This(row.GetRaw(hdr, "two_lower!")).Is(Pack(SuStr("hello world")))
}

func BenchmarkRow_GetRawVal(b *testing.B) {
	var row Row
	hdr := SimpleHeader([]string{"foo", "bar"})
	for range b.N {
		row.GetRawVal(hdr, "baz", nil, nil)
	}
}
