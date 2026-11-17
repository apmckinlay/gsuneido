// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var aiDir = ".ai"

const rawLog = false

type loggedAttachment struct {
	ContentPart
	Filename string `json:"filename,omitempty"`
}

func (agent *Agent) ensureLogFile() {
	if agent.logFile != nil {
		return
	}
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		return
	}
	stamp := time.Now().Format("20060102_150405")
	path := filepath.Join(aiDir, fmt.Sprintf("ai%s.md", stamp))
	f, err := os.Create(path)
	if err != nil {
		return
	}
	agent.logFile = f
	if rawLog {
		rawPath := filepath.Join(aiDir, fmt.Sprintf("ai%sraw.txt", stamp))
		raw, err := os.Create(rawPath)
		if err != nil {
			f.Close()
			return
		}
		agent.rawLogFile = raw
		agent.client.RawLog = raw
	}
	if agent.model != "" {
		agent.logWrite("## {{ Model }}\n\n" + agent.model + "\n\n")
	}
	if agent.loadedContent != "" {
		agent.logWrite(agent.loadedContent)
		agent.logWrite("\n## {{ Continued }}\n\n---\n\n")
		agent.loadedContent = ""
	}
}

func (agent *Agent) logMessage(role, content string) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	var marker string
	switch role {
	case "user":
		marker = "## {{ User }}\n\n"
	case "assistant":
		marker = "## {{ Assistant }}\n\n"
	case "system":
		marker = "## {{ Prompt }}\n\n"
	default:
		marker = "## {{ " + role + " }}\n\n"
	}
	agent.logWrite(marker + strings.Trim(content, "\r\n") + "\n\n")
}

func (agent *Agent) logThink(content string) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	agent.logWrite("## {{ Think }}\n\n" + strings.Trim(content, "\r\n") + "\n\n")
}

func (agent *Agent) logAttachments(msg Message) {
	// Skip only the prompt part, not text-file attachments that follow it.
	attachments := msg.Parts[msg.attachmentStart():]
	if len(attachments) == 0 {
		return
	}
	logged := make([]loggedAttachment, 0, len(attachments))
	for _, part := range attachments {
		logged = append(logged, loggedAttachment{ContentPart: part, Filename: part.Filename})
	}
	data, _ := json.Marshal(logged)
	agent.logMessage("Attachment", string(data))
}

func (agent *Agent) logAnnotations(annotations []json.RawMessage) {
	if len(annotations) == 0 {
		return
	}
	data, _ := json.Marshal(annotations)
	agent.logMessage("Annotations", string(data))
}

func (agent *Agent) flushThink() {
	if !agent.thinkDirty {
		return
	}
	agent.logThink(agent.thinkBuf.String())
	agent.thinkBuf.Reset()
	agent.thinkDirty = false
}

func (agent *Agent) logAssistantToolCalls(msg Message) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	b, _ := json.Marshal(msg)
	agent.logWrite("## {{ AssistantTool }}\n\n" + string(b) + "\n\n")
}

func (agent *Agent) logToolResult(msg Message) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	b, _ := json.Marshal(msg)
	agent.logWrite("## {{ ToolResult }}\n\n" + string(b) + "\n\n")
}

func (agent *Agent) logEffort(level string) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	agent.logWrite("## {{ Effort }}\n\n" + effortLabel(level) + "\n\n")
}

func (agent *Agent) logResponseModel(responseModel string) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	agent.logWrite("## {{ ResponseModel }}\n\n" + responseModel + "\n\n")
}

func (agent *Agent) closeLogFile() {
	agent.flushThink()
	if agent.logFile != nil {
		agent.logFile.Close()
		agent.logFile = nil
	}
	if agent.rawLogFile != nil {
		agent.rawLogFile.Close()
		agent.rawLogFile = nil
	}
	agent.client.RawLog = nil
}

func (agent *Agent) logWrite(s string) {
	if _, err := agent.logFile.WriteString(s); err != nil {
		log.Println("log write error:", err)
	}
}

// LoadConversation loads a conversation from a markdown file.
// The file should be in the format created by the logging.
// The current prompt is used, not the one from the file.
// The loaded conversation is copied to a new log file and logging continues there.
// outfn is called to restore the UI to its original state.
// Panics if a request is in progress.
func (agent *Agent) LoadConversation(path string, outfn OutFn) error {
	agent.mu.Lock()
	if agent.inProgress {
		agent.mu.Unlock()
		panic("LoadConversation: request in progress")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		agent.mu.Unlock()
		return err
	}
	agent.flushThink()
	agent.closeLogFile()
	agent.resetHistory()
	agent.loadedContent = string(data)
	agent.inProgress = true
	agent.mu.Unlock()

	defer func() {
		agent.flushThink()
		agent.mu.Lock()
		agent.inProgress = false
		agent.mu.Unlock()
	}()

	return agent.parseConversation(string(data), outfn)
}

