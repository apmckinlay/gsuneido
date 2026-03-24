// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Agent struct {
	client        *OpenAIClient
	toolClient    *ToolClient
	model         string
	prompt        string
	history       []Message
	outfn         OutFn
	cancel        context.CancelFunc
	logFile       *os.File
	rawLogFile    *os.File
	mu            sync.Mutex
	inProgress    bool
	thinkBuf      strings.Builder
	thinkDirty    bool
	loadedContent string  // content from LoadConversation to copy when log is created
	usage         *Usage  // token usage from last response
	totalCost     float64 // cumulative cost across all requests
}

// OutFn is the push callback for streaming output.
// what is one of "user", "think", "output", "tool", "complete"
// approval is non-nil when a tool call requires user approval.
type OutFn func(what, data string, approval *ToolApproval)

type ToolApproval struct {
	Before string
	After  string
	once   sync.Once
	ch     chan approvalDecision
}

type approvalDecision struct {
	allow bool
	text  string
}

func newToolApproval() *ToolApproval {
	return &ToolApproval{ch: make(chan approvalDecision, 1)}
}

func (a *ToolApproval) Allow(text string) {
	a.once.Do(func() {
		a.ch <- approvalDecision{allow: true, text: text}
	})
}

func (a *ToolApproval) Deny(text string) {
	a.once.Do(func() {
		a.ch <- approvalDecision{allow: false, text: text}
	})
}

func (a *ToolApproval) Wait(ctx context.Context) (approvalDecision, error) {
	select {
	case d := <-a.ch:
		return d, nil
	case <-ctx.Done():
		return approvalDecision{}, ctx.Err()
	}
}

// NewAgent creates an agent.
// prompt is optional.
func NewAgent(baseURL, apiKey, model, prompt string, outfn OutFn) *Agent {
	toolClient, err := NewToolClient()
	if err != nil {
		panic("NewAgent: " + err.Error())
	}

	agent := &Agent{
		client:     NewOpenAIClient(baseURL, apiKey),
		toolClient: toolClient,
		model:      model,
		prompt:     prompt,
		outfn:      outfn,
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

// LastUsage returns the token usage from the last response.
// Returns nil if no usage information is available.
func (agent *Agent) LastUsage() *Usage {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	return agent.usage
}

// TotalCost returns the cumulative cost across all requests.
func (agent *Agent) TotalCost() float64 {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	return agent.totalCost
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
	agent.usage = nil
	agent.totalCost = 0
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
			agent.emit("output", "ERROR: "+err.Error())
			agent.emit("complete", "")
			return
		}

		// Clear reasoning from previous assistant messages (already sent once)
		agent.clearReasoning()

		// Handle tool calls
		if len(toolCalls) > 0 && agent.toolClient != nil {
			agent.processToolCalls(ctx, content, reasoning, toolCalls)
			continue // Continue the loop to get next response
		}

		// No tool calls, we're done
		agent.history = append(agent.history,
			Message{Role: "assistant", Content: content, Reasoning: truncate(reasoning)})
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
	truncateLimit = truncateHead + truncateTail
	truncateHead  = 600
	truncateTail  = 1400
)

// truncate limits string size by keeping head and tail portions
func truncate(s string) string {
	if len(s) <= truncateLimit {
		return s
	}
	head := s[:truncateHead]
	tail := s[len(s)-truncateTail:]
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
		Plugins: []Plugin{
			{ID: "context-compression"},
			{ID: "response-healing"},
		},
	}
	if agent.toolClient != nil {
		req.Tools = agent.toolClient.GetTools()
		req.ToolChoice = "auto"
	}
	return req
}

// doStream performs the streaming request and returns accumulated content and reasoning
func (agent *Agent) doStream(ctx context.Context, req *ChatRequest) (
	content, reasoning string, toolCalls []ToolCall, err error) {

	var contentBuilder, reasoningBuilder strings.Builder
	toolCallsState := make([]*streamToolCall, 0, 2)
	inThink := false

	err = agent.client.Stream(ctx, req, func(chunk *ChatCompletionChunk) error {
		return agent.handleStreamChunk(chunk, &contentBuilder, &reasoningBuilder, &toolCallsState, &inThink)
	})
	if err != nil {
		return
	}

	toolCalls = make([]ToolCall, 0, len(toolCallsState))
	for _, st := range toolCallsState {
		if st == nil {
			continue
		}
		st.call.Function.Arguments = st.args.String()
		toolCalls = append(toolCalls, st.call)
	}
	return contentBuilder.String(), reasoningBuilder.String(), toolCalls, nil
}

