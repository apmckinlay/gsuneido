// MIT License
//
// Copyright (c) 2025 gSuneido
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestOpenAIClientChat(t *testing.T) {
	mockResp := ChatResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   "gpt-4",
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: "Hello! How can I help you?",
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 8,
			TotalTokens:      18,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.T(t).This(r.URL.Path).Is("/chat/completions")
		assert.T(t).This(r.Header.Get("Authorization")).Is("Bearer test-key")
		assert.T(t).This(r.Header.Get("Content-Type")).Is("application/json")

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		assert.T(t).This(req.Model).Is("gpt-4")
		assert.T(t).This(len(req.Messages)).Is(1)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockResp); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "test-key")
	resp, err := client.Chat(context.Background(), &ChatRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	})
	assert.T(t).This(err).Is(nil)
	assert.T(t).This(resp.ID).Is("test-id")
	assert.T(t).This(resp.Model).Is("gpt-4")
	assert.T(t).This(len(resp.Choices)).Is(1)
	assert.T(t).This(resp.Choices[0].Message.Content).Is("Hello! How can I help you?")
	assert.T(t).This(resp.Usage.TotalTokens).Is(18)
}

func TestOpenAIClientChatError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(ChatResponse{
			Error: &Error{
				Message: "Invalid API key",
				Type:    "invalid_request_error",
				Code:    "invalid_api_key",
			},
		}); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := NewOpenAIClient(server.URL, "bad-key")
	_, err := client.Chat(context.Background(), &ChatRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	})
	assert.T(t).This(err).Isnt(nil)
}
