// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// ContentPart is a text, image, or PDF part of a chat message.
type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
	File     *FileData `json:"file,omitempty"`
	Filename string    `json:"-"` // local metadata, persisted separately from provider fields
}

type ImageURL struct {
	URL string `json:"url"`
}

type FileData struct {
	Filename string `json:"filename"`
	Data     string `json:"file_data"`
}

// Attachment limits apply to raw bytes, before encoding or filename labels.
const (
	MaxTextAttachmentSize  = 1_000_000
	MaxAttachmentSize      = 10_000_000
	MaxTotalAttachmentSize = 25_000_000
)

// ReadAttachment snapshots a file and returns its unencoded byte count.
// Callers must apply their file-access permissions before calling this function.
func ReadAttachment(path string) (ContentPart, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return ContentPart{}, 0, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return ContentPart{}, 0, err
	}
	if !info.Mode().IsRegular() {
		return ContentPart{}, 0, fmt.Errorf("attachment must be a regular file")
	}
	r := bufio.NewReader(f)
	header, err := r.Peek(512)
	if err != nil && err != io.EOF {
		return ContentPart{}, 0, err
	}
	mime := http.DetectContentType(header)
	text := false
	switch mime {
	case "application/pdf", "image/png", "image/jpeg", "image/gif", "image/webp":
	default:
		text = strings.HasPrefix(mime, "text/")
		// A broken PDF or image must not silently become a text attachment.
		switch strings.ToLower(filepath.Ext(path)) {
		case ".csv", ".log", ".txt", ".md":
			text = true
		case ".pdf", ".png", ".jpg", ".jpeg", ".gif", ".webp":
			text = false
		}
		if !text {
			return ContentPart{}, 0, fmt.Errorf("unsupported attachment: select a text, PDF, PNG, JPEG, GIF, or WebP file")
		}
	}
	limit, limitError := MaxAttachmentSize, "attachment exceeds the 10 MB limit"
	if text {
		limit, limitError = MaxTextAttachmentSize, "text attachment exceeds the 1 MB limit"
	}
	if info.Size() > int64(limit) {
		return ContentPart{}, 0, fmt.Errorf("%s", limitError)
	}
	data, err := io.ReadAll(io.LimitReader(r, int64(limit)+1))
	if err != nil {
		return ContentPart{}, 0, err
	}
	if len(data) > limit {
		return ContentPart{}, 0, fmt.Errorf("%s", limitError)
	}
	if text {
		if !utf8.Valid(data) || bytes.IndexByte(data, 0) >= 0 {
			return ContentPart{}, 0, fmt.Errorf("text attachment must be UTF-8 without NUL bytes")
		}
		return ContentPart{Type: "text", Filename: filepath.Base(path),
			Text: "Attached file: " + filepath.Base(path) + "\n\n" + string(data)}, len(data), nil
	}
	url := "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data)
	if mime == "application/pdf" {
		return ContentPart{Type: "file", Filename: filepath.Base(path), File: &FileData{
			Filename: filepath.Base(path), Data: url}}, len(data), nil
	}
	return ContentPart{Type: "image_url", Filename: filepath.Base(path), ImageURL: &ImageURL{URL: url}}, len(data), nil
}

// isAttachmentText returns true for the labeled text part
// that ReadAttachment creates for text files ("Attached file: name\n\nbody").
// The "\n\n" distinguishes it from a prompt that is only the file manifest
// (e.g. "Attached file: name" when the user entered no text).
func (part ContentPart) isAttachmentText() bool {
	return part.Type == "text" && strings.HasPrefix(part.Text, "Attached file: ") &&
		strings.Contains(part.Text, "\n\n")
}

// hasPromptPart returns true if Parts[0] is the prompt text rather than an attachment.
// The prompt part is only present when Content is non-empty,
// and it always duplicates Content.
func (m Message) hasPromptPart() bool {
	return m.Content != "" && len(m.Parts) > 0 &&
		m.Parts[0].Type == "text" && m.Parts[0].Text == m.Content
}

// attachmentStart returns the index in Parts of the first attachment.
func (m Message) attachmentStart() int {
	if m.hasPromptPart() {
		return 1
	}
	return 0
}

func (part ContentPart) attachmentSize() int {
	var data string
	switch {
	case part.isAttachmentText():
		_, body, _ := strings.Cut(part.Text, "\n\n")
		return len(body)
	case part.Type == "file" && part.File != nil:
		data = part.File.Data
	case part.Type == "image_url" && part.ImageURL != nil:
		data = part.ImageURL.URL
	default:
		return 0
	}
	_, encoded, ok := strings.Cut(data, ";base64,")
	if !ok {
		return 0
	}
	size := base64.StdEncoding.DecodedLen(len(encoded))
	if strings.HasSuffix(encoded, "==") {
		size -= 2
	} else if strings.HasSuffix(encoded, "=") {
		size--
	}
	return max(0, size)
}

