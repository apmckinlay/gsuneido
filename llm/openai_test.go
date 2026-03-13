// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNewOpenAIClient(t *testing.T) {
	client := NewOpenAIClient("https://api.example.com", "test-key")
	assert.T(t).This(client.BaseURL).Is("https://api.example.com")
	assert.T(t).This(client.APIKey).Is("test-key")
	assert.T(t).True(client.HTTPClient != nil)
	assert.T(t).This(client.HTTPClient.Timeout).Is(60 * time.Second)
}

func TestOpenAIClientStreamSuccess(t *testing.T) {
	// Create test server that returns SSE stream
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.T(t).This(r.URL.Path).Is("/chat/completions")
		assert.T(t).This(r.Header.Get("Content-Type")).Is("application/json")
		assert.T(t).This(r.Header.Get("Authorization")).Is("Bearer test-key")

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Send SSE chunks
		chunks := []string{
			`{"id":"1","choices":[{"delta":{"role":"assistant"}}]}`,
			`{"id":"2","choices":[{"delta":{"content":"Hello"}}]}`,
			`{"id":"3","choices":[{"delta":{"content":" world"}}]}`,
			`{"id":"4","choices":[{"delta":{},"finish_reason":"stop"}]}`,
		}
		for _, chunk := range chunks {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{
		Model:    "test-model",
		Messages: []Message{{Role: "user", Content: "Hi"}},
	}

	var received []string
	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			received = append(received, chunk.Choices[0].Delta.Content)
		}
		return nil
	})

	assert.T(t).This(err).Is(nil)
	assert.T(t).This(strings.Join(received, "")).Is("Hello world")
}

func TestOpenAIClientStreamWithReasoning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		chunks := []string{
			`{"id":"1","choices":[{"delta":{"reasoning_content":"thinking..."}}]}`,
			`{"id":"2","choices":[{"delta":{"content":"answer"}}]}`,
		}
		for _, chunk := range chunks {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	var reasoning, content string
	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		if len(chunk.Choices) > 0 {
			reasoning += chunk.Choices[0].Delta.ReasoningContent
			content += chunk.Choices[0].Delta.Content
		}
		return nil
	})

	assert.T(t).This(err).Is(nil)
	assert.T(t).This(reasoning).Is("thinking...")
	assert.T(t).This(content).Is("answer")
}

func TestOpenAIClientStreamHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":{"message":"internal error"}}`)
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).True(err != nil)
	assert.T(t).True(strings.Contains(err.Error(), "api error"))
	assert.T(t).True(strings.Contains(err.Error(), "internal error"))
}

func TestOpenAIClientStreamHTTPErrorNoJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "Bad Gateway")
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).True(err != nil)
	assert.T(t).True(strings.Contains(err.Error(), "http error"))
	assert.T(t).True(strings.Contains(err.Error(), "Bad Gateway"))
}

func TestOpenAIClientStreamChunkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "data: %s\n\n", `{"id":"1","error":{"message":"rate limited"}}`)
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).True(err != nil)
	assert.T(t).True(strings.Contains(err.Error(), "api error"))
	assert.T(t).True(strings.Contains(err.Error(), "rate limited"))
}

func TestOpenAIClientStreamCallbackError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "data: %s\n\n", `{"id":"1","choices":[{"delta":{"content":"test"}}]}`)
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	expectedErr := fmt.Errorf("callback error")
	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return expectedErr
	})

	assert.T(t).This(err).Is(expectedErr)
}

func TestOpenAIClientStreamInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, "data: invalid json\n\n")
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).True(err != nil)
	assert.T(t).True(strings.Contains(err.Error(), "unmarshal chunk"))
}

func TestOpenAIClientStreamNoDataChunks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		// No data sent
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).True(err != nil)
	assert.T(t).True(strings.Contains(err.Error(), "missing data chunks"))
}