var sectionRx = regexp.MustCompile(`(?m)^## \{\{ `)

func (agent *Agent) parseConversation(content string, outfn OutFn) error {
	sections := sectionRx.Split(content, -1)
	for _, section := range sections {
		if section == "" {
			continue
		}
		before, after, ok := strings.Cut(section, " }}")
		if !ok {
			continue
		}
		role := before
		body := after
		body = strings.TrimPrefix(body, "\n\n")
		body = strings.TrimSuffix(body, "\n\n")
		historyBody := strings.TrimRight(body, "\r\n")
		logBody := body
		switch role {
		case "Model", "Effort", "ResponseModel", "Prompt", "Continued":
			logBody = strings.TrimSpace(body)
		}
		switch role {
		case "Model":
			agent.model = logBody
		case "Effort":
			level := normalizeEffort(logBody)
			if level == "" || efforts[level] {
				agent.effort = level
				if outfn != nil {
					outfn("effort", effortLabel(level), nil)
				}
			} else {
				log.Println("parseConversation invalid effort:", logBody)
			}
		case "ResponseModel":
			agent.responseModel = logBody
			if outfn != nil {
				outfn("responseModel", logBody, nil)
			}
		case "User":
			agent.history = append(agent.history, Message{Role: "user", Content: historyBody})
			if outfn != nil {
				outfn("user", body, nil)
			}
		case "Attachment":
			if len(agent.history) == 0 || agent.history[len(agent.history)-1].Role != "user" {
				return fmt.Errorf("parseConversation Attachment: missing user message")
			}
			var parts []loggedAttachment
			if err := json.Unmarshal([]byte(body), &parts); err != nil {
				return fmt.Errorf("parseConversation Attachment: %w", err)
			}
			msg := &agent.history[len(agent.history)-1]
			if len(msg.Parts) == 0 && msg.Content != "" {
				msg.Parts = append(msg.Parts, ContentPart{Type: "text", Text: msg.Content})
			}
			for _, part := range parts {
				part.ContentPart.Filename = part.Filename
				msg.Parts = append(msg.Parts, part.ContentPart)
			}
		case "Think":
			if outfn != nil {
				outfn("think", body, nil)
			}
		case "Assistant":
			agent.history = append(agent.history, Message{Role: "assistant", Content: historyBody})
			if outfn != nil {
				outfn("output", body, nil)
			}
		case "Annotations":
			if len(agent.history) == 0 || agent.history[len(agent.history)-1].Role != "assistant" {
				return fmt.Errorf("parseConversation Annotations: missing assistant message")
			}
			var annotations []json.RawMessage
			if err := json.Unmarshal([]byte(body), &annotations); err != nil {
				return fmt.Errorf("parseConversation Annotations: %w", err)
			}
			msg := &agent.history[len(agent.history)-1]
			msg.Annotations = mergeAnnotations(msg.Annotations, annotations)
		case "Prompt", "Continued":
			// skip - use current prompt, Continued is just a marker
		case "AssistantTool":
			var msg Message
			if err := json.Unmarshal([]byte(body), &msg); err == nil {
				msg.Reasoning = "" // clear - already sent when originally created
				agent.history = append(agent.history, msg)
				if outfn != nil && len(msg.ToolCalls) > 0 {
					for _, tc := range msg.ToolCalls {
						summary, err := agent.toolClient.FormatToolCallForDisplay(tc)
						if err != nil {
							name := strings.TrimPrefix(tc.Function.Name, "suneido_")
							summary = name + " " + tc.Function.Arguments + "\n"
						}
						if !strings.HasSuffix(summary, "\n\n") {
							summary = strings.TrimRight(summary, "\r\n") + "\n\n"
						}
						outfn("tool", summary, nil)
					}
				}
			} else {
				log.Println("parseConversation AssistantTool:", err)
			}
		case "ToolResult":
			var msg Message
			if err := json.Unmarshal([]byte(body), &msg); err == nil {
				agent.history = append(agent.history, msg)
			} else {
				log.Println("parseConversation ToolResult:", err)
			}
		case "tool":
			if outfn != nil {
				outfn("tool", body, nil)
			}
		default:
			log.Println("ERROR: parseConversation unknown role:", role)
		}
	}
	agent.pruneAttachments()
	if outfn != nil {
		outfn("complete", "", nil)
	}
	return nil
}
