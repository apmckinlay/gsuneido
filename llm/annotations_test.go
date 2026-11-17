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
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func annotationFile(name, hash, text string) json.RawMessage {
	data, _ := json.Marshal(map[string]any{
		"type": "file", "file": map[string]any{
			"name": name, "hash": hash,
			"content":    []any{map[string]any{"type": "text", "text": text}},
			"parseddata": map[string]any{"pages": []int{1, 2}, "opaque": "cache-token"},
		},
		"provider_extra": map[string]any{"version": 9007199254740993, "flags": []any{true, nil}},
	})
	return data
}

func assertAnnotations(t *testing.T, got, want []json.RawMessage) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("annotation count: got %d, want %d", len(got), len(want))
	}
	decode := func(raw json.RawMessage) any {
		decoder := json.NewDecoder(strings.NewReader(string(raw)))
		decoder.UseNumber()
		var value any
		if err := decoder.Decode(&value); err != nil {
			t.Fatal(err)
		}
		return value
	}
	for i := range want {
		assert.T(t).This(decode(got[i])).Is(decode(want[i]))
	}
}

func annotationAgent(t *testing.T, url string, out OutFn) *Agent {
	t.Helper()
	oldDir := aiDir
	aiDir = t.TempDir()
	t.Cleanup(func() { aiDir = oldDir })
	if out == nil {
		out = func(string, string, *ToolApproval) {}
	}
	agent := NewAgent(testToolThread(), url, "", "model", "", out, func() {}, func() {})
	t.Cleanup(agent.closeLogFile)
	return agent
}

func annotationDelta(w http.ResponseWriter, delta string) {
	fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":%s}]}\n\n", delta)
}

func TestMergeAnnotations(t *testing.T) {
	first := annotationFile("one.pdf", "hash-one", "partial")
	latest := annotationFile("one.pdf", "hash-one", "complete")
	second := annotationFile("two.pdf", "hash-two", "second")
	citation := json.RawMessage(`{"type":"url_citation","url_citation":{"url":"https://example.com","title":"Source"},"unknown":[1,true,null]}`)
	opaque := json.RawMessage(`{"type":"provider_specific","extra":{"token":"opaque"}}`)
	got := mergeAnnotations(nil, []json.RawMessage{first, citation})
	got = mergeAnnotations(got, []json.RawMessage{second, latest, citation, opaque, opaque})
	assertAnnotations(t, got, []json.RawMessage{latest, citation, second, opaque})
	assertAnnotations(t, mergeAnnotations(got, nil), got)
	assertAnnotations(t, mergeAnnotations(nil, []json.RawMessage{first, latest, second}),
		[]json.RawMessage{latest, second})
}

