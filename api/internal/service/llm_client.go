package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// anthropicRequest represents the Anthropic Messages API request body.
type anthropicRequest struct {
	Model     string              `json:"model"`
	MaxTokens int                 `json:"max_tokens"`
	System    string              `json:"system"`
	Messages  []anthropicMessage  `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []anthropicContent `json:"content"`
	Error   *anthropicError    `json:"error,omitempty"`
}

type anthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// AnthropicClient implements LLMClient using the Anthropic Messages API.
type AnthropicClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewAnthropicClient creates a new AnthropicClient.
func NewAnthropicClient(apiKey, model string) *AnthropicClient {
	return &AnthropicClient{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Chat sends messages to the Anthropic API and returns the response.
func (c *AnthropicClient) Chat(ctx context.Context, systemPrompt string, messages []ChatMessage) (string, error) {
	anthropicMsgs := make([]anthropicMessage, 0, len(messages))
	for _, m := range messages {
		anthropicMsgs = append(anthropicMsgs, anthropicMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	reqBody := anthropicRequest{
		Model:     c.model,
		MaxTokens: 1024,
		System:    systemPrompt,
		Messages:  anthropicMsgs,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic API error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var anthropicResp anthropicResponse
	if err := json.Unmarshal(respBytes, &anthropicResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("empty response from anthropic API")
	}

	return anthropicResp.Content[0].Text, nil
}
