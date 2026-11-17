// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func historyTextAttachment(name, body string) ContentPart {
	return ContentPart{Type: "text", Filename: name,
		Text: "Attached file: " + name + "\n\n" + body}
}

func historyBinaryAttachment(name, mime string, size int) ContentPart {
	// Encode repeated 'x' bytes directly, without allocating a second raw payload.
	data := "data:" + mime + ";base64," + strings.Repeat("eHh4", size/3)
	switch size % 3 {
	case 1:
		data += "eA=="
	case 2:
		data += "eHg="
	}
	if mime == "application/pdf" {
		return ContentPart{Type: "file", Filename: name,
			File: &FileData{Filename: name, Data: data}}
	}
	return ContentPart{Type: "image_url", Filename: name,
		ImageURL: &ImageURL{URL: data}}
}

func historyAttachmentMessage(prompt string, parts ...ContentPart) Message {
	msg := Message{Role: "user", Content: prompt}
	if prompt != "" {
		msg.Parts = append(msg.Parts, ContentPart{Type: "text", Text: prompt})
	}
	msg.Parts = append(msg.Parts, parts...)
	return msg
}

func TestContentPartAttachmentSize(t *testing.T) {
	assert := assert.T(t)
	test := func(part ContentPart, size int) {
		assert.This(part.attachmentSize()).Is(size)
	}
	// Text counts UTF-8 bytes, excluding the filename header, not characters.
	test(historyTextAttachment("notes.txt", "a\xc3\xa9\n\nbody"), 9)
	test(historyTextAttachment("empty.txt", ""), 0)
	test(ContentPart{Type: "text", Text: "ordinary user prompt"}, 0)
	test(ContentPart{Type: "text", Text: "[notes.txt omitted due to size limit]"}, 0)
	test(historyBinaryAttachment("empty.pdf", "application/pdf", 0), 0)
	test(historyBinaryAttachment("one.pdf", "application/pdf", 1), 1)
	test(historyBinaryAttachment("two.pdf", "application/pdf", 2), 2)
	test(historyBinaryAttachment("three.pdf", "application/pdf", 3), 3)
	test(historyBinaryAttachment("one.png", "image/png", 1), 1)
	test(historyBinaryAttachment("two.png", "image/png", 2), 2)
	test(historyBinaryAttachment("three.png", "image/png", 3), 3)
}

func TestAttachmentHistoryExactThreshold(t *testing.T) {
	assert := assert.T(t)
	assert.This(MaxTotalAttachmentSize).Is(25_000_000)
	pdf := historyBinaryAttachment("document.pdf", "application/pdf", 10_000_000)
	image := historyBinaryAttachment("picture.png", "image/png", 10_000_000)
	text := historyTextAttachment("notes.txt", strings.Repeat("x", 1_000_000))
	agent := &Agent{history: []Message{
		historyAttachmentMessage("read files", pdf, text, text),
		historyAttachmentMessage("", image, text, text, text),
	}}
	agent.pruneAttachments()
	assert.This(agent.history[0]).Is(historyAttachmentMessage("read files", pdf, text, text))
	assert.This(agent.history[1]).Is(historyAttachmentMessage("", image, text, text, text))
}

