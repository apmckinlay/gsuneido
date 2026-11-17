// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestAgentAttachment(t *testing.T) {
	assert := assert.T(t)
	oldDir := aiDir
	aiDir = t.TempDir()
	defer func() { aiDir = oldDir }()
	path := filepath.Join(t.TempDir(), "document.pdf")
	assert.This(os.WriteFile(path, []byte("%PDF-1.7\n"), 0600)).Is(nil)
	attachment, _, err := ReadAttachment(path)
	assert.This(err).Is(nil)
	assert.This(os.Remove(path)).Is(nil)
	path = filepath.Join(t.TempDir(), "notes.md")
	assert.This(os.WriteFile(path, []byte("# Notes\n\n## {{ User }}\n\nraw-only-file-content\n"), 0600)).Is(nil)
	textAttachment, _, err := ReadAttachment(path)
	assert.This(err).Is(nil)
	assert.This(os.Remove(path)).Is(nil)

	requests := make(chan ChatRequest, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		requests <- req
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"response\"}}]}\n\ndata: [DONE]\n\n")
	}))
	defer server.Close()
	done := make(chan struct{}, 2)
	agent := NewAgent(testToolThread(), server.URL, "", "model", "", func(what, data string, _ *ToolApproval) {
		if strings.HasPrefix(data, "ERROR:") {
			t.Error(data)
		}
	}, func() {}, func() { done <- struct{}{} })
	defer agent.closeLogFile()
	wait := func() {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("attachment request did not complete")
		}
	}
	text := "Explain ## {{ User }} markers\n\nAttached file: document.pdf"
	agent.Input(text, attachment, textAttachment)
	wait()
	req := <-requests
	assert.That(req.Stream)
	assert.This(len(req.Messages)).Is(1)
	assert.This(req.Messages[0].Content).Is(text)
	wireAttachment, wireTextAttachment := attachment, textAttachment
	wireAttachment.Filename, wireTextAttachment.Filename = "", ""
	assert.This(req.Messages[0].Parts).Is([]ContentPart{{Type: "text", Text: text}, wireAttachment, wireTextAttachment})

	logPath := agent.logFile.Name()
	agent.closeLogFile()
	data, err := os.ReadFile(logPath)
	assert.This(err).Is(nil)
	assert.That(strings.Contains(string(data), "## {{ User }}\n\n"+text+"\n\n## {{ Attachment }}\n\n["))
	assert.That(strings.Contains(string(data), `"type":"text","text":"Attached file: notes.md`))
	var display string
	assert.This(agent.LoadConversation(logPath, func(what, data string, _ *ToolApproval) {
		if what == "user" {
			display = data
		}
	})).Is(nil)
	assert.This(display).Is(text)
	assert.That(!strings.Contains(display, "base64"))
	assert.That(!strings.Contains(display, "raw-only-file-content"))
	assert.This(agent.history[0]).Is(Message{Role: "user", Content: text,
		Parts: []ContentPart{{Type: "text", Text: text}, attachment, textAttachment}})

	// The restored attachment is resent on follow-up even without the source file.
	agent.Input("Follow up")
	wait()
	req = <-requests
	assert.That(req.Stream)
	assert.This(len(req.Messages)).Is(3)
	assert.This(req.Messages[0].Parts[1]).Is(wireAttachment)
	assert.This(req.Messages[0].Parts[2]).Is(wireTextAttachment)
	assert.This(req.Messages[2].Content).Is("Follow up")
	assert.This(len(req.Messages[2].Parts)).Is(0)
}