func (agent *Agent) pruneAttachments() {
	total := 0
	for _, msg := range agent.history {
		if msg.Role != "user" {
			continue
		}
		for _, part := range msg.Parts[msg.attachmentStart():] {
			total += part.attachmentSize()
		}
	}
	for _, msg := range agent.history {
		if total <= MaxTotalAttachmentSize {
			return
		}
		if msg.Role != "user" {
			continue
		}
		for i := msg.attachmentStart(); i < len(msg.Parts); i++ {
			part := &msg.Parts[i]
			size := part.attachmentSize()
			if size == 0 {
				continue
			}
			name := part.Filename
			if name == "" && part.File != nil {
				name = part.File.Filename
			}
			if name == "" && part.isAttachmentText() {
				label, _, _ := strings.Cut(part.Text, "\n\n")
				name = strings.TrimPrefix(label, "Attached file: ")
			}
			if name == "" {
				// Existing image logs stored filenames only in the user message.
				manifest := msg.Content
				if pos := strings.LastIndex(manifest, "\n\nAttached "); pos >= 0 {
					manifest = manifest[pos+2:]
				}
				if file, ok := strings.CutPrefix(manifest, "Attached file: "); ok {
					name = file
				} else if list, ok := strings.CutPrefix(manifest, "Attached files (in attachment order):\n"); ok {
					index := i - msg.attachmentStart()
					files := strings.Split(list, "\n")
					if index >= 0 && index < len(files) {
						name, _ = strings.CutPrefix(files[index], fmt.Sprintf("%d. ", index+1))
					}
				}
			}
			if name == "" {
				name = "attachment"
			}
			if part.Type == "file" {
				agent.dropFileAnnotations(name)
			}
			*part = ContentPart{Type: "text", Text: "[" + name + " omitted due to size limit]"}
			total -= size
			if total <= MaxTotalAttachmentSize {
				return
			}
		}
	}
}

type annotationHeader struct {
	Type string `json:"type"`
	File struct {
		Hash string `json:"hash"`
		Name string `json:"name"`
	} `json:"file"`
}

// Keep annotations opaque so provider fields and parsed contents survive replay.
func mergeAnnotations(dst, src []json.RawMessage) []json.RawMessage {
	for _, annotation := range src {
		var incoming annotationHeader
		_ = json.Unmarshal(annotation, &incoming)
		found := false
		for i, previous := range dst {
			var existing annotationHeader
			_ = json.Unmarshal(previous, &existing)
			if bytes.Equal(previous, annotation) ||
				(incoming.Type == "file" && existing.Type == "file" &&
					incoming.File.Hash != "" && incoming.File.Hash == existing.File.Hash) {
				dst[i] = annotation
				found = true
				break
			}
		}
		if !found {
			dst = append(dst, annotation)
		}
	}
	return dst
}

func (agent *Agent) dropFileAnnotations(name string) {
	for i := range agent.history {
		annotations := agent.history[i].Annotations
		kept := annotations[:0]
		for _, annotation := range annotations {
			var header annotationHeader
			_ = json.Unmarshal(annotation, &header)
			// An unnamed parsed file cannot safely be retained after PDF eviction.
			if header.Type == "file" &&
				(header.File.Name == name || header.File.Name == "" || name == "attachment") {
				continue
			}
			kept = append(kept, annotation)
		}
		clear(annotations[len(kept):])
		agent.history[i].Annotations = kept
	}
}

func (m Message) MarshalJSON() ([]byte, error) {
	// Text-only messages retain their original wire format and log representation.
	type plain Message
	if len(m.Parts) == 0 {
		return json.Marshal(plain(m))
	}
	return json.Marshal(struct {
		plain
		Content []ContentPart `json:"content"`
	}{plain: plain(m), Content: m.Parts})
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type plain Message
	var wire struct {
		plain
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*m = Message(wire.plain)
	content := strings.TrimSpace(string(wire.Content))
	if strings.HasPrefix(content, "[") {
		if err := json.Unmarshal(wire.Content, &m.Parts); err != nil {
			return err
		}
		// A leading labeled text-file part is an attachment (the message had no
		// prompt text), not the prompt. Otherwise hasPromptPart would skip it.
		if len(m.Parts) > 0 && m.Parts[0].Type == "text" && !m.Parts[0].isAttachmentText() {
			m.Content = m.Parts[0].Text
		}
		return nil
	}
	if len(wire.Content) == 0 {
		return nil
	}
	return json.Unmarshal(wire.Content, &m.Content)
}
