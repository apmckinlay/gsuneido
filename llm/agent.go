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
	"sync"
	"time"
)

type Agent struct {
	client        *OpenAIClient
	mcpClient     *MCPClient
	model         string
	prompt        string
	history       []Message
	outfn         OutFn
	cancel        context.CancelFunc
	logFile       *os.File
	mu            sync.Mutex
	inProgress    bool
	thinkBuf      strings.Builder
	thinkDirty    bool
	loadedContent string // content from LoadConversation to copy when log is created
}

var aiDir = ".ai"

// OutFn is the push callback for streaming output.
// what is one of "user", "think", "output", "tool", "complete"
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
		agent.logWrite("## {{ Model }}\n\n" + agent.model + "\n\n")
	}
	// Copy loaded conversation content if any
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
	agent.loadedContent = string(data) // will be copied when log file is created
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
		// Find the closing }}
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
				outfn("user", body)
			}
		case "Think":
			if outfn != nil {
				outfn("think", body)
			}
		case "Assistant":
			agent.history = append(agent.history, Message{Role: "assistant", Content: historyBody})
			if outfn != nil {
				outfn("output", body)
			}
		case "Prompt", "Continued":
			// skip - use current prompt, Continued is just a marker
		case "AssistantTool":
			var msg Message
			if err := json.Unmarshal([]byte(body), &msg); err == nil {
				agent.history = append(agent.history, msg)
				if outfn != nil && len(msg.ToolCalls) > 0 {
					for _, tc := range msg.ToolCalls {
						name := strings.TrimPrefix(tc.Function.Name, "suneido_")
						outfn("tool", "**"+name+"** "+tc.Function.Arguments+"<br>")
					}
				}
			} else {
				log.Println("parseConversation AssistantTool:", err)
			}
		case "ToolResult":
			var msg Message
			if err := json.Unmarshal([]byte(body), &msg); err == nil {
				agent.history = append(agent.history, msg)
				// Tool result content is typically just logged, not displayed to UI
			} else {
				log.Println("parseConversation ToolResult:", err)
			}
		}
	}
	if outfn != nil {
		outfn("complete", "")
	}
	return nil
}

func (agent *Agent) Input(input string) {
	agent.mu.Lock()
	if agent.inProgress {
		agent.mu.Unlock()
		panic("Input: request already in progress")
	}
	agent.inProgress = true
	agent.mu.Unlock()
	go agent.request(input)
}

// Interrupt stops the current request
func (agent *Agent) Interrupt() {
	if agent.cancel != nil {
		agent.cancel()
	}
}

// SetModel sets the model to use for requests
// Panics if a request is in progress.
func (agent *Agent) SetModel(model string) {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	if agent.inProgress {
		panic("SetModel: request in progress")
	}
	if agent.logFile != nil && model != agent.model {
		agent.logWrite("## {{ Model }}\n\n" + model + "\n\n")
	}
	agent.model = model
}

// ClearHistory clears the conversation history and starts a new log file
// The prompt is retained/restored.
// Panics if a request is in progress.
func (agent *Agent) ClearHistory() {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	if agent.inProgress {
		panic("ClearHistory: request in progress")
	}
	agent.flushThink()
	agent.closeLogFile()
	agent.loadedContent = ""
	agent.resetHistory()
}

// request sends the request and streams the response to outfn
func (agent *Agent) request(input string) {
	defer func() {
		agent.mu.Lock()
		agent.inProgress = false
		agent.mu.Unlock()
	}()
	ctx, cancel := context.WithCancel(context.Background())
	agent.cancel = cancel
	defer cancel()

	agent.history = append(agent.history, Message{Role: "user", Content: input})
	agent.logMessage("user", input)

	for {
		req := agent.buildRequest()

		content, toolCalls, err := agent.doStream(ctx, req)
		if err != nil {
			agent.emit("output", "Error: "+err.Error())
			agent.emit("complete", "")
			return
		}

		// Handle tool calls
		if len(toolCalls) > 0 && agent.mcpClient != nil {
			agent.processToolCalls(ctx, content.String(), toolCalls)
			continue // Continue the loop to get next response
		}

		// No tool calls, we're done
		agent.history = append(agent.history, Message{Role: "assistant", Content: content.String()})
		agent.logMessage("assistant", content.String())
		agent.emit("complete", "")
		return
	}
}