func TestAgentStreamedPDFOutput(t *testing.T) {
	assert := assert.T(t)
	oldDir := aiDir
	aiDir = t.TempDir()
	defer func() { aiDir = oldDir }()
	answer := strings.Repeat("streamed answer\n", 500)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		assert.This(json.NewDecoder(r.Body).Decode(&req)).Is(nil)
		assert.That(req.Stream)
		w.Header().Set("Content-Type", "text/event-stream")
		data, err := json.Marshal(DeltaMessage{Content: "<think>reasoning</think>" + answer})
		assert.This(err).Is(nil)
		fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":%s}]}\n\ndata: [DONE]\n\n", data)
	}))
	defer server.Close()
	var output, thinking strings.Builder
	agent := &Agent{client: NewOpenAIClient(server.URL, ""), outfn: func(what, data string, _ *ToolApproval) {
		if what == "output" {
			output.WriteString(data)
		} else if what == "think" {
			thinking.WriteString(data)
		}
	}}
	defer agent.closeLogFile()
	msg, err := agent.doStream(context.Background(), &ChatRequest{Messages: []Message{
		{Role: "user", Content: "read", Parts: []ContentPart{
			{Type: "text", Text: "read"},
			{Type: "file", File: &FileData{Filename: "file.pdf", Data: "data:application/pdf;base64,eA=="}},
		}},
	}})
	assert.This(err).Is(nil)
	assert.This(msg.Content).Is(answer)
	assert.This(output.String()).Is(answer)
	assert.This(msg.Reasoning).Is("reasoning")
	assert.This(thinking.String()).Is("reasoning")
}

func TestAgentAttachmentSections(t *testing.T) {
	assert := assert.T(t)
	pdf := ContentPart{Type: "file", File: &FileData{
		Filename: "document ## {{ User }}.pdf", Data: "data:application/pdf;base64,cGRm"}}
	image := ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: "data:image/png;base64,aW1hZ2U="}}
	notes := ContentPart{Type: "text", Text: "Attached file: notes.txt\n\ncontents\n"}
	test := func(text string, attachments ...ContentPart) {
		content := "## {{ User }}\n\n" + text + "\n\n"
		for _, attachment := range attachments {
			data, err := json.Marshal([]ContentPart{attachment})
			assert.This(err).Is(nil)
			content += "## {{ Attachment }}\n\n" + string(data) + "\n\n"
		}
		content += "## {{ Assistant }}\n\nreply\n\n## {{ User }}\n\nfollow up\n\n"
		agent := &Agent{}
		var users []string
		assert.This(agent.parseConversation(content, func(what, data string, _ *ToolApproval) {
			if what == "user" {
				users = append(users, data)
			}
		})).Is(nil)
		assert.This(users).Is([]string{text, "follow up"})
		assert.This(len(agent.history)).Is(3)
		assert.This(agent.history[0].Content).Is(text)
		var parts []ContentPart
		if text != "" && len(attachments) > 0 {
			parts = append(parts, ContentPart{Type: "text", Text: text})
		}
		parts = append(parts, attachments...)
		assert.This(agent.history[0].Parts).Is(parts)
		assert.This(agent.history[1].Role).Is("assistant")
		assert.This(agent.history[2]).Is(Message{Role: "user", Content: "follow up"})
	}
	test("plain text")
	test("Explain ## {{ User }} markers", pdf)
	test("", image)
	test("both attachments", pdf, image)
	test("text file", notes)
	test("", notes)
	test("mixed attachments", notes, pdf, image)

	reject := func(content, want string) {
		agent := &Agent{}
		err := agent.parseConversation(content, nil)
		assert.That(err != nil)
		assert.That(strings.Contains(err.Error(), want))
	}
	reject("## {{ Attachment }}\n\n[]\n\n", "missing user message")
	reject("## {{ Assistant }}\n\nreply\n\n## {{ Attachment }}\n\n[]\n\n", "missing user message")
	reject("## {{ User }}\n\nhello\n\n## {{ Attachment }}\n\ninvalid\n\n", "parseConversation Attachment:")
}

func TestAgentLogging(t *testing.T) {
	// Create a temp directory for .ai
	tmpDir := t.TempDir()
	oldDir := aiDir
	aiDir = tmpDir
	defer func() { aiDir = oldDir }()

	prompt := "You are a helpful assistant."
	agent := NewAgent(testToolThread(), "", "", "test-model", prompt, func(what, data string, _ *ToolApproval) {}, func() {}, func() {})

	// Verify prompt is stored
	assert.T(t).This(agent.prompt).Is(prompt)
	assert.T(t).This(agent.model).Is("test-model")

	// Simulate logging a user message
	agent.logMessage("user", "Hello")
	assert.T(t).True(agent.logFile != nil)

	// Simulate logging an assistant message
	agent.logMessage("assistant", "Hi there!")
	agent.logResponseModel("provider/resolved-model")

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
	assert.T(t).True(strings.Contains(content, "## {{ ResponseModel }}"))
	assert.T(t).True(strings.Contains(content, "test-model"))
	assert.T(t).True(strings.Contains(content,
		"## {{ ResponseModel }}\n\nprovider/resolved-model"))
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

	agent := NewAgent(testToolThread(), "", "", "test-model", "", func(what, data string, _ *ToolApproval) {}, func() {}, func() {})

	agent.emit("think", "effort line")
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
	assert.T(t).True(strings.Contains(content, "effort line"))
}