func TestOpenAIClientStreamNonSSEError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"error":{"message":"something went wrong"}}`)
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).True(err != nil)
	assert.T(t).True(strings.Contains(err.Error(), "api error"))
	assert.T(t).True(strings.Contains(err.Error(), "something went wrong"))
}

func TestOpenAIClientStreamConnectionError(t *testing.T) {
	// Use an invalid URL to trigger a connection error
	client := NewOpenAIClient("http://127.0.0.1:1", "test-key") // Invalid port
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).True(err != nil)
	assert.T(t).True(strings.Contains(err.Error(), "send request"))
}

func TestOpenAIClientStreamSetsStreamFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify stream flag is set
		var req ChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.T(t).This(err).Is(nil)
		assert.T(t).True(req.Stream)

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).This(err).Is(nil)
}

func TestOpenAIClientRawLog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "data: %s\n\n", `{"id":"1","choices":[{"delta":{"content":"test"}}]}`)
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	var logBuf bytes.Buffer
	client := NewOpenAIClient(server.URL, "test-key")
	client.RawLog = &logBuf

	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).This(err).Is(nil)
	logOutput := logBuf.String()
	assert.T(t).True(strings.Contains(logOutput, "========================================"))
	assert.T(t).True(strings.Contains(logOutput, "----------------------------------------"))
	assert.T(t).True(strings.Contains(logOutput, "test-model"))
}

func TestOpenAIClientRawLogHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":{"message":"test error"}}`)
	}))
	defer server.Close()

	var logBuf bytes.Buffer
	client := NewOpenAIClient(server.URL, "test-key")
	client.RawLog = &logBuf

	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		return nil
	})

	assert.T(t).True(err != nil)
	logOutput := logBuf.String()
	assert.T(t).True(strings.Contains(logOutput, "500"))
	assert.T(t).True(strings.Contains(logOutput, "test error"))
}

func TestOpenAIClientStreamWithToolCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		chunks := []string{
			`{"id":"1","choices":[{"delta":{"tool_calls":[{"id":"call_123","type":"function","function":{"name":"get_weather","arguments":"{\"loc"}}]}}]}`,
			`{"id":"2","choices":[{"delta":{"tool_calls":[{"function":{"arguments":"ation\":"}}]}}]}`,
			`{"id":"3","choices":[{"delta":{"tool_calls":[{"function":{"arguments":"\"NYC\"}"}}]}}]}`,
		}
		for _, chunk := range chunks {
			fmt.Fprintf(w, "data: %s\n\n", chunk)
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Weather?"}}}

	var toolCalls []ToolCall
	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		if len(chunk.Choices) > 0 && len(chunk.Choices[0].Delta.ToolCalls) > 0 {
			toolCalls = append(toolCalls, chunk.Choices[0].Delta.ToolCalls...)
		}
		return nil
	})

	assert.T(t).This(err).Is(nil)
	assert.T(t).This(len(toolCalls)).Is(3)
	assert.T(t).This(toolCalls[0].ID).Is("call_123")
	assert.T(t).This(toolCalls[0].Function.Name).Is("get_weather")
}

func TestOpenAIClientStreamEmptyLinesIgnored(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, "\n\n")
		fmt.Fprintf(w, "data: %s\n\n", `{"id":"1","choices":[{"delta":{"content":"test"}}]}`)
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	var content string
	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		if len(chunk.Choices) > 0 {
			content += chunk.Choices[0].Delta.Content
		}
		return nil
	})

	assert.T(t).This(err).Is(nil)
	assert.T(t).This(content).Is("test")
}

func TestOpenAIClientStreamWhitespaceLinesIgnored(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, "   \n")
		fmt.Fprintf(w, "data: %s\n\n", `{"id":"1","choices":[{"delta":{"content":"test"}}]}`)
		fmt.Fprint(w, "\t\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	req := &ChatRequest{Model: "test-model", Messages: []Message{{Role: "user", Content: "Hi"}}}

	var content string
	err := client.Stream(context.Background(), req, func(chunk *ChatCompletionChunk) error {
		if len(chunk.Choices) > 0 {
			content += chunk.Choices[0].Delta.Content
		}
		return nil
	})

	assert.T(t).This(err).Is(nil)
	assert.T(t).This(content).Is("test")
}