func TestAttachmentHistoryOldestFirst(t *testing.T) {
	assert := assert.T(t)
	image := historyBinaryAttachment("picture.png", "image/png", 5_000_000)
	old := historyTextAttachment("old.txt", "x")
	// A prompt can resemble an attachment header without becoming an attachment.
	prompt := "Attached file: not-an-attachment.txt\n\n" + strings.Repeat("p", 1_000_000)
	system := Message{Role: "system", Content: "system instructions", Parts: []ContentPart{image}}
	assistant := Message{Role: "assistant", Content: "calling tool", Reasoning: "reasoning",
		ToolCalls: []ToolCall{{ID: "call1", Type: "function",
			Function: ToolCallFunction{Name: "test_tool", Arguments: `{}`}}}}
	tool := Message{Role: "tool", Content: "tool output", ToolCallID: "call1",
		Parts: []ContentPart{image}}
	agent := &Agent{history: []Message{
		system,
		historyAttachmentMessage("", old, image),
		assistant,
		tool,
		historyAttachmentMessage(prompt, image, image, image, image),
	}}
	agent.pruneAttachments()
	marker := ContentPart{Type: "text", Text: "[old.txt omitted due to size limit]"}
	assert.This(agent.history[0]).Is(system)
	assert.This(agent.history[1]).Is(historyAttachmentMessage("", marker, image))
	assert.This(agent.history[2]).Is(assistant)
	assert.This(agent.history[3]).Is(tool)
	assert.This(agent.history[4]).Is(historyAttachmentMessage(prompt, image, image, image, image))

	// Existing omission markers must neither consume the budget nor get rewritten.
	agent.pruneAttachments()
	assert.This(agent.history[1]).Is(historyAttachmentMessage("", marker, image))
	assert.This(agent.history[4]).Is(historyAttachmentMessage(prompt, image, image, image, image))
	agent.history = append(agent.history, historyAttachmentMessage("", old))
	agent.pruneAttachments()
	assert.This(agent.history[1]).Is(historyAttachmentMessage("", marker,
		ContentPart{Type: "text", Text: "[picture.png omitted due to size limit]"}))
	assert.This(agent.history[4]).Is(historyAttachmentMessage(prompt, image, image, image, image))
	assert.This(agent.history[5]).Is(historyAttachmentMessage("", old))
}

func TestAttachmentHistoryMultipleEvictions(t *testing.T) {
	assert := assert.T(t)
	image := historyBinaryAttachment("retained.png", "image/png", 5_000_000)
	first, second := image, image
	first.Filename = "first.png"
	second.Filename = "second.png"
	agent := &Agent{history: []Message{
		historyAttachmentMessage("keep this prompt", first, second, image, image),
		historyAttachmentMessage("latest", image, image, image),
	}}
	agent.pruneAttachments()
	assert.This(agent.history[0]).Is(historyAttachmentMessage("keep this prompt",
		ContentPart{Type: "text", Text: "[first.png omitted due to size limit]"},
		ContentPart{Type: "text", Text: "[second.png omitted due to size limit]"},
		image, image))
	assert.This(agent.history[1]).Is(historyAttachmentMessage("latest", image, image, image))
}

func TestAttachmentHistoryBuildRequestAndClear(t *testing.T) {
	assert := assert.T(t)
	image := historyBinaryAttachment("picture.png", "image/png", 5_000_000)
	pdf := historyBinaryAttachment("old-document.pdf", "application/pdf", 5_000_000)
	agent := &Agent{client: &OpenAIClient{}, model: "test-model", prompt: "system prompt"}
	agent.resetHistory()
	agent.history = append(agent.history,
		historyAttachmentMessage("old", pdf),
		historyAttachmentMessage("new", image, image, image, image, image))
	req := agent.buildRequest()
	assert.This(req.Model).Is("test-model")
	assert.This(req.Messages).Is(agent.history)
	assert.This(req.Messages[1]).Is(historyAttachmentMessage("old",
		ContentPart{Type: "text", Text: "[old-document.pdf omitted due to size limit]"}))
	assert.This(req.Messages[2]).Is(historyAttachmentMessage("new", image, image, image, image, image))

	agent.ClearHistory()
	assert.This(agent.history).Is([]Message{{Role: "system", Content: "system prompt"}})
	agent.history = append(agent.history,
		historyAttachmentMessage("fresh", image, image, image, image, image))
	req = agent.buildRequest()
	assert.This(len(req.Messages)).Is(2)
	assert.This(req.Messages[1]).Is(historyAttachmentMessage("fresh", image, image, image, image, image))
}

