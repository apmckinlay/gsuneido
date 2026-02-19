// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"strings"
)

type Agent struct {
	client    *OpenAIClient
	mcpClient *MCPClient
	model     string
	history   []Message
	outfn     OutFn
	cancel    context.CancelFunc
}

// OutFn is the push callback for streaming output.
// what is one of "think", "output", "tool", "complete"
type OutFn func(what, data string)

// NewAgent creates an agent.
// prompt and mcpClient are optional.
func NewAgent(baseURL, apiKey, model, prompt string, mcpClient *MCPClient, outfn OutFn) *Agent {
	history := []Message{}
	if prompt != "" {
		history = append(history, Message{Role: "system", Content: prompt})
	}
	return &Agent{
		client:    NewOpenAIClient(baseURL, apiKey),
		mcpClient: mcpClient,
		model:     model,
		history:   history,
		outfn:     outfn,
	}
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
	agent.model = model
}

// ClearHistory clears the conversation history
func (agent *Agent) ClearHistory() {
	agent.history = []Message{}
}

// request sends the request and streams the response to outfn
func (agent *Agent) request(input string) {
	ctx, cancel := context.WithCancel(context.Background())
	agent.cancel = cancel
	defer cancel()

	agent.history = append(agent.history, Message{Role: "user", Content: input})

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
				idx := tc.Index
				if idx < 0 {
					idx = 0
				}
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
			agent.history = append(agent.history, Message{
				Role:      "assistant",
				Content:   content.String(),
				ToolCalls: toolCallsList,
			})

			// Process each tool call
			for _, tc := range toolCallsList {
				agent.emit("tool",
					"**"+tc.Function.Name+"** "+tc.Function.Arguments+"\n\n")
				result, err := agent.mcpClient.CallToolFromLLM(ctx, tc)
				if err != nil {
					agent.emit("tool", "\n\nError: "+err.Error())
				}
				// Add tool result to history
				agent.history = append(agent.history, Message{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
				})
			}

			// Continue the loop to get next response
			continue
		}

		// No tool calls, we're done
		agent.history = append(agent.history, Message{Role: "assistant", Content: content.String()})
		agent.emit("complete", "")
		return
	}
}

func (agent *Agent) emit(what, data string) {
	agent.outfn(what, data)
}