type streamToolCall struct {
	call ToolCall
	args strings.Builder
}

// handleStreamChunk processes a single streaming chunk
func (agent *Agent) handleStreamChunk(chunk *ChatCompletionChunk, content *strings.Builder,
	reasoning *strings.Builder, toolCallsState *[]*streamToolCall, inThink *bool) error {

	// Capture usage if present (typically in final chunk)
	if chunk.Usage != nil {
		agent.mu.Lock()
		agent.usage = chunk.Usage
		agent.totalCost += chunk.Usage.Cost
		agent.mu.Unlock()
	}

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
	agent.accumulateToolCalls(delta.ToolCalls, toolCallsState)

	text := delta.Content
	return agent.processContentText(text, content, reasoning, inThink)
}

// accumulateToolCalls accumulates streaming tool call data by index
func (agent *Agent) accumulateToolCalls(deltaToolCalls []ToolCall, toolCallsState *[]*streamToolCall) {
	for _, tc := range deltaToolCalls {
		idx := max(tc.Index, 0)
		if idx >= len(*toolCallsState) {
			*toolCallsState = append(*toolCallsState, make([]*streamToolCall, idx-len(*toolCallsState)+1)...)
		}
		st := (*toolCallsState)[idx]
		if st == nil {
			st = &streamToolCall{call: ToolCall{}}
			(*toolCallsState)[idx] = st
		}
		if tc.ID != "" {
			st.call.ID = tc.ID
		}
		if tc.Type != "" {
			st.call.Type = tc.Type
		}
		if tc.Function.Name != "" {
			st.call.Function.Name = tc.Function.Name
		}
		if tc.Function.Arguments != "" {
			st.args.WriteString(tc.Function.Arguments)
		}
	}
}

// processContentText handles the content text, tracking think blocks
func (agent *Agent) processContentText(text string, content *strings.Builder,
	reasoning *strings.Builder, inThink *bool) error {
	if *inThink {
		if idx := strings.Index(text, "</think"); idx >= 0 {
			if idx > 0 {
				reasoning.WriteString(text[:idx])
				agent.emit("think", text[:idx])
			}
			*inThink = false
			agent.flushThink()
			return nil
		}
		if text == "" {
			return nil
		}
		reasoning.WriteString(text)
		agent.emit("think", text)
	} else if idx := strings.Index(text, "<think"); idx >= 0 {
		if idx > 0 {
			content.WriteString(text[:idx])
			agent.emit("output", text[:idx])
		}
		*inThink = true
		if gtIdx := strings.Index(text[idx:], ">"); gtIdx >= 0 {
			after := text[idx+gtIdx+1:]
			if after != "" {
				reasoning.WriteString(after)
				agent.emit("think", after)
			}
		}
	} else if text != "" {
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
		Reasoning: truncate(reasoning),
		ToolCalls: toolCallsList,
	}
	agent.history = append(agent.history, assistantMsg)
	agent.logAssistantToolCalls(assistantMsg)

	// Process each tool call
	for _, tc := range toolCallsList {
		agent.executeSingleToolCall(ctx, tc)
	}
}

// approvalFnKey is the context key for injecting an approval function into tool handlers.
type approvalFnKey struct{}

// requireApproval gets the approval function from context and calls it.
// before is the original code (empty for create), after is the new code.
// It panics if no approval function is found, or returns an error if denied.
func requireApproval(ctx context.Context, toolName, before, after string) error {
	approvalFn, ok := ctx.Value(approvalFnKey{}).(func(before, after string) (bool, error))
	if !ok {
		panic(toolName + ": missing approval function")
	}
	allowed, err := approvalFn(before, after)
	if err != nil {
		return err
	}
	if !allowed {
		return fmt.Errorf("DENIED")
	}
	return nil
}

