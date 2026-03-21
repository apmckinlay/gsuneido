// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core"
)

var _ = addTool(toolSpec{
	name:        "suneido_check_code",
	description: "Checks Suneido code for syntax and compilation errors without executing it. Returns compiler warnings only.",
	params:      []stringParam{{name: "code", description: "Suneido code to check (as the body of a function)", required: true}},
	summarize: func(args map[string]any) string {
		code := argString(args, "code")
		return mdSummary("Check Code") + "\n" + summarizeCodeBlock(code)
	},
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

func checkTool(code string) (result checkCodeOutput, err error) {
	var savedCode string
	defer func() {
		if r := recover(); r != nil {
			msg := convertPositions(fmt.Sprintf("%v", r), savedCode)
			err = fmt.Errorf("check error: %s", msg)
		}
	}()

	th := core.NewThread(core.MainThread)
	defer th.Close()

	src := funcWrap(code)
	savedCode = src
	_, warnings := compile.Checked(th, src)
	if warnings == nil {
		warnings = []string{}
	}
	for i, w := range warnings {
		warnings[i] = convertPositions(w, src)
	}

	result = checkCodeOutput{
		Code:     code,
		Warnings: warnings,
	}
	return
}
