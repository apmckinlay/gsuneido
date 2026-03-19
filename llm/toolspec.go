// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type stringParam struct {
	name        string
	description string
	required    bool
	kind        paramKind
}

type paramKind int

const (
	paramString paramKind = iota
	paramNumber
	paramBool
)

type toolSpec struct {
	name        string
	description string
	params      []stringParam
	summarize   func(args map[string]any) string
	handler     func(ctx context.Context, args map[string]any) (any, error)
}

func (spec toolSpec) inputSchema() map[string]any {
	return inputSchema(spec.params)
}

var toolSpecs []toolSpec

// addTool adds a tool specification to the toolSpecs list.
func addTool(spec toolSpec) bool {
	toolSpecs = append(toolSpecs, spec)
	return true
}

func inputSchema(params []stringParam) map[string]any {
	schema := map[string]any{"type": "object", "properties": map[string]any{}}
	if len(params) == 0 {
		return schema
	}
	props := make(map[string]any, len(params))
	required := make([]string, 0, len(params))
	for _, p := range params {
		t := "string"
		switch p.kind {
		case paramNumber:
			t = "integer"
		case paramBool:
			t = "boolean"
		}
		props[p.name] = map[string]any{
			"type":        t,
			"description": p.description,
		}
		if p.required {
			required = append(required, p.name)
		}
	}
	schema["properties"] = props
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func requireString(args map[string]any, name string) (string, error) {
	s, ok := args[name].(string)
	if !ok {
		return "", errors.New(name + " must be a string")
	}
	return s, nil
}

func optionalString(args map[string]any, name string) string {
	s, _ := args[name].(string)
	return s
}

func optionalInt(args map[string]any, name string, def int) (int, error) {
	val, ok := args[name]
	if !ok || val == nil {
		return def, nil
	}
	switch n := val.(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case float64:
		if n != float64(int(n)) {
			return 0, errors.New(name + " must be an integer")
		}
		return int(n), nil
	case float32:
		if n != float32(int(n)) {
			return 0, errors.New(name + " must be an integer")
		}
		return int(n), nil
	case string:
		if n == "" {
			return def, nil
		}
		parsed, err := strconv.Atoi(n)
		if err != nil {
			return 0, errors.New(name + " must be an integer")
		}
		return parsed, nil
	default:
		return 0, errors.New(name + " must be an integer")
	}
}

func optionalBool(args map[string]any, name string, def bool) (bool, error) {
	val, ok := args[name]
	if !ok || val == nil {
		return def, nil
	}
	switch b := val.(type) {
	case bool:
		return b, nil
	case string:
		if b == "" {
			return def, nil
		}
		switch b {
		case "true", "True", "TRUE", "1":
			return true, nil
		case "false", "False", "FALSE", "0":
			return false, nil
		default:
			return false, errors.New(name + " must be a boolean")
		}
	default:
		return false, errors.New(name + " must be a boolean")
	}
}

func mdSummary(tool string, args ...string) string {
	var b strings.Builder
	b.WriteString("**")
	b.WriteString(tool)
	b.WriteString("**")
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		b.WriteString(" ")
		b.WriteString(arg)
	}
	return b.String()
}

func mdInline(s string) string {
	return "`" + strings.ReplaceAll(s, "`", "\\`") + "`"
}

func mdAny(v any) string {
	switch x := v.(type) {
	case nil:
		return "`null`"
	case string:
		return mdInline(oneLine(x, 40))
	case bool:
		if x {
			return "`true`"
		}
		return "`false`"
	case float64:
		if x == float64(int(x)) {
			return "`" + strconv.Itoa(int(x)) + "`"
		}
		return "`" + fmt.Sprintf("%g", x) + "`"
	default:
		return mdInline(oneLine(fmt.Sprintf("%v", x), 140))
	}
}

const summarizeCodeMaxLines = 10

func summarizeCodeBlock(code string) string {
	code = strings.ReplaceAll(code, "\r\n", "\n")
	code = strings.ReplaceAll(code, "\r", "\n")
	trimmed := strings.TrimRight(code, "\n")
	if trimmed == "" {
		return "```suneido\n```"
	}
	lines := strings.Split(trimmed, "\n")
	if len(lines) > summarizeCodeMaxLines {
		half := summarizeCodeMaxLines / 2
		head := lines[:half]
		tail := lines[len(lines)-half:]
		lines = append(append(head, "..."), tail...)
	}
	return "```suneido\n" + strings.Join(lines, "\n") + "\n```"
}

func oneLine(s string, max int) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\r\n", " ↩ ")
	s = strings.ReplaceAll(s, "\n", " ↩ ")
	if len(s) > max {
		return s[:max-1] + "…"
	}
	return s
}

func argString(args map[string]any, name string) string {
	s, _ := args[name].(string)
	return s
}

func argBool(args map[string]any, name string) (bool, bool) {
	v, ok := args[name]
	if !ok || v == nil {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

func argInt(args map[string]any, name string, def int) int {
	v, ok := args[name]
	if !ok || v == nil {
		return def
	}
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case float32:
		return int(n)
	default:
		return def
	}
}

// argReqStr returns mdInline of a required string argument, or "" if empty.
func argReqStr(args map[string]any, name string) string {
	s := argString(args, name)
	if s == "" {
		return ""
	}
	return mdInline(s)
}

// argOptStr returns "label:`value`" if the string arg is non-empty, else "".
func argOptStr(args map[string]any, name string) string {
	s := argString(args, name)
	if s == "" {
		return ""
	}
	name = strings.ReplaceAll(name, "_", "-")
	return name + ":" + mdInline(s)
}

// argOptInt returns "label:`value`" if the int arg differs from def, else "".
func argOptInt(args map[string]any, name string, def int) string {
	n := argInt(args, name, def)
	if n == def {
		return ""
	}
	name = strings.ReplaceAll(name, "_", "-")
	return name + ":`" + strconv.Itoa(n) + "`"
}

// argOptBool returns label if the bool arg is true, else "".
func argOptBool(args map[string]any, name string) string {
	b, ok := argBool(args, name)
	if !ok || !b {
		return ""
	}
	return strings.ReplaceAll(name, "_", "-")
}