// executeSingleToolCall executes a single tool call and logs the result
func (agent *Agent) executeSingleToolCall(ctx context.Context, tc ToolCall) {
	name := strings.TrimPrefix(tc.Function.Name, "suneido_")
	toolOutput, err := agent.toolClient.FormatToolCallForDisplay(tc)
	if err != nil {
		toolOutput = name + " " + tc.Function.Arguments + "\n"
	}
	toolOutput = ensureTrailingNewlines(toolOutput, 2)
	var approvalText string
	if agent.needsApproval(tc.Function.Name) {
		agent.emit("tool", toolOutput)
		approval := newToolApproval()
		approvalFn := func(before, after string) (bool, error) {
			approval.Before = before
			approval.After = after
			agent.emitToolWithApproval("", approval)
			allowed, text, err := agent.waitForApproval(ctx, approval)
			approvalText = text
			if allowed {
				agent.emit("tool", "Allowed\n\n")
			} else {
				agent.emit("tool", "Denied\n\n")
			}
			return allowed, err
		}
		ctx = context.WithValue(ctx, approvalFnKey{}, approvalFn)
	} else {
		agent.emit("tool", toolOutput)
	}
	result, err := agent.toolClient.CallToolFromLLM(ctx, tc)
	if err != nil {
		agent.emit("tool", "ERROR: "+err.Error()+"\n\n")
		result = "ERROR: " + err.Error()
	} else if tc.Function.Name == "suneido_execute" {
		agent.emitExecToolResult(result)
	}
	if approvalText != "" {
		result += "\nUSER FEEDBACK: " + approvalText
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

func (agent *Agent) waitForApproval(ctx context.Context, approval *ToolApproval) (bool, string, error) {
	decision, err := approval.Wait(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return false, "", errors.New("tool approval canceled")
		}
		return false, "", err
	}
	if !decision.allow {
		return false, decision.text, nil
	}
	return true, decision.text, nil
}

func (*Agent) needsApproval(toolName string) bool {
	return strings.HasPrefix(toolName, "suneido_create_") ||
		strings.HasPrefix(toolName, "suneido_delete_") ||
		strings.HasPrefix(toolName, "suneido_edit_")
}

func (agent *Agent) emitExecToolResult(result string) {
	var execOut execOutput
	if err := json.Unmarshal([]byte(result), &execOut); err != nil {
		agent.emit("tool", "=> "+result+"\n")
		return
	}
	results := execOut.Results
	agent.emit("tool", "=> "+results+"\n")
	if execOut.Print != "" {
		agent.emit("tool", execOut.Print+"\n")
	}
}

func (agent *Agent) emit(what, data string) {
	if what == "think" {
		agent.thinkBuf.WriteString(data)
		agent.thinkDirty = true
		agent.outfn(what, data, nil)
		return
	}
	agent.flushThink()
	data = truncate(data)
	if what == "tool" {
		data = capTrailingNewlines(data, 2)
	}
	agent.outfn(what, data, nil)
}

func (agent *Agent) emitToolWithApproval(data string, approval *ToolApproval) {
	agent.flushThink()
	agent.outfn("tool", capTrailingNewlines(truncate(data), 2), approval)
}

// capTrailingNewlines limits trailing newlines to at most n,
// normalizing CRLF to LF. Returns the original string when no change is needed.
func capTrailingNewlines(s string, n int) string {
	if n < 0 {
		n = 0
	}
	i := len(s) - 1
	count := 0
	allLF := true
	for i >= 0 {
		switch s[i] {
		case '\n':
			count++
			i--
			if i >= 0 && s[i] == '\r' {
				allLF = false
				i--
			}
		case '\r':
			allLF = false
			count++
			i--
		default:
			goto done
		}
	}
done:
	if count <= n && allLF {
		return s
	}
	base := s[:i+1]
	if count > n {
		count = n
	}
	if count == 0 {
		return base
	}
	return base + strings.Repeat("\n", count)
}

// ensureTrailingNewlines returns s ending with exactly n newlines.
// Returns the original string when no change is needed.
func ensureTrailingNewlines(s string, n int) string {
	capped := capTrailingNewlines(s, n)
	// Count trailing newlines in capped result
	count := 0
	for i := len(capped) - 1; i >= 0 && capped[i] == '\n'; i-- {
		count++
	}
	if count == n {
		return capped
	}
	return capped + strings.Repeat("\n", n-count)
}
