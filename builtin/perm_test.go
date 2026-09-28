// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPerm_ServerEval(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	// not during Auth
	assert.This(func() { perm_ServerEval(th, []Value{SuStr("x")}) }).
		Panics("Perm.ServerEval can only be called from Auth")
	// during Auth
	perms := &Perms{}
	th.SetPerms(perms)
	perm_ServerEval(th, []Value{SuStr("F1")})
	assert.True(perms.ServerEvalAllowed("F1"))
	assert.False(perms.ServerEvalAllowed("F2"))
}