func TestAgentClearHistory(t *testing.T) {
	agent := NewAgent(testToolThread(), "", "", "test-model", "test prompt", func(what, data string, _ *ToolApproval) {}, func() {}, func() {})

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

	agent := NewAgent(testToolThread(), "", "", "initial-model", "test prompt", func(what, data string, _ *ToolApproval) {}, func() {}, func() {})

	// Trigger log file creation
	agent.logMessage("user", "Hello")

	// Change model
	agent.SetModel("new-model")
	assert.T(t).This(agent.model).Is("new-model")
	assert.T(t).This(agent.GetModel()).Is("new-model")

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

func TestAgentEmitEffort(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir := aiDir
	aiDir = tmpDir
	defer func() { aiDir = oldDir }()

	var emitted []string
	agent := NewAgent(testToolThread(), "", "", "test-model", "test prompt",
		func(what, data string, _ *ToolApproval) {
			if what == "effort" {
				emitted = append(emitted, data)
			}
		}, func() {}, func() {})

	agent.emitEffort() // default
	agent.SetEffort(EffortHigh)
	agent.emitEffort()
	agent.SetEffort(EffortLow)
	agent.emitEffort()
	agent.SetEffort(EffortDefault)
	agent.emitEffort()
	agent.closeLogFile()

	assert.T(t).This(emitted).Is([]string{"default", "high", "low", "default"})

	files, err := os.ReadDir(tmpDir)
	assert.T(t).This(err).Is(nil)
	data, err := os.ReadFile(filepath.Join(tmpDir, files[0].Name()))
	assert.T(t).This(err).Is(nil)
	content := string(data)

	assert.T(t).True(strings.Contains(content, "## {{ Effort }}"))
	assert.T(t).True(strings.Contains(content, "high"))
	assert.T(t).True(strings.Contains(content, "low"))
	assert.T(t).True(strings.Contains(content, "## {{ Effort }}\n\ndefault"))
}

func TestAgentSetEffort(t *testing.T) {
	agent := NewAgent(testToolThread(), "", "", "test-model", "", func(what, data string, _ *ToolApproval) {}, func() {}, func() {})

	// default: provider decides
	assert.T(t).This(agent.GetEffort()).Is("")
	req := agent.buildRequest()
	assert.T(t).That(req.Reasoning == nil)
	assert.T(t).That(req.IncludeReasoning)

	// explicit level
	agent.SetEffort(EffortHigh)
	assert.T(t).This(agent.GetEffort()).Is("high")
	req = agent.buildRequest()
	assert.T(t).That(req.Reasoning != nil)
	assert.T(t).This(req.Reasoning.Effort).Is("high")
	assert.T(t).That(!req.IncludeReasoning)

	// "default" maps back to the provider default
	agent.SetEffort(EffortDefault)
	assert.T(t).This(agent.GetEffort()).Is("")
	req = agent.buildRequest()
	assert.T(t).That(req.Reasoning == nil)

	// empty also means provider default
	agent.SetEffort("")
	assert.T(t).This(agent.GetEffort()).Is("")

	// invalid level
	assert.T(t).This(func() { agent.SetEffort("bogus") }).
		Panics("invalid effort level")
}

func TestEffortJSON(t *testing.T) {
	data, err := json.Marshal(ChatRequest{Reasoning: &Reasoning{Effort: EffortMedium}})
	assert.T(t).This(err).Is(nil)
	assert.T(t).True(strings.Contains(string(data), `"reasoning":{"effort":"medium"}`))

	data, err = json.Marshal(ChatRequest{})
	assert.T(t).This(err).Is(nil)
	assert.T(t).True(!strings.Contains(string(data), "reasoning"))
}

func TestAgentLoadConversation(t *testing.T) {
	// Create a temp file with conversation content
	tmpFile := filepath.Join(t.TempDir(), "test.md")
	content := `## {{ Model }}

test-router

## {{ Effort }}

high

## {{ ResponseModel }}

provider/model

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
	agent := NewAgent(testToolThread(), "", "", "test-model", currentPrompt, func(what, data string, _ *ToolApproval) {}, func() {}, func() {})
	err = agent.LoadConversation(tmpFile, nil)
	assert.T(t).This(err).Is(nil)
	assert.T(t).This(agent.GetModel()).Is("test-router")
	assert.T(t).This(agent.GetEffort()).Is("high")
	assert.T(t).This(agent.GetResponseModel()).Is("provider/model")

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

func TestAgentLoadConversationEmitsEffort(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.md")
	content := `## {{ Effort }}

high

## {{ User }}

hi

## {{ Assistant }}

hello
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.T(t).This(err).Is(nil)

	var efforts []string
	outfn := func(what, data string, _ *ToolApproval) {
		if what == "effort" {
			efforts = append(efforts, data)
		}
	}
	agent := NewAgent(testToolThread(), "", "", "", "", outfn, func() {}, func() {})
	err = agent.LoadConversation(tmpFile, outfn)
	assert.T(t).This(err).Is(nil)
	assert.T(t).This(agent.GetEffort()).Is("high")
	assert.T(t).This(efforts).Is([]string{"high"})
}

func TestAgentLoadAndResave(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir := aiDir
	aiDir = tmpDir
	defer func() { aiDir = oldDir }()

	// Create original conversation
	originalContent := `## {{ Prompt }}

Original prompt.

## {{ Effort }}

medium

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
	agent := NewAgent(testToolThread(), "", "", "test-model", currentPrompt, func(what, data string, _ *ToolApproval) {}, func() {}, func() {})
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
	assert.T(t).True(strings.Contains(newContent, "## {{ Effort }}"))
	assert.T(t).True(strings.Contains(newContent, "medium"))
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

	agent := NewAgent(testToolThread(), "", "", "", "", func(what, data string, _ *ToolApproval) {}, func() {}, func() {})
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
		agent := NewAgent(testToolThread(), "", "", "", "", func(what, data string, _ *ToolApproval) {
			emitted = append(emitted, emission{what, data})
		}, func() {}, func() {})
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
		agent := NewAgent(testToolThread(), "", "", "", "", func(what, data string, _ *ToolApproval) {
			outs = append(outs, out{what: what, data: data})
		}, func() {}, func() {})
		agent.emitExecToolResults(result)
		return outs
	}

	// with results and print
	outs := collect(`{"code":"","warnings":[],"results":[{"type":"Number","value":"1"},{"type":"Number","value":"2"}],"print":"a\nb"}`)
	assert.T(t).This(len(outs)).Is(1)
	assert.T(t).This(outs[0].data).Is("```\na\nb\n=> 1, 2\n```\n\n")

	// no results and no print: no output
	outs = collect(`{"code":"","warnings":[],"results":[]}`)
	assert.T(t).This(len(outs)).Is(0)

	// truncated result is marked
	outs = collect(`{"code":"","warnings":[],"results":[{"type":"Text","value":"abc...","is_truncated":true}],"print":""}`)
	assert.T(t).This(outs[0].data).Is("```\n=> abc...\n```\n\n")
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
		agent := NewAgent(testToolThread(), "", "", "", "", func(what, data string, approval *ToolApproval) {}, func() {}, func() {})
		approval := newToolApproval()
		go approval.Allow("approved")
		allowed, text, err := agent.waitForApproval(context.Background(), approval)
		assert.T(t).This(err).Is(nil)
		assert.T(t).This(allowed).Is(true)
		assert.T(t).This(text).Is("approved")
	})

	t.Run("deny returns error", func(t *testing.T) {
		agent := NewAgent(testToolThread(), "", "", "", "", func(what, data string, approval *ToolApproval) {}, func() {}, func() {})
		approval := newToolApproval()
		go approval.Deny("not now")
		allowed, text, err := agent.waitForApproval(context.Background(), approval)
		assert.T(t).This(err).Is(nil)
		assert.T(t).This(allowed).Is(false)
		assert.T(t).This(text).Is("not now")
	})

	t.Run("blocks until decision", func(t *testing.T) {
		agent := NewAgent(testToolThread(), "", "", "", "", func(what, data string, approval *ToolApproval) {}, func() {}, func() {})
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
