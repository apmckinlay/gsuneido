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

	query := `
				(cus
			join by(ck)
					(ivc
				union /*NOT DISJOINT*/
						(ivc
					union /*NOT DISJOINT*/
						(ivc
						where ik is '3'))))
		union /*NOT DISJOINT*/
					((cus
				join by(ck)
					ivc)
			union /*NOT DISJOINT*/
					((cus
					where c2 <= '3')
				join by(ck)
					ivc))`
	th := &core.Thread{}
	s := compile.EvalString(th, `QueryStrategy("`+query+`", formatted:)`)
	fmt.Println(core.ToStr(s))
	compile.EvalString(th, `QueryHash("`+query+`", details:)`)
}