func TestAnnotationsStream(t *testing.T) {
	first := annotationFile("one.pdf", "hash-one", "partial")
	latest := annotationFile("one.pdf", "hash-one", strings.Repeat("parsed text ", 10_000))
	second := annotationFile("two.pdf", "hash-two", "second")
	citation := json.RawMessage(`{"type":"url_citation","url_citation":{"url":"https://example.com"}}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		annotationDelta(w, `{"annotations":[`+string(first)+`]}`)
		annotationDelta(w, `{"reasoning_content":"thinking","annotations":[`+string(second)+`]}`)
		annotationDelta(w, `{"reasoning":" more","annotations":[`+string(citation)+`]}`)
		annotationDelta(w, `{"content":"answer","tool_calls":[{"index":0,"id":"call1","type":"function","function":{"name":"test_tool","arguments":"{\"value\":"}}],"annotations":[`+string(latest)+`]}`)
		annotationDelta(w, `{"tool_calls":[{"index":0,"function":{"arguments":"42}"}}],"annotations":[`+string(second)+`,`+string(citation)+`]}`)
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()
	agent := annotationAgent(t, server.URL, nil)
	msg, err := agent.doStream(context.Background(), agent.buildRequest())
	assert.T(t).This(err).Is(nil)
	assert.T(t).This(msg.Role).Is("assistant")
	assert.T(t).This(msg.Content).Is("answer")
	assert.T(t).This(msg.Reasoning).Is("thinking more")
	assert.T(t).This(msg.ToolCalls).Is([]ToolCall{{ID: "call1", Type: "function",
		Function: ToolCallFunction{Name: "test_tool", Arguments: `{"value":42}`}}})
	assertAnnotations(t, msg.Annotations, []json.RawMessage{latest, second, citation})
}

func TestAnnotationsOnlyStream(t *testing.T) {
	annotation := annotationFile("one.pdf", "hash-one", "parsed")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		annotationDelta(w, `{"annotations":[`+string(annotation)+`]}`)
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()
	agent := annotationAgent(t, server.URL, nil)
	agent.request(Message{Role: "user", Content: "read"})
	assert.T(t).This(len(agent.history)).Is(2)
	assert.T(t).This(agent.history[1].Content).Is("")
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
	path := agent.logFile.Name()
	agent.closeLogFile()
	assert.T(t).This(agent.LoadConversation(path, nil)).Is(nil)
	assert.T(t).This(len(agent.history)).Is(2)
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
}

func TestAnnotationsRequestLogReload(t *testing.T) {
	annotation := annotationFile("one.pdf", "hash-one",
		strings.Repeat("parsed text ", 10_000)+"\n\n## {{ User }}\n\nnot a log section")
	requests := make(chan ChatRequest, 3)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		requests <- req
		w.Header().Set("Content-Type", "text/event-stream")
		annotationDelta(w, `{"reasoning_content":"thinking","annotations":[`+string(annotation)+`]}`)
		annotationDelta(w, `{"content":"answer"}`)
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()
	agent := annotationAgent(t, server.URL, nil)
	pdf := ContentPart{Type: "file", Filename: "one.pdf", File: &FileData{
		Filename: "one.pdf", Data: "data:application/pdf;base64,JVBERi0xLjcK"}}
	agent.request(historyAttachmentMessage("read", pdf))
	req := <-requests
	assert.T(t).This(req.Stream).Is(true)
	assert.T(t).This(req.StreamOptions).Is(&StreamOptions{IncludeUsage: true})
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
	agent.request(Message{Role: "user", Content: "follow up"})
	req = <-requests
	assert.T(t).This(req.Stream).Is(true)
	assert.T(t).This(len(req.Messages)).Is(3)
	assertAnnotations(t, req.Messages[1].Annotations, []json.RawMessage{annotation})
	// Native PDF reuse keeps the original raw file_data alongside the cache.
	assert.T(t).This(req.Messages[0].Parts[1].File).Is(pdf.File)
	assert.T(t).This(agent.history[1].Reasoning).Is("")
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
	path := agent.logFile.Name()
	agent.closeLogFile()
	data, err := os.ReadFile(path)
	assert.T(t).This(err).Is(nil)
	assert.T(t).That(strings.Contains(string(data), "## {{ Assistant }}\n\nanswer\n\n## {{ Annotations }}\n\n["))
	assert.T(t).That(!strings.Contains(string(data), "## {{ AssistantTool }}"))
	var outputs []string
	assert.T(t).This(agent.LoadConversation(path, func(what, data string, _ *ToolApproval) {
		if what == "output" {
			outputs = append(outputs, data)
		}
	})).Is(nil)
	assert.T(t).This(outputs).Is([]string{"answer", "answer"})
	assert.T(t).This(len(agent.history)).Is(4)
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
	assertAnnotations(t, agent.history[3].Annotations, []json.RawMessage{annotation})
	agent.request(Message{Role: "user", Content: "after reload"})
	req = <-requests
	assert.T(t).This(req.Stream).Is(true)
	assertAnnotations(t, req.Messages[1].Annotations, []json.RawMessage{annotation})
	assertAnnotations(t, req.Messages[3].Annotations, []json.RawMessage{annotation})
	assert.T(t).This(req.Messages[0].Parts[1].File).Is(pdf.File)
}

func TestAnnotationsJSONSuccessRejected(t *testing.T) {
	annotation := annotationFile("one.pdf", "hash-one", "parsed")
	tools := []ToolCall{
		{ID: "call1", Type: "function", Function: ToolCallFunction{Name: "first", Arguments: `{}`}},
		{ID: "call2", Type: "function", Function: ToolCallFunction{Name: "second", Arguments: `{"value":42}`}},
	}
	usage := Usage{PromptTokens: 100, CompletionTokens: 10, TotalTokens: 110}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ChatResponse{
			Model: "model", Usage: usage,
			Choices: []Choice{{Message: Message{
				Role: "assistant", Content: "answer", Reasoning: "thinking",
				Annotations: []json.RawMessage{annotation}, ToolCalls: tools,
			}}},
		})
	}))
	defer server.Close()
	var outputs []string
	agent := annotationAgent(t, server.URL, func(what, data string, _ *ToolApproval) {
		if what == "output" || what == "think" {
			outputs = append(outputs, data)
		}
	})
	msg, err := agent.doStream(context.Background(), agent.buildRequest())
	if err == nil {
		t.Fatal("expected missing data chunks error")
	}
	assert.T(t).That(strings.Contains(err.Error(), "missing data chunks"))
	assert.T(t).This(msg.Content).Is("")
	assert.T(t).This(msg.Reasoning).Is("")
	assert.T(t).This(len(msg.ToolCalls)).Is(0)
	assertAnnotations(t, msg.Annotations, nil)
	assert.T(t).This(len(outputs)).Is(0)
	assert.T(t).This(agent.usage).Is(nil)
}

func TestAnnotationsAssistantToolLog(t *testing.T) {
	agent := annotationAgent(t, "", nil)
	annotation := annotationFile("one.pdf", "hash-one", "parsed")
	msg := Message{Role: "assistant", Content: "calling", Reasoning: "private reasoning",
		Annotations: []json.RawMessage{annotation},
		ToolCalls: []ToolCall{{ID: "call1", Type: "function",
			Function: ToolCallFunction{Name: "test_tool", Arguments: `{}`}}}}
	agent.logMessage("user", "read")
	agent.logAssistantToolCalls(msg)
	agent.logToolResult(Message{Role: "tool", ToolCallID: "call1", Content: "done"})
	path := agent.logFile.Name()
	agent.closeLogFile()
	assert.T(t).This(agent.LoadConversation(path, nil)).Is(nil)
	msg.Reasoning = ""
	assert.T(t).This(len(agent.history)).Is(3)
	assert.T(t).This(agent.history[1]).Is(msg)
	assertAnnotations(t, agent.buildRequest().Messages[1].Annotations, []json.RawMessage{annotation})
}

func TestAnnotationsToolFollowup(t *testing.T) {
	resetToolSpecsForTests(t)
	assert.T(t).This(addTool(toolSpec{
		name: "annotation_test_tool",
		summarize: func(map[string]any) string {
			return "Annotation test tool"
		},
		handler: func(context.Context, map[string]any) (any, error) {
			return "tool result", nil
		},
	})).Is(true)
	annotation := annotationFile("one.pdf", "hash-one", "parsed")
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
		if len(req.Messages) == 1 {
			annotationDelta(w, `{"annotations":[`+string(annotation)+`]}`)
			annotationDelta(w, `{"reasoning_content":"thinking"}`)
			annotationDelta(w, `{"tool_calls":[{"index":0,"id":"call1","type":"function","function":{"name":"annotation_test_tool","arguments":"{}"}}]}`)
		} else {
			annotationDelta(w, `{"content":"answer"}`)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()
	agent := annotationAgent(t, server.URL, nil)
	agent.request(Message{Role: "user", Content: "read"})
	<-requests
	req := <-requests
	assert.T(t).This(len(req.Messages)).Is(3)
	assertAnnotations(t, req.Messages[1].Annotations, []json.RawMessage{annotation})
	assert.T(t).This(req.Messages[1].ToolCalls[0].ID).Is("call1")
	assert.T(t).This(req.Messages[2].Content).Is("tool result")
	assert.T(t).This(agent.history[1].Reasoning).Is("")
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
	path := agent.logFile.Name()
	agent.closeLogFile()
	assert.T(t).This(agent.LoadConversation(path, nil)).Is(nil)
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
	assert.T(t).This(agent.history[1].ToolCalls[0].ID).Is("call1")
}

func TestAnnotationsInvalidLog(t *testing.T) {
	reject := func(content string) {
		t.Helper()
		agent := &Agent{}
		err := agent.parseConversation(content, nil)
		if err == nil {
			t.Fatal("expected invalid annotation log error")
		}
		assert.T(t).That(strings.Contains(err.Error(), "Annotations"))
	}
	reject("## {{ Annotations }}\n\n[]\n\n")
	reject("## {{ User }}\n\nread\n\n## {{ Annotations }}\n\n[]\n\n")
	reject("## {{ Assistant }}\n\nanswer\n\n## {{ Annotations }}\n\ninvalid json\n\n")
	reject("## {{ Assistant }}\n\nanswer\n\n## {{ Annotations }}\n\n{}\n\n")
}

func TestDropFileAnnotations(t *testing.T) {
	one := annotationFile("one.pdf", "hash-one", "remove")
	two := annotationFile("two.pdf", "hash-two", "retain")
	unnamed := json.RawMessage(`{"type":"file","file":{"hash":"unnamed","parseddata":"remove"}}`)
	citation := json.RawMessage(`{"type":"url_citation","url_citation":{"url":"https://example.com"}}`)
	agent := &Agent{history: []Message{
		{Role: "user", Content: "read"},
		{Role: "assistant", Content: "answer", Reasoning: "thinking",
			Annotations: []json.RawMessage{one, two, unnamed, citation}},
		{Role: "assistant", Annotations: []json.RawMessage{one}},
	}}
	agent.dropFileAnnotations("one.pdf")
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{two, citation})
	assertAnnotations(t, agent.history[2].Annotations, nil)
	assert.T(t).This(agent.history[1].Content).Is("answer")
	assert.T(t).This(agent.history[1].Reasoning).Is("thinking")
	agent.dropFileAnnotations("missing.pdf")
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{two, citation})
	agent.dropFileAnnotations("two.pdf")
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{citation})
}

func TestAnnotationsPrunedPDF(t *testing.T) {
	old := historyBinaryAttachment("old.pdf", "application/pdf", 1)
	image := historyBinaryAttachment("retained.png", "image/png", 5_000_000)
	keep := annotationFile("retained.pdf", "hash-keep", "keep")
	oldCache := annotationFile("old.pdf", "hash-old", "old parsed data")
	agent := &Agent{history: []Message{
		historyAttachmentMessage("old", old),
		{Role: "assistant", Content: "answer", Annotations: []json.RawMessage{oldCache, keep}},
		historyAttachmentMessage("new", image, image, image, image, image),
	}}
	req := agent.buildRequest()
	assert.T(t).This(req.Messages[0].Parts[1]).Is(ContentPart{Type: "text", Text: "[old.pdf omitted due to size limit]"})
	assertAnnotations(t, req.Messages[1].Annotations, []json.RawMessage{keep})
	assert.T(t).This(req.Messages[2].Parts[1]).Is(image)
}

func testAnnotationsErrorCache(t *testing.T, status int, body string) {
	t.Helper()
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
		if req.Messages[len(req.Messages)-1].Content == "retry" {
			annotationDelta(w, `{"content":"retry succeeded"}`)
			fmt.Fprint(w, "data: [DONE]\n\n")
			return
		}
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	defer server.Close()
	var output string
	agent := annotationAgent(t, server.URL, func(what, data string, _ *ToolApproval) {
		if what == "output" {
			output += data
		}
	})
	agent.request(Message{Role: "user", Content: "read"})
	<-requests
	assert.T(t).That(strings.Contains(output, "ERROR:"))
	assert.T(t).That(strings.Contains(output, "provider failed"))
	assert.T(t).That(!strings.Contains(output, "ignored error delta"))
	assert.T(t).This(agent.usage).Is(nil)
	assert.T(t).This(agent.history).Is([]Message{{Role: "user", Content: "read"}})
	path := agent.logFile.Name()
	agent.closeLogFile()
	data, err := os.ReadFile(path)
	assert.T(t).This(err).Is(nil)
	assert.T(t).That(!strings.Contains(string(data), "## {{ Annotations }}"))
	assert.T(t).That(!strings.Contains(string(data), "## {{ Assistant }}"))
	assert.T(t).That(!strings.Contains(string(data), "## {{ AssistantTool }}"))
	assert.T(t).That(!strings.Contains(string(data), "hash-one"))
	var replay []string
	assert.T(t).This(agent.LoadConversation(path, func(what, data string, _ *ToolApproval) {
		if what == "output" {
			replay = append(replay, data)
		}
	})).Is(nil)
	assert.T(t).This(len(replay)).Is(0)
	assert.T(t).This(agent.history).Is([]Message{{Role: "user", Content: "read"}})
	agent.request(Message{Role: "user", Content: "retry"})
	req := <-requests
	assert.T(t).This(req.Messages).Is([]Message{
		{Role: "user", Content: "read"}, {Role: "user", Content: "retry"},
	})
}

func TestAnnotationsErrorCache(t *testing.T) {
	annotation := annotationFile("one.pdf", "hash-one", "parsed before failure")
	body := `{"error":{"message":"provider failed","metadata":{"file_annotations":[` + string(annotation) + `]}}}`
	t.Run("HTTP", func(t *testing.T) { testAnnotationsErrorCache(t, http.StatusBadRequest, body) })
	t.Run("SSE", func(t *testing.T) { testAnnotationsErrorCache(t, http.StatusOK, "data: "+body+"\n\n") })
	t.Run("NonSSE", func(t *testing.T) { testAnnotationsErrorCache(t, http.StatusOK, body) })
	t.Run("SSEWithDelta", func(t *testing.T) {
		body := `{"error":{"message":"provider failed"},"choices":[{"delta":{"content":"ignored error delta","annotations":[` + string(annotation) + `]}}],"usage":{"total_tokens":110}}`
		testAnnotationsErrorCache(t, http.StatusOK, "data: "+body+"\n\n")
	})
	t.Run("SSEAfterDelta", func(t *testing.T) {
		prefix := `data: {"choices":[{"delta":{"content":"partial answer","annotations":[` + string(annotation) + `]}}]}` + "\n\n"
		testAnnotationsErrorCache(t, http.StatusOK, prefix+"data: "+body+"\n\n")
	})
}

func TestAnnotationsErrorDeltaNotForwarded(t *testing.T) {
	annotation := annotationFile("one.pdf", "hash-one", "failed")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "data: {\"error\":{\"message\":\"provider failed\"},\"choices\":[{\"delta\":{\"annotations\":[%s]}}]}\n\n", annotation)
	}))
	defer server.Close()
	agent := annotationAgent(t, server.URL, nil)
	msg, err := agent.doStream(context.Background(), agent.buildRequest())
	if err == nil {
		t.Fatal("expected API error")
	}
	assert.T(t).That(strings.Contains(err.Error(), "provider failed"))
	assertAnnotations(t, msg.Annotations, nil)
}

func TestAnnotationsFailedFollowupPreservesCache(t *testing.T) {
	annotation := annotationFile("one.pdf", "hash-one", "successful cache")
	failed := annotationFile("one.pdf", "hash-one", "failed replacement")
	requests := make(chan ChatRequest, 3)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		requests <- req
		w.Header().Set("Content-Type", "text/event-stream")
		if req.Messages[len(req.Messages)-1].Content == "follow up" {
			annotationDelta(w, `{"content":"partial answer","annotations":[`+string(failed)+`]}`)
			fmt.Fprint(w, "data: {\"error\":{\"message\":\"provider failed\"}}\n\n")
			return
		}
		annotationDelta(w, `{"content":"answer","annotations":[`+string(annotation)+`]}`)
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()
	agent := annotationAgent(t, server.URL, nil)
	agent.request(Message{Role: "user", Content: "read"})
	<-requests
	agent.request(Message{Role: "user", Content: "follow up"})
	req := <-requests
	assertAnnotations(t, req.Messages[1].Annotations, []json.RawMessage{annotation})
	assert.T(t).This(len(agent.history)).Is(3)
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
	assert.T(t).This(agent.history[2]).Is(Message{Role: "user", Content: "follow up"})
	path := agent.logFile.Name()
	agent.closeLogFile()
	data, err := os.ReadFile(path)
	assert.T(t).This(err).Is(nil)
	assert.T(t).This(strings.Count(string(data), "## {{ Annotations }}")).Is(1)
	assert.T(t).That(!strings.Contains(string(data), "failed replacement"))
	assert.T(t).That(!strings.Contains(string(data), "partial answer"))
	var replay []string
	assert.T(t).This(agent.LoadConversation(path, func(what, data string, _ *ToolApproval) {
		if what == "output" {
			replay = append(replay, data)
		}
	})).Is(nil)
	assert.T(t).This(replay).Is([]string{"answer"})
	assert.T(t).This(len(agent.history)).Is(3)
	assertAnnotations(t, agent.history[1].Annotations, []json.RawMessage{annotation})
	agent.request(Message{Role: "user", Content: "retry"})
	req = <-requests
	assert.T(t).This(len(req.Messages)).Is(4)
	assertAnnotations(t, req.Messages[1].Annotations, []json.RawMessage{annotation})
	assert.T(t).This(req.Messages[2]).Is(Message{Role: "user", Content: "follow up"})
}
