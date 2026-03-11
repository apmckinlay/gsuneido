// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var aiDir = ".ai"

func (agent *Agent) ensureLogFile() {
	if agent.logFile != nil {
		return
	}
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		return
	}
	filename := fmt.Sprintf("ai%s.md", time.Now().Format("20060102_150405"))
	path := filepath.Join(aiDir, filename)
	f, err := os.Create(path)
	if err != nil {
		return
	}
	agent.logFile = f
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
	agent.logWrite(marker + content + "\n\n")
}

func (agent *Agent) logThink(content string) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	agent.logWrite("## {{ Think }}\n\n" + content + "\n\n")
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

func (agent *Agent) closeLogFile() {
	agent.flushThink()
	if agent.logFile != nil {
		agent.logFile.Close()
		agent.logFile = nil
	}
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

func (agent *Agent) parseConversation(content string, outfn OutFn) error {
	for section := range strings.SplitSeq(content, "## {{ ") {
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
		case "Model", "Prompt", "Continued":
			logBody = strings.TrimSpace(body)
		}
		switch role {
		case "Model":
			agent.model = logBody
		case "User":
			agent.history = append(agent.history, Message{Role: "user", Content: historyBody})
			if outfn != nil {
				outfn("user", body, nil)
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
		case "Prompt", "Continued":
			// skip - use current prompt, Continued is just a marker
		case "AssistantTool":
			var msg Message
			if err := json.Unmarshal([]byte(body), &msg); err == nil {
				msg.Reasoning = "" // clear - already sent when originally created
				agent.history = append(agent.history, msg)
				if outfn != nil && len(msg.ToolCalls) > 0 {
					for _, tc := range msg.ToolCalls {
						name := strings.TrimPrefix(tc.Function.Name, "suneido_")
						outfn("tool", name+" "+tc.Function.Arguments+"\n", nil)
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
		default:
			log.Println("ERROR: parseConversation unknown role:", role)
		}
	}
	if outfn != nil {
		outfn("complete", "", nil)
	}
	return nil
}
