// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/llm"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestAgentInputAttachmentErrors(t *testing.T) {
	assert := assert.T(t)
	dir := t.TempDir()
	file := filepath.Join(dir, "not-an-image.png")
	assert.This(os.WriteFile(file, []byte("not an image"), 0600)).Is(nil)
	valid := filepath.Join(dir, "valid.pdf")
	assert.This(os.WriteFile(valid, []byte("%PDF-1.7\n"), 0600)).Is(nil)
	test := func(files Value, want string) {
		th := &Thread{}
		// A nil agent ensures validation fails before starting any request.
		result := agent_Input(th, &suAgent{}, []Value{SuStr("draft"), files})
		assert.That(th.ReturnThrow)
		assert.That(strings.Contains(ToStr(result), want))
	}
	test(SuStr(file), "unsupported attachment")
	test(SuStr(filepath.Join(dir, "missing.pdf")), "AiAgent.Input:")
	test(SuObjectOf(SuStr(valid), SuStr(file)), "unsupported attachment")
	test(SuObjectOf(SuStr(valid), SuStr(filepath.Join(dir, "missing.pdf"))), "AiAgent.Input:")

	enableSandbox(dir)
	defer DisableSandbox()
	test(SuStr(file), "absolute paths disabled in sandbox")
	test(SuStr("../outside.pdf"), "parent paths disabled in sandbox")
	test(SuStr("not-an-image.png"), "unsupported attachment")
	test(SuObjectOf(SuStr("valid.pdf"), SuStr(file)), "absolute paths disabled in sandbox")
	test(SuObjectOf(SuStr("valid.pdf"), SuStr("../outside.pdf")), "parent paths disabled in sandbox")

	outside := filepath.Join(t.TempDir(), "outside.pdf")
	assert.This(os.WriteFile(outside, []byte("%PDF-1.7\n"), 0600)).Is(nil)
	if os.Symlink(outside, filepath.Join(dir, "link.pdf")) == nil {
		test(SuStr("link.pdf"), "path outside sandbox")
	}
}

func TestAgentInputAttachmentLimits(t *testing.T) {
	assert := assert.T(t)
	t.Chdir(t.TempDir())
	dir := t.TempDir()
	pdf := func(name string, size int64) string {
		path := filepath.Join(dir, name)
		assert.This(os.WriteFile(path, []byte("%PDF-1.7\n"), 0600)).Is(nil)
		assert.This(os.Truncate(path, size)).Is(nil)
		return path
	}
	text := filepath.Join(dir, "oversize.txt")
	assert.This(os.WriteFile(text, []byte(strings.Repeat(" ", 1_000_001)), 0600)).Is(nil)
	valid := pdf("valid.pdf", 9)
	oversize := pdf("oversize.pdf", 10_000_001)
	first := pdf("first.pdf", 10_000_000)
	second := pdf("second.pdf", 10_000_000)
	third := pdf("third.pdf", 5_000_001)
	test := func(files Value) {
		th := &Thread{}
		// No agent: even a valid earlier file must not start a partial request.
		result := agent_Input(th, &suAgent{}, []Value{SuStr("draft"), files})
		assert.That(th.ReturnThrow)
		assert.That(strings.HasPrefix(ToStr(result), "AiAgent.Input:"))
		assert.That(strings.Contains(ToStr(result), "limit"))
	}
	test(SuStr(text))
	test(SuObjectOf(SuStr(valid), SuStr(text)))
	test(SuStr(oversize))
	test(SuObjectOf(SuStr(valid), SuStr(oversize)))
	test(SuObjectOf(SuStr(first), SuStr(second), SuStr(third)))

	assert.This(os.Truncate(third, 5_000_000)).Is(nil)
	started := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()
	done := make(chan struct{}, 1)
	th := &Thread{}
	agent := llm.NewAgent(th, server.URL, "", "model", "",
		func(_ string, data string, _ *llm.ToolApproval) {
			if strings.HasPrefix(data, "ERROR:") {
				t.Error(data)
			}
		}, func() {}, func() { done <- struct{}{} })
	defer func() {
		close(release)
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("request did not complete")
		}
		agent.ClearHistory()
	}()
	agent.Input("hold request open")
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("request did not start")
	}
	// Exactly 25,000,000 raw bytes must pass validation and reach the busy guard,
	// without serializing or sending a large HTTP request.
	assert.This(func() {
		agent_Input(th, &suAgent{agent: agent}, []Value{SuStr("draft"),
			SuObjectOf(SuStr(first), SuStr(second), SuStr(third))})
	}).Panics("Input: request already in progress")
	assert.That(!th.ReturnThrow)
}

