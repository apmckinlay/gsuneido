// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAIClient implements Client for OpenAI-compatible APIs.
type OpenAIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	RawLog     io.Writer // optional raw HTTP payload logging
}

// NewOpenAIClient creates a new OpenAI-compatible client.
func NewOpenAIClient(baseURL, apiKey string) *OpenAIClient {
	return &OpenAIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Chat sends a chat request to an OpenAI-compatible API.
func (c *OpenAIClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	c.logRawRequest(body)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.logRawResponse(resp.Status, respBody)
		var chatResp ChatResponse
		if err := json.Unmarshal(respBody, &chatResp); err == nil && chatResp.Error != nil {
			return nil, fmt.Errorf("api error: %s", chatResp.Error.Message)
		}
		return nil, fmt.Errorf("http error: %s: %s", resp.Status,
			strings.TrimSpace(string(respBody)))
	}
	c.logRawResponse("", respBody)

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("api error: %s", chatResp.Error.Message)
	}

	return &chatResp, nil
}

// StreamChunk is a callback function for streaming responses.
type StreamChunk func(chunk *ChatCompletionChunk) error

// Stream sends a streaming chat completion request to an OpenAI-compatible API.
func (c *OpenAIClient) Stream(ctx context.Context, req *ChatRequest, onChunk StreamChunk) error {
	req.Stream = true
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	c.logRawRequest(body)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions",
		bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Streaming responses can legitimately run for several minutes.
	// Use request context cancellation, not http.Client.Timeout,
	// to avoid aborting long streams while reading the body.
	streamClient := *c.HTTPClient
	streamClient.Timeout = 0
	resp, err := streamClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		c.logRawResponse(resp.Status, body)
		var chatResp ChatResponse
		if err := json.Unmarshal(body, &chatResp); err == nil && chatResp.Error != nil {
			return fmt.Errorf("api error: %s", chatResp.Error.Message)
		}
		return fmt.Errorf("http error: %s: %s", resp.Status,
			strings.TrimSpace(string(body)))
	}

	responseBody := io.Reader(resp.Body)
	if c.RawLog != nil {
		_, _ = io.WriteString(c.RawLog, "\n----------------------------------------\n")
		responseBody = io.TeeReader(resp.Body, c.RawLog)
	}

	scanner := bufio.NewScanner(responseBody)
	seenData := false
	var nonSSE strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(line, "data: ") {
			nonSSE.WriteString(trimmed)
			continue
		}
		seenData = true

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk ChatCompletionChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf("unmarshal chunk: %w", err)
		}

		if chunk.Error != nil {
			return fmt.Errorf("api error: %s", chunk.Error.Message)
		}

		if err := onChunk(&chunk); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !seenData {
		if nonSSE.Len() > 0 {
			raw := nonSSE.String()
			var chatResp ChatResponse
			if err := json.Unmarshal([]byte(raw), &chatResp); err == nil && chatResp.Error != nil {
				return fmt.Errorf("api error: %s", chatResp.Error.Message)
			}
			return fmt.Errorf("stream response missing data chunks: %s", raw)
		}
		return fmt.Errorf("stream response missing data chunks")
	}

	return nil
}

func (c *OpenAIClient) logRawRequest(body []byte) {
	if c.RawLog == nil {
		return
	}
	_, _ = io.WriteString(c.RawLog, "\n========================================\n")
	_, _ = c.RawLog.Write(body)
	_, _ = io.WriteString(c.RawLog, "\n")
}

func (c *OpenAIClient) logRawResponse(status string, body []byte) {
	if c.RawLog == nil {
		return
	}
	_, _ = io.WriteString(c.RawLog, "\n----------------------------------------\n")
	if status != "" {
		_, _ = io.WriteString(c.RawLog, status+"\n")
	}
	_, _ = c.RawLog.Write(body)
	_, _ = io.WriteString(c.RawLog, "\n")
}