func TestAttachmentHistoryFilenameMetadata(t *testing.T) {
	assert := assert.T(t)
	dir := t.TempDir()
	test := func(name, body string) {
		path := filepath.Join(dir, name)
		assert.This(os.WriteFile(path, []byte(body), 0600)).Is(nil)
		part, size, err := ReadAttachment(path)
		if err != nil {
			t.Fatal(err)
		}
		assert.This(part.Filename).Is(name)
		assert.This(part.attachmentSize()).Is(size)
		wire, err := json.Marshal(historyAttachmentMessage("read", part))
		assert.This(err).Is(nil)
		var decoded struct {
			Content []map[string]json.RawMessage `json:"content"`
		}
		assert.This(json.Unmarshal(wire, &decoded)).Is(nil)
		assert.This(len(decoded.Content)).Is(2)
		_, hasFilename := decoded.Content[1]["filename"]
		_, hasExportedFilename := decoded.Content[1]["Filename"]
		assert.That(!hasFilename && !hasExportedFilename)
		if part.File != nil {
			var file FileData
			assert.This(json.Unmarshal(decoded.Content[1]["file"], &file)).Is(nil)
			assert.This(file.Filename).Is(name)
		}
	}
	test("notes.txt", "hello\n")
	test("document.pdf", "%PDF-1.7\n")
	test("picture.png", "\x89PNG\r\n\x1a\n")
}

func TestAttachmentHistoryLogReload(t *testing.T) {
	assert := assert.T(t)
	path := filepath.Join(t.TempDir(), "conversation.md")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	logger := &Agent{client: &OpenAIClient{}, logFile: f}
	t.Cleanup(logger.closeLogFile)
	log := func(msg Message) {
		logger.logMessage("user", msg.Content)
		logger.logAttachments(msg)
	}
	old := historyBinaryAttachment("old-picture.png", "image/png", 1)
	image := historyBinaryAttachment("retained-picture.png", "image/png", 5_000_000)
	log(historyAttachmentMessage("", old))
	for range 5 {
		log(historyAttachmentMessage("read this", image))
	}
	logger.closeLogFile()
	loaded := &Agent{client: &OpenAIClient{}, prompt: "current system prompt"}
	assert.This(loaded.LoadConversation(path, nil)).Is(nil)
	assert.This(len(loaded.history)).Is(7)
	assert.This(loaded.history[0]).Is(Message{Role: "system", Content: "current system prompt"})
	assert.This(loaded.history[1]).Is(historyAttachmentMessage("",
		ContentPart{Type: "text", Text: "[old-picture.png omitted due to size limit]"}))
	for _, msg := range loaded.history[2:] {
		assert.This(msg).Is(historyAttachmentMessage("read this", image))
	}
	loaded.pruneAttachments()
	assert.This(loaded.history[1].Parts[0].attachmentSize()).Is(0)
}

func TestAttachmentHistoryOriginalFilenameLabels(t *testing.T) {
	assert := assert.T(t)
	oldImage := historyBinaryAttachment("old.png", "image/png", 1)
	oldPDF := historyBinaryAttachment("old.pdf", "application/pdf", 1)
	oldText := historyTextAttachment("old.txt", "x")
	image := historyBinaryAttachment("retained.png", "image/png", 5_000_000)
	test := func(prompt string, originals []ContentPart, names ...string) {
		content := "## {{ User }}\n\n" + prompt + "\n\n"
		for _, part := range originals {
			// The original log format had no top-level filename metadata.
			data, err := json.Marshal([]ContentPart{part})
			assert.This(err).Is(nil)
			content += "## {{ Attachment }}\n\n" + string(data) + "\n\n"
		}
		agent := &Agent{}
		assert.This(agent.parseConversation(content, nil)).Is(nil)
		agent.history = append(agent.history, historyAttachmentMessage("new", image, image, image, image, image))
		agent.buildRequest()
		assert.This(agent.history[0].Content).Is(prompt)
		for i, name := range names {
			assert.This(agent.history[0].Parts[i+1]).Is(ContentPart{Type: "text",
				Text: "[" + name + " omitted due to size limit]"})
		}
		assert.This(agent.history[1]).Is(historyAttachmentMessage("new", image, image, image, image, image))
	}
	test("read\n\nAttached file: old.png", []ContentPart{oldImage}, "old.png")
	test("Attached file: old.png", []ContentPart{oldImage}, "old.png")
	test("read\n\nAttached files (in attachment order):\n1. old.txt\n2. old.pdf\n3. old.png",
		[]ContentPart{oldText, oldPDF, oldImage}, "old.txt", "old.pdf", "old.png")
}