func TestAgentInputMultipleFiles(t *testing.T) {
	assert := assert.T(t)
	// Keep conversation logs out of the workspace.
	t.Chdir(t.TempDir())
	dir := t.TempDir()
	pdf := filepath.Join(dir, "document.pdf")
	first := filepath.Join(dir, "first.png")
	second := filepath.Join(dir, "second.png")
	pdfData := "%PDF-1.7\n"
	firstData := "\x89PNG\r\n\x1a\nfirst"
	secondData := "\x89PNG\r\n\x1a\nsecond"
	assert.This(os.WriteFile(pdf, []byte(pdfData), 0600)).Is(nil)
	assert.This(os.WriteFile(first, []byte(firstData), 0600)).Is(nil)
	assert.This(os.WriteFile(second, []byte(secondData), 0600)).Is(nil)
	requests := make(chan llm.ChatRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req llm.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		requests <- req
		assert.That(req.Stream)
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"reply\"}}]}\n\ndata: [DONE]\n\n")
	}))
	defer server.Close()
	done := make(chan struct{}, 1)
	th := &Thread{}
	agent := llm.NewAgent(th, server.URL, "", "model", "",
		func(_ string, data string, _ *llm.ToolApproval) {
			if strings.HasPrefix(data, "ERROR:") {
				t.Error(data)
			}
		}, func() {}, func() { done <- struct{}{} })
	defer agent.ClearHistory()
	su := &suAgent{agent: agent}
	test := func(input string, files Value, want string, attachments ...llm.ContentPart) {
		assert.This(agent_Input(th, su, []Value{SuStr(input), files})).Is(nil)
		assert.That(!th.ReturnThrow)
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("request did not complete")
		}
		req := <-requests
		assert.This(len(req.Messages)).Is(1)
		assert.This(req.Messages[0].Role).Is("user")
		assert.This(req.Messages[0].Content).Is(want)
		if len(attachments) == 0 {
			assert.This(len(req.Messages[0].Parts)).Is(0)
		} else {
			parts := append([]llm.ContentPart{{Type: "text", Text: want}}, attachments...)
			assert.This(req.Messages[0].Parts).Is(parts)
		}
		agent.ClearHistory()
	}
	pdfPart := llm.ContentPart{Type: "file", File: &llm.FileData{
		Filename: "document.pdf", Data: "data:application/pdf;base64," + base64.StdEncoding.EncodeToString([]byte(pdfData))}}
	firstPart := llm.ContentPart{Type: "image_url", ImageURL: &llm.ImageURL{
		URL: "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(firstData))}}
	secondPart := llm.ContentPart{Type: "image_url", ImageURL: &llm.ImageURL{
		URL: "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(secondData))}}
	text := func(name, data string) (Value, llm.ContentPart) {
		path := filepath.Join(dir, name)
		assert.This(os.WriteFile(path, []byte(data), 0600)).Is(nil)
		return SuStr(path), llm.ContentPart{Type: "text", Text: "Attached file: " + name + "\n\n" + data}
	}
	csv, csvPart := text("report.csv", "name,value\nalpha,1\n")
	log, logPart := text("application.log", "2026-11-13 INFO started\r\n")
	txt, txtPart := text("notes.txt", "plain UTF-8: \xc3\xa9 \xe4\xb8\xad\n")
	md, mdPart := text("README.md", "# Heading\n\n**markdown**\n")
	maxText, maxTextPart := text("maximum.txt", strings.Repeat(" ", 1_000_000))
	test("no files", SuObjectOf(), "no files")
	test("no files", EmptyStr, "no files")
	test("one file", SuStr(pdf), "one file\n\nAttached file: document.pdf", pdfPart)
	test("one file", SuObjectOf(SuStr(first)), "one file\n\nAttached file: first.png", firstPart)
	test("Compare first.png and second.png with document.pdf",
		SuObjectOf(SuStr(pdf), SuStr(first), SuStr(second)),
		"Compare first.png and second.png with document.pdf\n\nAttached files (in attachment order):\n1. document.pdf\n2. first.png\n3. second.png",
		pdfPart, firstPart, secondPart)
	test("", SuObjectOf(SuStr(first), SuStr(second)),
		"Attached files (in attachment order):\n1. first.png\n2. second.png", firstPart, secondPart)
	test("Read CSV", csv, "Read CSV\n\nAttached file: report.csv", csvPart)
	test("Read log", SuObjectOf(log), "Read log\n\nAttached file: application.log", logPart)
	test("Read text", txt, "Read text\n\nAttached file: notes.txt", txtPart)
	test("", md, "Attached file: README.md", mdPart)
	test("At text limit", maxText, "At text limit\n\nAttached file: maximum.txt", maxTextPart)
	test("Compare files", SuObjectOf(csv, SuStr(pdf), txt, SuStr(first), log, md),
		"Compare files\n\nAttached files (in attachment order):\n1. report.csv\n2. document.pdf\n3. notes.txt\n4. first.png\n5. application.log\n6. README.md",
		csvPart, pdfPart, txtPart, firstPart, logPart, mdPart)
}
