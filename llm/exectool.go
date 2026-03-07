// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/str"
)

var _ = addTool(toolSpec{
	name: "suneido_execute",
	description: "Executes Suneido code for its result or side effects.\n" +
		"Use this for calculations, data manipulation, or system commands.\n" +
		"Note: A single returned object will appear as the first result (e.g., [[1,2]]), while multiple return values appear as separate elements (e.g., [1,2]).",
	params: []stringParam{{name: "code", description: "Suneido code to execute (as the body of a function)", required: true}},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		code, err := requireString(args, "code")
		if err != nil {
			return nil, err
		}
		return execTool(code)
	},
})

type execOutput struct {
	Code     string   `json:"code" jsonschema:"The code that was executed"`
	Warnings []string `json:"warnings" jsonschema:"Compiler warnings"`
	Results  string   `json:"results" jsonschema:"0, 1, or multiple return values as Suneido-format strings"`
}

var _ = addTool(toolSpec{
	name:        "suneido_check_code",
	description: "Checks Suneido code for syntax and compilation errors without executing it. Returns compiler warnings only.",
	params:      []stringParam{{name: "code", description: "Suneido code to check (as the body of a function)", required: true}},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		code, err := requireString(args, "code")
		if err != nil {
			return nil, err
		}
		return checkTool(code)
	},
})

type checkCodeOutput struct {
	Code     string   `json:"code" jsonschema:"The code that was checked"`
	Warnings []string `json:"warnings" jsonschema:"Compiler warnings"`
}

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
	if warnings == nil {
		warnings = []string{}
	}
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

func checkTool(code string) (result checkCodeOutput, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("check error: %v", r)
		}
	}()

	th := core.NewThread(core.MainThread)
	defer th.Close()

	code = strings.TrimSpace(code)
	src := "function () {\n" + code + "\n}"
	_, warnings := compile.Checked(th, src)
	if warnings == nil {
		warnings = []string{}
	}

	result = checkCodeOutput{
		Code:     code,
		Warnings: warnings,
	}
	return
}
