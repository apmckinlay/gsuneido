// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestAgentLogging(t *testing.T) {
	// Create a temp directory for .ai
	tmpDir := t.TempDir()
	oldDir := aiDir
	aiDir = tmpDir
	defer func() { aiDir = oldDir }()

	prompt := "You are a helpful assistant."
	agent := NewAgent("", "", "test-model", prompt, func(what, data string, _ *ToolApproval) {})

	// Verify prompt is stored
	assert.T(t).This(agent.prompt).Is(prompt)
	assert.T(t).This(agent.model).Is("test-model")

	// Simulate logging a user message
	agent.logMessage("user", "Hello")
	assert.T(t).True(agent.logFile != nil)

	// Simulate logging an assistant message
	agent.logMessage("assistant", "Hi there!")

	// Simulate logging a tool call via assistant message with tool calls
	agent.logAssistantToolCalls(Message{
		Role:    "assistant",
		Content: "",
		ToolCalls: []ToolCall{{
			ID:       "call1",
			Type:     "function",
			Function: ToolCallFunction{Name: "test_tool", Arguments: `{"arg": "value"}`},
		}},
	})
	agent.logToolResult(Message{
		Role:       "tool",
		Content:    "tool result",
		ToolCallID: "call1",
	})

	agent.closeLogFile()

	// Read the log file
	files, err := os.ReadDir(tmpDir)
	assert.T(t).This(err).Is(nil)
	name := files[0].Name()
	assert.T(t).That(strings.HasPrefix(name, "ai") && strings.HasSuffix(name, ".md"))

	data, err := os.ReadFile(filepath.Join(tmpDir, name))
	assert.T(t).This(err).Is(nil)
	content := string(data)

	// Verify content uses new format
	assert.T(t).True(strings.Contains(content, "## {{ Model }}"))
	assert.T(t).True(strings.Contains(content, "test-model"))
	assert.T(t).True(strings.Contains(content, "## {{ User }}"))
	assert.T(t).True(strings.Contains(content, "Hello"))
	assert.T(t).True(strings.Contains(content, "## {{ Assistant }}"))
	assert.T(t).True(strings.Contains(content, "Hi there!"))
	assert.T(t).True(strings.Contains(content, "## {{ AssistantTool }}"))
	assert.T(t).True(strings.Contains(content, "test_tool"))
	assert.T(t).True(strings.Contains(content, "## {{ ToolResult }}"))
	assert.T(t).True(strings.Contains(content, "tool result"))
}

func TestAgentLogThink(t *testing.T) {
	// Create a temp directory for .ai
	tmpDir := t.TempDir()
	oldDir := aiDir
	aiDir = tmpDir
	defer func() { aiDir = oldDir }()

	agent := NewAgent("", "", "test-model", "", func(what, data string, _ *ToolApproval) {})

	agent.emit("think", "reasoning line")
	agent.emit("output", "done")
	agent.closeLogFile()

	files, err := os.ReadDir(tmpDir)
	assert.T(t).This(err).Is(nil)
	name := files[0].Name()
	assert.T(t).That(strings.HasPrefix(name, "ai") && strings.HasSuffix(name, ".md"))

	data, err := os.ReadFile(filepath.Join(tmpDir, name))
	assert.T(t).This(err).Is(nil)
	content := string(data)

	assert.T(t).True(strings.Contains(content, "## {{ Think }}"))
	assert.T(t).True(strings.Contains(content, "reasoning line"))
}

func TestAgentClearHistory(t *testing.T) {
	agent := NewAgent("", "", "test-model", "test prompt", func(what, data string, _ *ToolApproval) {})

	// Add some history
	agent.history = append(agent.history, Message{Role: "user", Content: "Hello"})
	agent.history = append(agent.history, Message{Role: "assistant", Content: "Hi"})

	// Clear history
	agent.ClearHistory()

	// Verify prompt is retained
	assert.T(t).This(agent.prompt).Is("test prompt")
	// Verify history only contains system prompt
	assert.T(t).This(len(agent.history)).Is(1)
	assert.T(t).This(agent.history[0].Role).Is("system")
	assert.T(t).This(agent.history[0].Content).Is("test prompt")
}

