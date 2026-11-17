// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestReadAttachment(t *testing.T) {
	assert := assert.T(t)
	dir := t.TempDir()
	test := func(name, data, mime string) ContentPart {
		path := filepath.Join(dir, name)
		assert.This(os.WriteFile(path, []byte(data), 0600)).Is(nil)
		part, size, err := ReadAttachment(path)
		assert.This(err).Is(nil)
		assert.This(size).Is(len(data))
		assert.This(part.Filename).Is(name)
		url := "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString([]byte(data))
		if mime == "application/pdf" {
			assert.This(part.Type).Is("file")
			assert.This(part.File.Filename).Is(name)
			assert.This(part.File.Data).Is(url)
		} else {
			assert.This(part.Type).Is("image_url")
			assert.This(part.ImageURL.URL).Is(url)
		}
		return part
	}
	test("document.pdf", "%PDF-1.7\n", "application/pdf")
	test("image.png", "\x89PNG\r\n\x1a\n", "image/png")
	test("image.jpg", "\xff\xd8\xff", "image/jpeg")
	test("image.gif", "GIF89a", "image/gif")
	test("image.webp", "RIFF\x00\x00\x00\x00WEBPVP", "image/webp")
	// Detect the actual content, not just a supported filename extension.
	test("renamed.bin", "%PDF-1.7\n", "application/pdf")
	text := func(name, data string) {
		path := filepath.Join(dir, name)
		assert.This(os.WriteFile(path, []byte(data), 0600)).Is(nil)
		part, size, err := ReadAttachment(path)
		assert.This(err).Is(nil)
		assert.This(size).Is(len(data))
		assert.This(part).Is(ContentPart{Type: "text", Filename: name, Text: "Attached file: " + name + "\n\n" + data})
	}
	text("data.csv", "name,value\r\nfirst,1\r\n")
	text("server.log", "2026-11-13 started\n")
	text("color.log", "\x1b[31merror\x1b[0m\n")
	text("notes.txt", "hello \xc3\xa9\xe4\xb8\xad\n")
	text("notes.md", "# Heading\n\n## {{ User }}\n\nText\n")
	text("config.json", `{"key":"value"}`)
	text("code.go", "package main\n")
	text("empty.txt", "")

	reject := func(name, data string) {
		path := filepath.Join(dir, name)
		assert.This(os.WriteFile(path, []byte(data), 0600)).Is(nil)
		_, _, err := ReadAttachment(path)
		assert.That(err != nil)
		assert.That(strings.Contains(err.Error(), "unsupported attachment"))
	}
	reject("empty.pdf", "")
	reject("fake.png", "not an image")
	reject("file.zip", "PK\x03\x04arbitrary binary")
	_, _, err := ReadAttachment(filepath.Join(dir, "missing.pdf"))
	assert.That(err != nil)
	_, _, err = ReadAttachment(dir)
	assert.That(err != nil)

	path := filepath.Join(dir, "large.pdf")
	f, err := os.Create(path)
	assert.This(err).Is(nil)
	_, err = f.WriteString("%PDF-1.7\n")
	assert.This(err).Is(nil)
	assert.This(f.Truncate(MaxAttachmentSize + 1)).Is(nil)
	assert.This(f.Close()).Is(nil)
	_, _, err = ReadAttachment(path)
	assert.That(err != nil)
	assert.That(strings.Contains(err.Error(), "10 MB"))
}

func TestReadAttachmentLimits(t *testing.T) {
	assert := assert.T(t)
	dir := t.TempDir()
	test := func(name, data string, size int, wantError string) {
		path := filepath.Join(dir, name)
		f, err := os.Create(path)
		assert.This(err).Is(nil)
		_, err = f.WriteString(data)
		assert.This(err).Is(nil)
		assert.This(f.Truncate(int64(size))).Is(nil)
		assert.This(f.Close()).Is(nil)
		_, n, err := ReadAttachment(path)
		if wantError != "" {
			assert.That(err != nil)
			assert.That(strings.Contains(err.Error(), wantError))
		} else {
			assert.This(err).Is(nil)
			assert.This(n).Is(size)
		}
	}
	test("exact.txt", strings.Repeat("x", MaxTextAttachmentSize), MaxTextAttachmentSize, "")
	test("large.log", "text", MaxTextAttachmentSize+1, "1 MB")
	test("exact.pdf", "%PDF-1.7\n", MaxAttachmentSize, "")
	test("large.pdf", "%PDF-1.7\n", MaxAttachmentSize+1, "10 MB")
	test("exact.png", "\x89PNG\r\n\x1a\n", MaxAttachmentSize, "")
	test("large.png", "\x89PNG\r\n\x1a\n", MaxAttachmentSize+1, "10 MB")
	// Validate all text bytes, not just the MIME-sniffing prefix.
	bad := strings.Repeat("x", 512) + "\xff"
	test("invalid.txt", bad, len(bad), "UTF-8")
	bad = strings.Repeat("x", 512) + "\x00"
	test("binary.txt", bad, len(bad), "NUL")
}

