// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/str"
)

var _ = addTool(toolSpec{
	name: "suneido_execute",
	description: "Executes Suneido code for its result or side effects.\n" +
		"Use this for calculations, data manipulation, or system commands.\n" +
		"A single returned object will appear as the first result (e.g., [[1,2]])\n" +
		"multiple return values appear as separate elements (e.g., [1,2]).\n" +
		"Errors will include the call stack trace",
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
	Print    string   `json:"print,omitempty" jsonschema:"Output from Print calls"`
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

// srcPrefix is prepended to user code when wrapping in a function
const srcPrefix = "function () {\n"

// offsetToLine converts a byte offset in src (which includes srcPrefix)
// to a 1-based line number within src,
// then subtracts 1 for the srcPrefix line to give the line in user code.
func offsetToLine(src string, offset int) int {
	if offset > len(src) {
		offset = len(src)
	}
	line := strings.Count(src[:offset], "\n") + 1
	return line - 1 // subtract the srcPrefix line
}

var atOffsetRe = regexp.MustCompile(`@(\d+)`)

// convertPositions replaces @<byteoffset> with @line:<linenum> in s.
// src is the full wrapped source (including srcPrefix).
func convertPositions(s, src string) string {
	return atOffsetRe.ReplaceAllStringFunc(s, func(m string) string {
		n, _ := strconv.Atoi(m[1:])
		return "@line:" + strconv.Itoa(offsetToLine(src, n))
	})
}

func execTool(code string) (result execOutput, err error) {
	var savedCode string
	var th *core.Thread
	defer func() {
		if r := recover(); r != nil {
			msg := convertPositions(fmt.Sprintf("%v", r), savedCode)
			if th != nil {
				if stack := th.StackString(6); stack != "" {
					msg += "\n" + strings.TrimRight(stack, "\n")
				}
			}
			err = fmt.Errorf("execute error: %s", msg)
		}
	}()

	th = core.NewThread(core.MainThread)
	defer th.Close()

	var printBuf strings.Builder
	suneido := core.Suneido.Clone()
	suneido.Set(core.SuStr("Print"), &core.SuBuiltin1{
		Fn: func(s core.Value) core.Value {
			printBuf.WriteString(core.ToStr(s))
			return nil
		},
		BuiltinParams: core.BuiltinParams{ParamSpec: core.ParamSpec1},
	})
	th.Suneido.Store(suneido)

	code = strings.TrimSpace(code)
	src := srcPrefix + code + "\n}"
	savedCode = src
	v, warnings := compile.Checked(th, src)
	if warnings == nil {
		warnings = []string{}
	}
	for i, w := range warnings {
		warnings[i] = convertPositions(w, src)
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
		Print:    printBuf.String(),
	}
	return
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

	code = strings.TrimSpace(code)
	src := srcPrefix + code + "\n}"
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
