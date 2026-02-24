// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"os"
	"strings"
	"sync"
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

		content, reasoning, toolCalls, err := agent.doStream(ctx, req)
		if err != nil {
			agent.emit("output", "Error: "+err.Error())
			agent.emit("complete", "")
			return
		}

		// Clear reasoning from previous assistant messages (already sent once)
		agent.clearReasoning()

		// Handle tool calls
		if len(toolCalls) > 0 && agent.mcpClient != nil {
			agent.processToolCalls(ctx, content, reasoning, toolCalls)
			continue // Continue the loop to get next response
		}

		// No tool calls, we're done
		agent.history = append(agent.history,
			Message{Role: "assistant", Content: content, Reasoning: truncateReasoning(reasoning)})
		agent.logMessage("assistant", content)
		agent.emit("complete", "")
		return
	}
}

// clearReasoning removes reasoning from all assistant messages in history
// to avoid accumulating it in context after it's been sent once
func (agent *Agent) clearReasoning() {
	for i := range agent.history {
		if agent.history[i].Role == "assistant" {
			agent.history[i].Reasoning = ""
		}
	}
}

const (
	reasoningLimit = reasoningHead + reasoningTail
	reasoningHead  = 600
	reasoningTail  = 1400
)

// truncateReasoning limits reasoning size by keeping head and tail portions
func truncateReasoning(s string) string {
	if len(s) <= reasoningLimit {
		return s
	}
	head := s[:reasoningHead]
	tail := s[len(s)-reasoningTail:]
	// Trim partial line from end of head
	if idx := strings.LastIndexByte(head, '\n'); idx >= 0 {
		head = head[:idx]
	}
	// Trim partial line from start of tail
	if idx := strings.IndexByte(tail, '\n'); idx >= 0 {
		tail = tail[idx+1:]
	}
	return head + "\n...[truncated]...\n" + tail
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

// doStream performs the streaming request and returns accumulated content and reasoning
func (agent *Agent) doStream(ctx context.Context, req *ChatRequest) (
	content, reasoning string, toolCalls []ToolCall, err error) {

	var contentBuilder, reasoningBuilder strings.Builder
	toolCallsMap := make(map[int]*ToolCall)
	inThink := false

	err = agent.client.Stream(ctx, req, func(chunk *ChatCompletionChunk) error {
		return agent.handleStreamChunk(chunk, &contentBuilder, &reasoningBuilder, &toolCallsMap, &inThink)
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
	return contentBuilder.String(), reasoningBuilder.String(), toolCalls, nil
}

// handleStreamChunk processes a single streaming chunk
func (agent *Agent) handleStreamChunk(chunk *ChatCompletionChunk, content *strings.Builder,
	reasoning *strings.Builder, toolCallsMap *map[int]*ToolCall, inThink *bool) error {

	if len(chunk.Choices) == 0 {
		return nil
	}
	delta := chunk.Choices[0].Delta

	// Handle reasoning_content field (used by DeepSeek R1)
	if delta.ReasoningContent != "" {
		reasoning.WriteString(delta.ReasoningContent)
		agent.emit("think", delta.ReasoningContent)
		return nil
	}
	// Handle reasoning field (alternative used by some providers)
	if delta.Reasoning != "" {
		reasoning.WriteString(delta.Reasoning)
		agent.emit("think", delta.Reasoning)
		return nil
	}

	// Handle tool calls - accumulate deltas by index
	agent.accumulateToolCalls(delta.ToolCalls, toolCallsMap)

	text := delta.Content
	return agent.processContentText(text, content, reasoning, inThink)
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
func (agent *Agent) processContentText(text string, content *strings.Builder,
	reasoning *strings.Builder, inThink *bool) error {
	if *inThink {
		if strings.Contains(text, "</think") {
			*inThink = false
			agent.flushThink()
			return nil
		}
		reasoning.WriteString(text)
		agent.emit("think", text)
	} else if strings.Contains(text, "<think") {
		*inThink = true
		after := text[strings.Index(text, ">")+1:]
		reasoning.WriteString(after)
		agent.emit("think", after)
	} else {
		content.WriteString(text)
		agent.emit("output", text)
	}
	return nil
}

// processToolCalls processes tool calls and adds results to history
func (agent *Agent) processToolCalls(ctx context.Context, content, reasoning string, toolCallsList []ToolCall) {
	// Add assistant message with tool calls to history
	assistantMsg := Message{
		Role:      "assistant",
		Content:   content,
		Reasoning: truncateReasoning(reasoning),
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
