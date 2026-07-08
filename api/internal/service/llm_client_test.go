package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestAnthropicClient returns a client pointed at a throwaway httptest
// server driven by handler.
func newTestAnthropicClient(t *testing.T, handler http.HandlerFunc) *AnthropicClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := NewAnthropicClient("test-key", "claude-test")
	client.baseURL = server.URL
	return client
}

func TestAnthropicClient_Chat_Success(t *testing.T) {
	var captured anthropicRequest
	client := newTestAnthropicClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-key", r.Header.Get("x-api-key"))
		assert.Equal(t, "2023-06-01", r.Header.Get("anthropic-version"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(body, &captured))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"content":[{"type":"text","text":"こんにちは"}]}`))
	})

	reply, err := client.Chat(context.Background(), "system prompt", []ChatMessage{{Role: "user", Content: "hi"}})

	require.NoError(t, err)
	assert.Equal(t, "こんにちは", reply)
	assert.Equal(t, "claude-test", captured.Model)
	assert.Equal(t, "system prompt", captured.System)
	require.Len(t, captured.Messages, 1)
	assert.Equal(t, "user", captured.Messages[0].Role)
	assert.Equal(t, "hi", captured.Messages[0].Content)
}

func TestAnthropicClient_Chat_StatusErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    error
	}{
		{"429 maps to ErrLLMRateLimit", http.StatusTooManyRequests, ErrLLMRateLimit},
		{"503 maps to ErrLLMUnavailable", http.StatusServiceUnavailable, ErrLLMUnavailable},
		{"502 maps to ErrLLMUnavailable", http.StatusBadGateway, ErrLLMUnavailable},
		{"504 maps to ErrLLMUnavailable", http.StatusGatewayTimeout, ErrLLMUnavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestAnthropicClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(`{"error":{"type":"error","message":"boom"}}`))
			})

			_, err := client.Chat(context.Background(), "sys", []ChatMessage{{Role: "user", Content: "hi"}})

			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAnthropicClient_Chat_GenericServerError(t *testing.T) {
	client := newTestAnthropicClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"type":"error","message":"boom"}}`))
	})

	_, err := client.Chat(context.Background(), "sys", []ChatMessage{{Role: "user", Content: "hi"}})

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrLLMRateLimit)
	assert.NotErrorIs(t, err, ErrLLMUnavailable)
	assert.Contains(t, err.Error(), "500")
}

func TestAnthropicClient_Chat_ClientError(t *testing.T) {
	client := newTestAnthropicClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"type":"invalid_request","message":"bad"}}`))
	})

	_, err := client.Chat(context.Background(), "sys", []ChatMessage{{Role: "user", Content: "hi"}})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
}

func TestAnthropicClient_Chat_MalformedResponseBody(t *testing.T) {
	client := newTestAnthropicClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{not-json`))
	})

	_, err := client.Chat(context.Background(), "sys", []ChatMessage{{Role: "user", Content: "hi"}})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal response")
}

func TestAnthropicClient_Chat_EmptyContent(t *testing.T) {
	client := newTestAnthropicClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"content":[]}`))
	})

	_, err := client.Chat(context.Background(), "sys", []ChatMessage{{Role: "user", Content: "hi"}})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

func TestAnthropicClient_Chat_TransportError(t *testing.T) {
	// Start a server, capture its URL, then close it so the connection is
	// refused and the HTTP round trip fails.
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	url := server.URL
	server.Close()

	client := NewAnthropicClient("test-key", "claude-test")
	client.baseURL = url

	_, err := client.Chat(context.Background(), "sys", []ChatMessage{{Role: "user", Content: "hi"}})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "send request")
}
