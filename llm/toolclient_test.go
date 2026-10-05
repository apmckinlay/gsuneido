// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestToolClientGetTools(t *testing.T) {
	resetToolSpecsForTests(t)
	_ = addTool(toolSpec{
		name:        "test_tool",
		description: "A test tool",
		summarize: func(args map[string]any) string {
			return mdSummary("Test Tool")
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			return "test result", nil
		},
	})

	client, err := NewToolClient(testToolThread())
	assert.T(t).This(err, nil)
	defer client.Close()

	tools := client.GetTools()
	assert.T(t).True(len(tools) == 1)
	assert.T(t).This(tools[0].Function.Name, "test_tool")
}

func TestToolClientCallTool(t *testing.T) {
	resetToolSpecsForTests(t)
	_ = addTool(toolSpec{
		name:        "echo",
		description: "Echo back the input",
		params: []stringParam{
			{name: "message", kind: paramString},
		},
		summarize: func(args map[string]any) string {
			return mdSummary("Echo", argReqStr(args, "message"))
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			s, _ := args["message"].(string)
			return s, nil
		},
	})

	client, err := NewToolClient(testToolThread())
	assert.T(t).This(err, nil)
	defer client.Close()

	result, err := client.CallTool(context.Background(), "echo", map[string]any{"message": "hello"})
	assert.T(t).This(err, nil)
	assert.T(t).This(result, "hello")
}

func TestToolClientCallToolFromLLM(t *testing.T) {
	resetToolSpecsForTests(t)
	_ = addTool(toolSpec{
		name:        "add",
		description: "Add two numbers",
		summarize: func(args map[string]any) string {
			return mdSummary("Add", mdAny(args["a"]), mdAny(args["b"]))
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			a := int(args["a"].(float64))
			b := int(args["b"].(float64))
			return intStr(a + b), nil
		},
	})

	client, err := NewToolClient(testToolThread())
	assert.T(t).This(err, nil)
	defer client.Close()

	tc := ToolCall{
		ID:   "call_123",
		Type: "function",
		Function: ToolCallFunction{
			Name:      "add",
			Arguments: `{"a": 2, "b": 3}`,
		},
	}

	result, err := client.CallToolFromLLM(context.Background(), tc)
	assert.T(t).This(err, nil)
	assert.T(t).This(result, "5")
}

func TestToolClientCallToolStructuredResult(t *testing.T) {
	resetToolSpecsForTests(t)
	_ = addTool(toolSpec{
		name: "obj",
		summarize: func(args map[string]any) string {
			return mdSummary("Obj")
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			return map[string]any{"ok": true}, nil
		},
	})
	client, err := NewToolClient(testToolThread())
	assert.T(t).This(err, nil)
	defer client.Close()

	result, err := client.CallTool(context.Background(), "obj", nil)
	assert.T(t).This(err, nil)

	var got map[string]any
	err = json.Unmarshal([]byte(result), &got)
	assert.T(t).This(err, nil)
	assert.T(t).This(got["ok"], true)
}

func TestToolClientFormatToolCallForDisplay(t *testing.T) {
	resetToolSpecsForTests(t)
	_ = addTool(toolSpec{
		name: "suneido_demo",
		summarize: func(args map[string]any) string {
			return mdSummary("Demo", "X: "+mdAny(args["x"]))
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			return nil, nil
		},
	})

	client, err := NewToolClient(testToolThread())
	assert.T(t).This(err, nil)
	defer client.Close()

	tc := ToolCall{
		ID:   "call_123",
		Type: "function",
		Function: ToolCallFunction{
			Name:      "suneido_demo",
			Arguments: `{"x": 3}`,
		},
	}

	result, err := client.FormatToolCallForDisplay(tc)
	assert.T(t).This(err, nil)
	assert.T(t).True(strings.Contains(result, "**Demo**"))
	assert.T(t).True(strings.Contains(result, "X: `3`"))
}

