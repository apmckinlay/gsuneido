// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/str"
)

func execTool(code string) (result execOutput, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("execute error: %v", r)
		}
	}()

	th := core.NewThread(core.MainThread)
	defer th.Close()

	code = strings.TrimSpace(code)
	src := "function () {\n" + code + "\n}"
	v, warnings := compile.Checked(th, src)
	fn := v.(*core.SuFunc)

	res := th.Call(fn)
	results := []string{}
	if res != nil {
		results = append(results, res.String())
	} else if len(th.ReturnMulti) > 0 {
		for i := len(th.ReturnMulti) - 1; i >= 0; i-- {
			results = append(results, th.ReturnMulti[i].String())
		}
	}
	result = execOutput{
		Code:     code,
		Warnings: warnings,
		Results:  str.Join("[, ]", results),
	}
	return
}