func TestAgentSetModel(t *testing.T) {
	// Create a temp directory for .ai
	tmpDir := t.TempDir()
	oldDir := aiDir
	aiDir = tmpDir
	defer func() { aiDir = oldDir }()

	agent := NewAgent("", "", "initial-model", "test prompt", func(what, data string, _ *ToolApproval) {})

	// Trigger log file creation
	agent.logMessage("user", "Hello")

	// Change model
	agent.SetModel("new-model")
	assert.T(t).This(agent.model).Is("new-model")

	agent.closeLogFile()

	// Read the log file
	files, err := os.ReadDir(tmpDir)
	assert.T(t).This(err).Is(nil)
	data, err := os.ReadFile(filepath.Join(tmpDir, files[0].Name()))
	assert.T(t).This(err).Is(nil)
	content := string(data)

	// Verify both models are logged with new format
	assert.T(t).True(strings.Contains(content, "## {{ Model }}"))
	assert.T(t).True(strings.Contains(content, "initial-model"))
	assert.T(t).True(strings.Contains(content, "new-model"))
}

func TestAgentLoadConversation(t *testing.T) {
	// Create a temp file with conversation content
	tmpFile := filepath.Join(t.TempDir(), "test.md")
	content := `## {{ Model }}

gpt-4

## {{ Prompt }}

You are a helpful assistant.

## {{ User }}

What is 2+2?

## {{ Assistant }}

2+2 equals 4.
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.T(t).This(err).Is(nil)

	// Agent has its own prompt - this should be used, not the file's
	currentPrompt := "Current prompt"
	agent := NewAgent("", "", "test-model", currentPrompt, func(what, data string, _ *ToolApproval) {})
	err = agent.LoadConversation(tmpFile, nil)
	assert.T(t).This(err).Is(nil)

	// Verify model was extracted
	assert.T(t).This(agent.model).Is("gpt-4")

	// Verify current prompt is retained (not replaced by file's prompt)
	assert.T(t).This(agent.prompt).Is(currentPrompt)

	// Verify history was recreated with current prompt
	assert.T(t).This(len(agent.history)).Is(3)
	assert.T(t).This(agent.history[0].Role).Is("system")
	assert.T(t).This(agent.history[0].Content).Is(currentPrompt)
	assert.T(t).This(agent.history[1].Role).Is("user")
	assert.T(t).This(agent.history[1].Content).Is("What is 2+2?")
	assert.T(t).This(agent.history[2].Role).Is("assistant")
	assert.T(t).This(agent.history[2].Content).Is("2+2 equals 4.")

	// ClearHistory should restore the current prompt
	agent.history = append(agent.history, Message{Role: "user", Content: "Another question"})
	agent.ClearHistory()
	assert.T(t).This(len(agent.history)).Is(1)
	assert.T(t).This(agent.history[0].Role).Is("system")
	assert.T(t).This(agent.history[0].Content).Is(currentPrompt)
}

func TestAgentLoadAndResave(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir := aiDir
	aiDir = tmpDir
	defer func() { aiDir = oldDir }()

	// Create original conversation
	originalContent := `## {{ Prompt }}

Original prompt.

## {{ User }}

Question?

## {{ Assistant }}

Answer.
`
	originalFile := filepath.Join(tmpDir, "original.md")
	err := os.WriteFile(originalFile, []byte(originalContent), 0644)
	assert.T(t).This(err).Is(nil)

	// Agent has its own prompt - this should be used, not the file's
	currentPrompt := "Current prompt"
	agent := NewAgent("", "", "test-model", currentPrompt, func(what, data string, _ *ToolApproval) {})
	err = agent.LoadConversation(originalFile, nil)
	assert.T(t).This(err).Is(nil)

	// Log a new message (this creates a new log file)
	agent.logMessage("user", "New question")
	agent.closeLogFile()

	// Find the new log file
	files, err := os.ReadDir(tmpDir)
	assert.T(t).This(err).Is(nil)
	var newFile string
	for _, f := range files {
		if f.Name() != "original.md" {
			newFile = filepath.Join(tmpDir, f.Name())
			break
		}
	}
	assert.T(t).True(newFile != "")

	data, err := os.ReadFile(newFile)
	assert.T(t).This(err).Is(nil)
	newContent := string(data)

	// Verify the current prompt is used (not the file's prompt)
	// Note: prompt is not logged, only user messages
	assert.T(t).True(strings.Contains(newContent, "New question"))
}

func TestAgentLoadConversationWithTools(t *testing.T) {
	assistantMsg := Message{
		Role:    "assistant",
		Content: "",
		ToolCalls: []ToolCall{{
			ID:       "call_abc",
			Type:     "function",
			Function: ToolCallFunction{Name: "myTool", Arguments: `{"x":1}`},
		}},
	}
	toolResultMsg := Message{
		Role:       "tool",
		Content:    "result text",
		ToolCallID: "call_abc",
	}
	marshalMsg := func(m Message) string {
		b, _ := json.Marshal(m)
		return string(b)
	}
	logContent := "## {{ User }}\n\ndo something\n\n" +
		"## {{ AssistantTool }}\n\n" + marshalMsg(assistantMsg) + "\n\n" +
		"## {{ ToolResult }}\n\n" + marshalMsg(toolResultMsg) + "\n\n" +
		"## {{ Assistant }}\n\ndone\n\n"

	tmpFile := filepath.Join(t.TempDir(), "conv.md")
	err := os.WriteFile(tmpFile, []byte(logContent), 0644)
	assert.T(t).This(err).Is(nil)

	agent := NewAgent("", "", "", "", func(what, data string, _ *ToolApproval) {})
	err = agent.LoadConversation(tmpFile, nil)
	assert.T(t).This(err).Is(nil)

	assert.T(t).This(len(agent.history)).Is(4)
	assert.T(t).This(agent.history[0].Role).Is("user")
	assert.T(t).This(agent.history[1].Role).Is("assistant")
	assert.T(t).This(len(agent.history[1].ToolCalls)).Is(1)
	assert.T(t).This(agent.history[1].ToolCalls[0].ID).Is("call_abc")
	assert.T(t).This(agent.history[1].ToolCalls[0].Function.Name).Is("myTool")
	assert.T(t).This(agent.history[2].Role).Is("tool")
	assert.T(t).This(agent.history[2].ToolCallID).Is("call_abc")
	assert.T(t).This(agent.history[2].Content).Is("result text")
	assert.T(t).This(agent.history[3].Role).Is("assistant")
	assert.T(t).This(agent.history[3].Content).Is("done")
}

func TestProcessContentTextThinkTags(t *testing.T) {
	type emission struct {
		what string
		data string
	}
	run := func(chunks []string) ([]emission, string, string) {
		var emitted []emission
		agent := NewAgent("", "", "", "", func(what, data string, _ *ToolApproval) {
			emitted = append(emitted, emission{what, data})
		})
		var content, reasoning strings.Builder
		inThink := false
		for _, chunk := range chunks {
			agent.processContentText(chunk, &content, &reasoning, &inThink)
		}
		return emitted, content.String(), reasoning.String()
	}

	// content before </think> in the same chunk must not be dropped
	emitted, content, reasoning := run([]string{
		"<think>", "for it</think>", "output",
	})
	assert.T(t).This(reasoning).Is("for it")
	assert.T(t).This(content).Is("output")
	_ = emitted

	// content before <think> in the same chunk must not be dropped
	emitted, content, reasoning = run([]string{
		"output before<think>think after",
	})
	assert.T(t).This(content).Is("output before")
	assert.T(t).This(reasoning).Is("think after")
	_ = emitted

	// normal multi-chunk think with no tags
	_, content, reasoning = run([]string{
		"<think>", "hello ", "world", "</think>", "output",
	})
	assert.T(t).This(reasoning).Is("hello world")
	assert.T(t).This(content).Is("output")
}

func TestAgentEmitExecToolResult(t *testing.T) {
	type out struct {
		what string
		data string
	}
	collect := func(result string) []out {
		outs := []out{}
		agent := NewAgent("", "", "", "", func(what, data string, _ *ToolApproval) {
			outs = append(outs, out{what: what, data: data})
		})
		agent.emitExecToolResult(result)
		return outs
	}

	// with results and print
	outs := collect(`{"code":"","warnings":[],"results":"[1,2]","print":"a\nb"}`)
	assert.T(t).This(len(outs)).Is(2)
	assert.T(t).This(outs[0].data).Is("=> [1,2]\n")
	assert.T(t).This(outs[1].data).Is("a\nb\n")

	// no results ([]): still emits => line but without the []
	outs = collect(`{"code":"","warnings":[],"results":"[]"}`)
	assert.T(t).This(len(outs)).Is(1)
	assert.T(t).This(outs[0].data).Is("=> \n")
}

func TestToolApprovalAllowAndDeny(t *testing.T) {
	t.Run("allow", func(t *testing.T) {
		a := newToolApproval()
		go a.Allow("ok")
		d, err := a.Wait(context.Background())
		assert.T(t).This(err).Is(nil)
		assert.T(t).This(d.allow).Is(true)
		assert.T(t).This(d.text).Is("ok")
	})

	t.Run("deny", func(t *testing.T) {
		a := newToolApproval()
		go a.Deny("no")
		d, err := a.Wait(context.Background())
		assert.T(t).This(err).Is(nil)
		assert.T(t).This(d.allow).Is(false)
		assert.T(t).This(d.text).Is("no")
	})
}

func TestAgentWaitForApproval(t *testing.T) {
	t.Run("allow proceeds", func(t *testing.T) {
		agent := NewAgent("", "", "", "", func(what, data string, approval *ToolApproval) {})
		approval := newToolApproval()
		go approval.Allow("approved")
		allowed, text, err := agent.waitForApproval(context.Background(), approval)
		assert.T(t).This(err).Is(nil)
		assert.T(t).This(allowed).Is(true)
		assert.T(t).This(text).Is("approved")
	})

	t.Run("deny returns error", func(t *testing.T) {
		agent := NewAgent("", "", "", "", func(what, data string, approval *ToolApproval) {})
		approval := newToolApproval()
		go approval.Deny("not now")
		allowed, text, err := agent.waitForApproval(context.Background(), approval)
		assert.T(t).This(err).Is(nil)
		assert.T(t).This(allowed).Is(false)
		assert.T(t).This(text).Is("not now")
	})

	t.Run("blocks until decision", func(t *testing.T) {
		agent := NewAgent("", "", "", "", func(what, data string, approval *ToolApproval) {})
		approval := newToolApproval()
		type result struct {
			allowed bool
			text    string
			err     error
		}
		done := make(chan result, 1)
		go func() {
			allowed, text, err := agent.waitForApproval(context.Background(), approval)
			done <- result{allowed: allowed, text: text, err: err}
		}()

		select {
		case r := <-done:
			t.Fatalf("expected wait to block, got %+v", r)
		case <-time.After(10 * time.Millisecond):
		}

		approval.Allow("ok")
		r := <-done
		assert.T(t).This(r.err).Is(nil)
		assert.T(t).This(r.allowed).Is(true)
		assert.T(t).This(r.text).Is("ok")
	})
}

func TestCapTrailingNewlines(t *testing.T) {
	assert := assert.T(t)
	// no trailing newlines - unchanged
	assert.This(capTrailingNewlines("hello", 2)).Is("hello")
	// fewer than cap - unchanged
	assert.This(capTrailingNewlines("hello\n", 2)).Is("hello\n")
	assert.This(capTrailingNewlines("hello\n\n", 2)).Is("hello\n\n")
	// exactly at cap - unchanged
	assert.This(capTrailingNewlines("hello\n\n\n", 3)).Is("hello\n\n\n")
	// over cap - trimmed
	assert.This(capTrailingNewlines("hello\n\n\n", 2)).Is("hello\n\n")
	assert.This(capTrailingNewlines("hello\n\n\n\n", 1)).Is("hello\n")
	// CRLF normalized to LF, count preserved (not inflated to cap)
	assert.This(capTrailingNewlines("hello\r\n", 2)).Is("hello\n")
	assert.This(capTrailingNewlines("hello\r\n\r\n", 2)).Is("hello\n\n")
	// CRLF over cap - trimmed and normalized
	assert.This(capTrailingNewlines("hello\r\n\r\n\r\n", 2)).Is("hello\n\n")
	// bare CR normalized
	assert.This(capTrailingNewlines("hello\r", 2)).Is("hello\n")
	// n=0 strips all trailing newlines
	assert.This(capTrailingNewlines("hello\n\n", 0)).Is("hello")
	assert.This(capTrailingNewlines("hello\r\n", 0)).Is("hello")
	// empty string
	assert.This(capTrailingNewlines("", 2)).Is("")
	// only newlines
	assert.This(capTrailingNewlines("\n\n\n", 2)).Is("\n\n")
}

func TestEnsureTrailingNewlines(t *testing.T) {
	assert := assert.T(t)
	// already correct
	assert.This(ensureTrailingNewlines("hello\n\n", 2)).Is("hello\n\n")
	// too few - adds newlines
	assert.This(ensureTrailingNewlines("hello", 2)).Is("hello\n\n")
	assert.This(ensureTrailingNewlines("hello\n", 2)).Is("hello\n\n")
	// too many - caps
	assert.This(ensureTrailingNewlines("hello\n\n\n", 2)).Is("hello\n\n")
	// CRLF normalized and padded to exactly n
	assert.This(ensureTrailingNewlines("hello\r\n", 2)).Is("hello\n\n")
	// n=0 strips all
	assert.This(ensureTrailingNewlines("hello\n\n", 0)).Is("hello")
}