// buildRequest creates the chat request with current history and tools
func (agent *Agent) buildRequest() *ChatRequest {
	req := &ChatRequest{
		Model:            agent.model,
		Messages:         agent.history,
		IncludeReasoning: true,
	}
	if agent.mcpClient != nil {
		req.Tools = agent.mcpClient.GetTools()
		req.ToolChoice = "auto"
	}
	return req
}

// doStream performs the streaming request and returns accumulated content
func (agent *Agent) doStream(ctx context.Context, req *ChatRequest) (
	content strings.Builder, toolCalls []ToolCall, err error) {

	toolCallsMap := make(map[int]*ToolCall)
	inThink := false

	err = agent.client.Stream(ctx, req, func(chunk *ChatCompletionChunk) error {
		return agent.handleStreamChunk(chunk, &content, &toolCallsMap, &inThink)
	})

	if err != nil {
		return
	}

	// Convert map to slice
	for i := 0; i < len(toolCallsMap); i++ {
		if tc, ok := toolCallsMap[i]; ok {
			toolCalls = append(toolCalls, *tc)
		}
	}
	return
}

// handleStreamChunk processes a single streaming chunk
func (agent *Agent) handleStreamChunk(chunk *ChatCompletionChunk, content *strings.Builder,
	toolCallsMap *map[int]*ToolCall, inThink *bool) error {

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
	agent.accumulateToolCalls(delta.ToolCalls, toolCallsMap)

	text := delta.Content
	return agent.processContentText(text, content, inThink)
}

// accumulateToolCalls accumulates streaming tool call data by index
func (agent *Agent) accumulateToolCalls(deltaToolCalls []ToolCall, toolCallsMap *map[int]*ToolCall) {
	for _, tc := range deltaToolCalls {
		idx := max(tc.Index, 0)
		if existing, ok := (*toolCallsMap)[idx]; ok {
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
			(*toolCallsMap)[idx] = &ToolCall{
				ID:   tc.ID,
				Type: tc.Type,
				Function: ToolCallFunction{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			}
		}
	}
}

// processContentText handles the content text, tracking think blocks
func (agent *Agent) processContentText(text string, content *strings.Builder, inThink *bool) error {
	if *inThink {
		if strings.Contains(text, "</think") {
			*inThink = false
			agent.flushThink()
			return nil
		}
		agent.emit("think", text)
	} else if strings.Contains(text, "<think") {
		*inThink = true
		after := text[strings.Index(text, ">")+1:]
		agent.emit("think", after)
	} else {
		content.WriteString(text)
		agent.emit("output", text)
	}
	return nil
}

// processToolCalls processes tool calls and adds results to history
func (agent *Agent) processToolCalls(ctx context.Context, content string, toolCallsList []ToolCall) {
	// Add assistant message with tool calls to history
	assistantMsg := Message{
		Role:      "assistant",
		Content:   content,
		ToolCalls: toolCallsList,
	}
	agent.history = append(agent.history, assistantMsg)
	agent.logAssistantToolCalls(assistantMsg)

	// Process each tool call
	for _, tc := range toolCallsList {
		agent.executeSingleToolCall(ctx, tc)
	}
}

// executeSingleToolCall executes a single tool call and logs the result
func (agent *Agent) executeSingleToolCall(ctx context.Context, tc ToolCall) {
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

func (agent *Agent) emit(what, data string) {
	if what == "think" {
		agent.thinkBuf.WriteString(data)
		agent.thinkDirty = true
		agent.outfn(what, data)
		return
	}
	agent.flushThink()
	agent.outfn(what, data)
}