func TestToolClientFormatToolCallForDisplayDefault(t *testing.T) {
	resetToolSpecsForTests(t)
	_ = addTool(toolSpec{
		name: "suneido_demo",
		summarize: func(args map[string]any) string {
			return mdSummary("Demo",
				argReqStr(args, "a"),
				argOptBool(args, "b"))
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			return nil, nil
		},
	})

	client, err := NewToolClient(testToolThread())
	assert.T(t).This(err, nil)
	defer client.Close()

	tc := ToolCall{
		ID:   "call_123",
		Type: "function",
		Function: ToolCallFunction{
			Name:      "suneido_demo",
			Arguments: `{"b": true, "a": "x"}`,
		},
	}

	result, err := client.FormatToolCallForDisplay(tc)
	assert.T(t).This(err, nil)
	assert.T(t).True(strings.Contains(result, "**Demo**"))
	assert.T(t).True(strings.Contains(result, "`x`"))
	assert.T(t).True(strings.Contains(result, "b"))
}

func TestToolClientPermissions(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	defer db.Close()
	d := dbms.NewDbmsLocal(db)
	prevGetDbms := core.GetDbms
	defer func() { core.GetDbms = prevGetDbms }()
	core.GetDbms = func() core.IDbms { return d }

	query.DoAdminTest(db, "create public (a) key(a)")
	query.DoAdminTest(db, "create private (a) key(a)")
	perms := &core.Perms{}
	perms.AddTable("public", "read")
	parent := core.NewThread(nil)
	parent.SetPerms(perms)
	parent.SetNewPerms(core.AllPerms)
	client, err := NewToolClient(parent)
	assert.This(err).Is(nil)
	defer client.Close()
	assert.That(client.thread != parent)
	assert.This(client.thread.NewPerms()).Is(nil)

	// Register a thread-aware query primitive without importing builtin,
	// which imports llm. It exercises the execute tool's execution permissions.
	oldRead := core.Global.GetIfPresent("LlmTestRead")
	defer core.Global.TestDef("LlmTestRead", oldRead)
	core.Global.TestDef("LlmTestRead", &core.SuBuiltinRaw{
		Fn: func(th *core.Thread, _ *core.ArgSpec, args []core.Value) core.Value {
			tran := th.Dbms().Transaction(false, th.Perms())
			defer tran.Complete()
			q := tran.Query(core.ToStr(args[0]), nil)
			defer q.Close()
			q.Get(th, core.Next)
			return core.True
		},
		ParamSpec: core.ParamSpec1,
	})
	ctx := context.Background()
	_, err = client.CallTool(ctx, "suneido_query", map[string]any{"query": "public"})
	assert.This(err).Is(nil)
	_, err = client.CallTool(ctx, "suneido_query", map[string]any{"query": "private"})
	assert.That(err != nil)
	assert.That(strings.Contains(err.Error(), "not authorized: private"))
	_, err = client.CallTool(ctx, "suneido_execute", map[string]any{"code": "LlmTestRead('public')"})
	assert.This(err).Is(nil)
	_, err = client.CallTool(ctx, "suneido_execute", map[string]any{"code": "LlmTestRead('private')"})
	assert.That(err != nil)
	assert.That(strings.Contains(err.Error(), "not authorized: private"))
}

func TestToolClientDefaultPermissions(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	defer db.Close()
	d := dbms.NewDbmsLocal(db)
	query.DoAdminTest(db, "create private (a) key(a)")
	oldGetDbms := core.GetDbms
	defer func() { core.GetDbms = oldGetDbms }()
	core.GetDbms = func() core.IDbms { return d }
	parent := core.NewThread(nil)
	parent.SetDbms(d)
	client, err := NewToolClient(parent)
	assert.This(err).Is(nil)
	defer client.Close()
	ctx := context.Background()
	_, err = client.CallTool(ctx, "suneido_query", map[string]any{"query": "private"})
	assert.That(err != nil)
	assert.That(strings.Contains(err.Error(), "not authorized: private"))

	missing, err := NewToolClient(nil)
	assert.This(err).Is(nil)
	defer missing.Close()
	_, err = missing.CallTool(ctx, "suneido_query", map[string]any{"query": "private"})
	assert.That(err != nil)
	assert.That(strings.Contains(err.Error(), "not authorized: private"))
	_, err = queryTool(ctx, "private")
	assert.That(err != nil)
	assert.That(strings.Contains(err.Error(), "not authorized: private"))
}

func resetToolSpecsForTests(t *testing.T) {
	saved := toolSpecs
	t.Cleanup(func() { toolSpecs = saved })
	toolSpecs = nil
}

func intStr(n int) string {
	if n == 0 {
		return "0"
	}
	var neg bool
	if n < 0 {
		neg = true
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
