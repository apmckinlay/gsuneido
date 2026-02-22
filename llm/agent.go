// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Agent struct {
	client    *OpenAIClient
	mcpClient *MCPClient
	model     string
	prompt    string
	history   []Message
	outfn     OutFn
	cancel    context.CancelFunc
	logFile   *os.File
}

var aiDir = ".ai"

// OutFn is the push callback for streaming output.
// what is one of "think", "output", "tool", "complete"
type OutFn func(what, data string)

// NewAgent creates an agent.
// prompt and mcpClient are optional.
func NewAgent(baseURL, apiKey, model, prompt string, mcpClient *MCPClient, outfn OutFn) *Agent {
	agent := &Agent{
		client:    NewOpenAIClient(baseURL, apiKey),
		mcpClient: mcpClient,
		model:     model,
		prompt:    prompt,
		outfn:     outfn,
	}
	agent.resetHistory()
	return agent
}

func (agent *Agent) resetHistory() {
	agent.history = nil
	if agent.prompt != "" {
		agent.history = append(agent.history, Message{Role: "system", Content: agent.prompt})
	}
}

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
	// Log the model at the start
	if agent.model != "" {
		agent.logWrite("## Model\n\n" + agent.model + "\n\n")
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
		marker = "## User\n\n"
	case "assistant":
		marker = "## Assistant\n\n"
	case "system":
		marker = "## Prompt\n\n"
	default:
		marker = "## " + role + "\n\n"
	}
	agent.logWrite(marker + content + "\n\n")
}

func (agent *Agent) logAssistantToolCalls(msg Message) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	b, _ := json.Marshal(msg)
	agent.logWrite("## AssistantTool\n\n" + string(b) + "\n\n")
}

func (agent *Agent) logToolResult(msg Message) {
	agent.ensureLogFile()
	if agent.logFile == nil {
		return
	}
	b, _ := json.Marshal(msg)
	agent.logWrite("## ToolResult\n\n" + string(b) + "\n\n")
}

