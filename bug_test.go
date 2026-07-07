// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBug(t *testing.T) {
	assert.TestOnlyIndividually(t)

	core.Libload = libload // dependency injection

	openDbms()
	defer db.CloseKeepMapped()

	query := `((((cus leftjoin (ivc extend b2 = i2)) union (cus join ivc)) union ((cus leftjoin ivc) union (cus leftjoin (ivc extend x2 = i2)))) union (((ivc leftjoin cus) union (cus join ivc)) union (((cus leftjoin ivc) rename i4 to y0) union (ivc leftjoin cus))))`
	th := &core.Thread{}
	s := compile.EvalString(th, `QueryStrategy("`+query+`", formatted:)`)
	fmt.Println(core.ToStr(s))
	
	fmt.Println("=== QueryHash ===")
	x := compile.EvalString(th, `QueryHash("`+query+`", details:)`)
	fmt.Println("QueryHash result:", core.ToStr(x))
	
	fmt.Println("=== QueryAltHash ===")
	y := compile.EvalString(th, `QueryAltHash("`+query+`", details:)`)
	fmt.Println("QueryAltHash result:", core.ToStr(y))
	
	assert.This(x).Is(y)
}
