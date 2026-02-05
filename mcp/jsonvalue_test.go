// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"encoding/json"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestJsonValue(t *testing.T) {
	assert := assert.T(t)
	assert.This(jsonValue(core.True, 0)).Is(true)
	assert.This(jsonValue(core.IntVal(123), 0)).Is(json.RawMessage("123"))
	assert.This(jsonValue(core.IntVal(-45), 0)).Is(json.RawMessage("-45"))

	d := core.NewDate(2026, 2, 3, 12, 34, 56, 7)
	assert.This(jsonValue(d, 0)).Is("2026-02-03T12:34:56.007")

	ob := core.SuObjectOf(core.IntVal(1), core.IntVal(2))
	assert.This(jsonValue(ob, 0)).Is([]any{json.RawMessage("1"), json.RawMessage("2")})

	ob2 := &core.SuObject{}
	ob2.Set(core.SuStr("a"), core.IntVal(1))
	ob2.Add(core.IntVal(9))
	v := jsonValue(ob2, 0)
	m, ok := v.(map[string]any)
	assert.That(ok)
	assert.This(m["0"]).Is(json.RawMessage("9"))
	assert.This(m["a"]).Is(json.RawMessage("1"))
	// list members are included in the same map via Iter2 keys
}
