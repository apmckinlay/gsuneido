// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typechecker

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/typechecker/internal/engine"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// a bug in a pass has to come back as an error the caller can turn into a
// diagnostic, not as a Go panic escaping into the interpreter. A nil class
// stands in for any such bug - RunPipeline dereferences it.
func TestCheckArgumentRecoversFromPassPanic(t *testing.T) {
	a := assert.T(t)
	result, err := checkArgument(SourceEntry{Name: "T", Src: "class { }"},
		"TypeInfer", nil, engine.NewTypeEnv(), engine.NewPassCtx(), nil)
	a.That(result == nil)
	a.That(err != nil)
	a.That(strings.HasPrefix(err.Error(), "T: "))
}

// the result slot still has to hold something the Suneido side can read
func TestFailedResult(t *testing.T) {
	a := assert.T(t)
	a.This(failedResult("TypeAnnotate", "class { }")).Is("class { }")
	ti, ok := failedResult("TypeInfer", "class { }").(TypeInfo)
	a.That(ok)
	a.This(len(ti.Methods)).Is(0)
	a.This(len(ti.Members)).Is(0)
}
