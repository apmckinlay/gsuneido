// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"errors"
	"strconv"
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
