// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

func TestRound0(t *testing.T) {
	dn := SuDnum{Dnum: dnum.FromFloat(123.456)}
	x := round(dn, Zero, dnum.HalfUp)
	assert.T(t).This(x).Is(IntVal(123))
	typ := fmt.Sprintf("%T", x)
	assert.T(t).This(typ).Is("*core.smi")
}
