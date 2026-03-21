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
)

var _ = addTool(toolSpec{
	name: "suneido_execute",
	description: "Executes Suneido code for its result or side effects.\n" +
		"Use this for calculations, data manipulation, or system commands.\n" +
		"Errors will include the call stack trace",
	params: []stringParam{{name: "code", required: true,
		description: "Suneido code to execute (as the body of a function)"}},
	summarize: func(args map[string]any) string {
		code := argString(args, "code")
		trimmed := strings.TrimSpace(code)
		if strings.Contains(trimmed, "\n") || strings.Contains(trimmed, "\r") {
			return mdSummary("Execute") + "\n" + summarizeCodeBlock(code)
		}
		return mdSummary("Execute", mdInline(trimmed))
	},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		code, err := requireString(args, "code")
		if err != nil {
			return nil, err
		}
		return execTool(code)
	},
})

type execOutput struct {
	Code     string   `json:"code" jsonschema:"the code that was executed"`
	Warnings []string `json:"warnings" jsonschema:"compiler warnings"`
	Results  string   `json:"results" jsonschema:"comma separated list of return values"`
	Print    string   `json:"print,omitempty" jsonschema:"output from Print calls"`
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

	src := funcWrap(code)
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
		results = append(results, displayOrType(th, res))
	} else if len(th.ReturnMulti) > 0 {
		for i := len(th.ReturnMulti) - 1; i >= 0; i-- {
			results = append(results, displayOrType(th, th.ReturnMulti[i]))
		}
	}
	result = execOutput{
		Code:     code,
		Warnings: warnings,
		Results:  strings.Join(results, ", "),
		Print:    printBuf.String(),
	}
	return
}

func funcWrap(code string) string {
	return "function () {\n" + code + "\n}"
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

const maxDisplayLen = 512

func displayOrType(th *core.Thread, val core.Value) string {
	display := core.Display(th, val)
	if len(display) > maxDisplayLen {
		return val.Type().String()
	}
	return display
}