func (agent *Agent) closeLogFile() {
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
func (agent *Agent) LoadConversation(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	agent.closeLogFile()
	agent.resetHistory()
	if err := agent.parseConversation(string(data)); err != nil {
		return err
	}
	// Copy the loaded conversation to a new log file
	agent.copyToNewLogFile(string(data))
	return nil
}

// copyToNewLogFile creates a new log file and copies existing conversation content
func (agent *Agent) copyToNewLogFile(content string) {
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
	agent.logWrite(content)
	agent.logWrite("\n## Continued\n\n---\n\n")
}

func (agent *Agent) parseConversation(content string) error {
	for section := range strings.SplitSeq(content, "## ") {
		if section == "" {
			continue
		}
		lines := strings.SplitN(section, "\n", 2)
		if len(lines) < 2 {
			continue
		}
		role := strings.TrimSpace(lines[0])
		body := strings.TrimSpace(lines[1])
		switch role {
		case "Model":
			agent.model = body
		case "User":
			agent.history = append(agent.history, Message{Role: "user", Content: body})
		case "Assistant":
			agent.history = append(agent.history, Message{Role: "assistant", Content: body})
		case "Prompt", "Continued":
			// skip - use current prompt, Continued is just a marker
		case "AssistantTool":
			var msg Message
			if err := json.Unmarshal([]byte(body), &msg); err == nil {
				agent.history = append(agent.history, msg)
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
		}
	}
	return nil
}

func (agent *Agent) Input(input string) {
	go agent.request(input)
}

// Interrupt stops the current request
func (agent *Agent) Interrupt() {
	if agent.cancel != nil {
		agent.cancel()
	}
}

// SetModel sets the model to use for requests
func (agent *Agent) SetModel(model string) {
	if agent.logFile != nil && model != agent.model {
		agent.logWrite("## Model\n\n" + model + "\n\n")
	}
	agent.model = model
}

// ClearHistory clears the conversation history and starts a new log file
// The prompt is retained/restored.
func (agent *Agent) ClearHistory() {
	agent.closeLogFile()
	agent.resetHistory()
}

// request sends the request and streams the response to outfn
func (agent *Agent) request(input string) {
	ctx, cancel := context.WithCancel(context.Background())
	agent.cancel = cancel
	defer cancel()

	agent.history = append(agent.history, Message{Role: "user", Content: input})
	agent.logMessage("user", input)

	for {
		req := &ChatRequest{
			Model:            agent.model,
			Messages:         agent.history,
			IncludeReasoning: true,
		}

		if agent.mcpClient != nil {
			req.Tools = agent.mcpClient.GetTools()
			req.ToolChoice = "auto"
		}

		var content strings.Builder
		var think strings.Builder
		toolCalls := make(map[int]*ToolCall) // accumulated tool calls by index
		inThink := false

		err := agent.client.Stream(ctx, req, func(chunk *ChatCompletionChunk) error {
			if len(chunk.Choices) == 0 {
				return nil
			}
			delta := chunk.Choices[0].Delta

			// Handle reasoning_content field (used by DeepSeek R1)
			if delta.ReasoningContent != "" {
				agent.emit("think", delta.ReasoningContent)
				return nil
			}
			// Handle reasoning field (alternative used by some providers)
			if delta.Reasoning != "" {
				agent.emit("think", delta.Reasoning)
				return nil
			}

			// Handle tool calls - accumulate deltas by index
			for _, tc := range delta.ToolCalls {
				idx := max(tc.Index, 0)
				if existing, ok := toolCalls[idx]; ok {
					// Append to existing tool call
					existing.Function.Arguments += tc.Function.Arguments
					if tc.Function.Name != "" {
						existing.Function.Name = tc.Function.Name
					}
					if tc.ID != "" {
						existing.ID = tc.ID
					}
				} else {
					// New tool call
					toolCalls[idx] = &ToolCall{
						ID:   tc.ID,
						Type: tc.Type,
						Function: ToolCallFunction{
							Name:      tc.Function.Name,
							Arguments: tc.Function.Arguments,
						},
					}
				}
			}

			text := delta.Content
			if inThink {
				if strings.Contains(text, "</think") {
					inThink = false
					agent.emit("think", think.String())
					think.Reset()
					return nil
				}
				think.WriteString(text)
			} else if strings.Contains(text, "<think") {
				inThink = true
				after := text[strings.Index(text, ">")+1:]
				think.WriteString(after)
			} else {
				content.WriteString(text)
				agent.emit("output", text)
			}
			return nil
		})

		if err != nil {
			agent.emit("output", "Error: "+err.Error())
			return
		}

		// Convert map to slice
		var toolCallsList []ToolCall
		for i := 0; i < len(toolCalls); i++ {
			if tc, ok := toolCalls[i]; ok {
				toolCallsList = append(toolCallsList, *tc)
			}
		}

		// Handle tool calls
		if len(toolCallsList) > 0 && agent.mcpClient != nil {
			// Add assistant message with tool calls to history
			assistantMsg := Message{
				Role:      "assistant",
				Content:   content.String(),
				ToolCalls: toolCallsList,
			}
			agent.history = append(agent.history, assistantMsg)
			agent.logAssistantToolCalls(assistantMsg)

			// Process each tool call
			for _, tc := range toolCallsList {
				name := strings.TrimPrefix(tc.Function.Name, "suneido_")
				agent.emit("tool", "**"+name+"** "+tc.Function.Arguments+"<br>")
				result, err := agent.mcpClient.CallToolFromLLM(ctx, tc)
				if err != nil {
					agent.emit("tool", "**Error:** "+err.Error()+"<br>")
					result = "Error: " + err.Error()
				}
				// Add tool result to history
				toolMsg := Message{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
				}
				agent.history = append(agent.history, toolMsg)
				agent.logToolResult(toolMsg)
			}

			// Continue the loop to get next response
			continue
		}

		// No tool calls, we're done
		agent.history = append(agent.history, Message{Role: "assistant", Content: content.String()})
		agent.logMessage("assistant", content.String())
		agent.emit("complete", "")
		return
	}
}

func (agent *Agent) emit(what, data string) {
	agent.outfn(what, data)
}