func TestMessageContentParts(t *testing.T) {
	assert := assert.T(t)
	test := func(msg Message, want string) {
		data, err := json.Marshal(msg)
		assert.This(err).Is(nil)
		assert.This(string(data)).Is(want)
		var restored Message
		assert.This(json.Unmarshal(data, &restored)).Is(nil)
		assert.This(restored).Is(msg)
	}
	test(Message{Role: "user", Content: "hello"}, `{"role":"user","content":"hello"}`)
	test(Message{Role: "user", Content: "describe", Parts: []ContentPart{
		{Type: "text", Text: "describe"},
		{Type: "image_url", ImageURL: &ImageURL{URL: "data:image/png;base64,aGVsbG8="}},
	}}, `{"role":"user","content":[{"type":"text","text":"describe"},{"type":"image_url","image_url":{"url":"data:image/png;base64,aGVsbG8="}}]}`)
	test(Message{Role: "user", Content: "summarize", Parts: []ContentPart{
		{Type: "text", Text: "summarize"},
		{Type: "file", File: &FileData{Filename: "document.pdf", Data: "data:application/pdf;base64,cGRm"}},
	}}, `{"role":"user","content":[{"type":"text","text":"summarize"},{"type":"file","file":{"filename":"document.pdf","file_data":"data:application/pdf;base64,cGRm"}}]}`)
	test(Message{Role: "user", Content: "read notes.txt", Parts: []ContentPart{
		{Type: "text", Text: "read notes.txt"},
		{Type: "text", Text: "Attached file: notes.txt\n\nhello\n"},
	}}, `{"role":"user","content":[{"type":"text","text":"read notes.txt"},{"type":"text","text":"Attached file: notes.txt\n\nhello\n"}]}`)
	// Some providers return null content with assistant tool calls.
	var msg Message
	assert.This(json.Unmarshal([]byte(`{"role":"assistant","content":null}`), &msg)).Is(nil)
	assert.This(msg.Content).Is("")
}

func TestMessagePromptPart(t *testing.T) {
	assert := assert.T(t)
	notes := ContentPart{Type: "text", Filename: "notes.txt", Text: "Attached file: notes.txt\n\nhello\n"}
	image := ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: "data:image/png;base64,aGVsbG8="}}
	prompt := ContentPart{Type: "text", Text: "read"}

	assert.That(!Message{Role: "user", Content: "read"}.hasPromptPart())
	assert.This(Message{Role: "user", Content: "read", Parts: []ContentPart{prompt, image}}.attachmentStart()).Is(1)
	assert.This(Message{Role: "user", Parts: []ContentPart{image}}.attachmentStart()).Is(0)
	// Content without a matching leading text part (e.g. foreign or malformed data)
	// must not cause the first attachment to be skipped.
	assert.This(Message{Role: "user", Content: "read", Parts: []ContentPart{image}}.attachmentStart()).Is(0)
	assert.This(Message{Role: "user", Content: "read", Parts: []ContentPart{notes}}.attachmentStart()).Is(0)

	// An attachment-only text file message must round trip without
	// mistaking the attachment for the prompt.
	wire := `{"role":"user","content":[{"type":"text","text":"Attached file: notes.txt\n\nhello\n"}]}`
	var msg Message
	assert.This(json.Unmarshal([]byte(wire), &msg)).Is(nil)
	assert.This(msg.Content).Is("")
	assert.This(len(msg.Parts)).Is(1)
	assert.This(msg.attachmentStart()).Is(0)
	assert.This(msg.Parts[0].attachmentSize()).Is(len("hello\n"))

	// A prompt that happens to be only the manifest is still the prompt.
	wire = `{"role":"user","content":[{"type":"text","text":"Attached file: notes.txt"},` +
		`{"type":"text","text":"Attached file: notes.txt\n\nhello\n"}]}`
	msg = Message{}
	assert.This(json.Unmarshal([]byte(wire), &msg)).Is(nil)
	assert.This(msg.Content).Is("Attached file: notes.txt")
	assert.This(msg.attachmentStart()).Is(1)
}

func TestPruneKeepsUnpromptedTextAttachment(t *testing.T) {
	assert := assert.T(t)
	big := ContentPart{Type: "text", Filename: "big.txt",
		Text: "Attached file: big.txt\n\n" + strings.Repeat("x", MaxTotalAttachmentSize+1)}
	agent := &Agent{history: []Message{{Role: "user", Parts: []ContentPart{big}}}}
	agent.pruneAttachments()
	// The leading attachment is counted and pruned, not skipped as a prompt.
	assert.This(agent.history[0].Parts[0]).Is(
		ContentPart{Type: "text", Text: "[big.txt omitted due to size limit]"})
}
