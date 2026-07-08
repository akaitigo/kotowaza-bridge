package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akaitigo/kotowaza-bridge/api/internal/domain/kotowaza"
	"github.com/akaitigo/kotowaza-bridge/api/internal/repository"
	"github.com/akaitigo/kotowaza-bridge/api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLLMClient struct {
	response string
	err      error
}

func (c *testLLMClient) Chat(_ context.Context, _ string, _ []service.ChatMessage) (string, error) {
	return c.response, c.err
}

// mockChatServicer is a hand-rolled mock of the ChatServicer interface. Injecting
// it lets the handler tests drive error branches that are hard to trigger through
// the real service.
type mockChatServicer struct {
	resp *service.ChatResponse
	err  error
}

func (m *mockChatServicer) Chat(_ context.Context, _ service.ChatRequest) (*service.ChatResponse, error) {
	return m.resp, m.err
}

func setupChatHandler(mockRepo *repository.MockKotowazaRepository, llm service.LLMClient) *chi.Mux {
	chatSvc := service.NewChatService(mockRepo, llm)
	h := NewChatHandler(chatSvc)
	r := chi.NewRouter()
	r.Post("/api/v1/chat", h.Chat)
	return r
}

func TestChatHandler_Chat(t *testing.T) {
	testID := uuid.MustParse("a1b2c3d4-0001-4000-8000-000000000001")
	testKotowaza := &kotowaza.Kotowaza{
		ID:           testID,
		Japanese:     "猿も木から落ちる",
		Reading:      "さるもきからおちる",
		Meaning:      "テスト",
		UsageExample: "テスト",
		CreatedAt:    time.Now(),
	}

	t.Run("returns 200 with chat response", func(t *testing.T) {
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return testKotowaza, nil
			},
		}
		llm := &testLLMClient{response: "いい質問ですね！"}

		r := setupChatHandler(mockRepo, llm)

		body, _ := json.Marshal(service.ChatRequest{
			KotowazaID: testID,
			Messages:   []service.ChatMessage{{Role: "user", Content: "教えて"}},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp service.ChatResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "assistant", resp.Message.Role)
		assert.Equal(t, "いい質問ですね！", resp.Message.Content)
	})

	t.Run("returns 400 for invalid JSON body", func(t *testing.T) {
		mockRepo := &repository.MockKotowazaRepository{}
		llm := &testLLMClient{}

		r := setupChatHandler(mockRepo, llm)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 404 when kotowaza not found", func(t *testing.T) {
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return nil, nil
			},
		}
		llm := &testLLMClient{}

		r := setupChatHandler(mockRepo, llm)

		body, _ := json.Marshal(service.ChatRequest{
			KotowazaID: uuid.New(),
			Messages:   []service.ChatMessage{{Role: "user", Content: "test"}},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for validation error", func(t *testing.T) {
		mockRepo := &repository.MockKotowazaRepository{}
		llm := &testLLMClient{}

		r := setupChatHandler(mockRepo, llm)

		body, _ := json.Marshal(service.ChatRequest{
			KotowazaID: uuid.New(),
			Messages:   []service.ChatMessage{},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 429 when LLM rate limited", func(t *testing.T) {
		testID := uuid.MustParse("a1b2c3d4-0001-4000-8000-000000000001")
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return &kotowaza.Kotowaza{
					ID:           testID,
					Japanese:     "猿も木から落ちる",
					Reading:      "さるもきからおちる",
					Meaning:      "テスト",
					UsageExample: "テスト",
					CreatedAt:    time.Now(),
				}, nil
			},
		}
		llm := &testLLMClient{err: fmt.Errorf("%w: status 429: rate limited", service.ErrLLMRateLimit)}

		r := setupChatHandler(mockRepo, llm)

		body, _ := json.Marshal(service.ChatRequest{
			KotowazaID: testID,
			Messages:   []service.ChatMessage{{Role: "user", Content: "test"}},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("returns 503 when LLM unavailable", func(t *testing.T) {
		testID := uuid.MustParse("a1b2c3d4-0001-4000-8000-000000000001")
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return &kotowaza.Kotowaza{
					ID:           testID,
					Japanese:     "猿も木から落ちる",
					Reading:      "さるもきからおちる",
					Meaning:      "テスト",
					UsageExample: "テスト",
					CreatedAt:    time.Now(),
				}, nil
			},
		}
		llm := &testLLMClient{err: fmt.Errorf("%w: status 503: service unavailable", service.ErrLLMUnavailable)}

		r := setupChatHandler(mockRepo, llm)

		body, _ := json.Marshal(service.ChatRequest{
			KotowazaID: testID,
			Messages:   []service.ChatMessage{{Role: "user", Content: "test"}},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("returns 504 when context deadline exceeded", func(t *testing.T) {
		testID := uuid.MustParse("a1b2c3d4-0001-4000-8000-000000000001")
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return &kotowaza.Kotowaza{
					ID:           testID,
					Japanese:     "猿も木から落ちる",
					Reading:      "さるもきからおちる",
					Meaning:      "テスト",
					UsageExample: "テスト",
					CreatedAt:    time.Now(),
				}, nil
			},
		}
		llm := &testLLMClient{err: fmt.Errorf("llm chat: %w", context.DeadlineExceeded)}

		r := setupChatHandler(mockRepo, llm)

		body, _ := json.Marshal(service.ChatRequest{
			KotowazaID: testID,
			Messages:   []service.ChatMessage{{Role: "user", Content: "test"}},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusGatewayTimeout, w.Code)
	})
}

// TestChatHandler_ServiceErrors drives the handler through a mocked ChatServicer,
// covering branches that the concrete service cannot easily reach.
func TestChatHandler_ServiceErrors(t *testing.T) {
	testID := uuid.MustParse("a1b2c3d4-0001-4000-8000-000000000001")

	setup := func(svc ChatServicer) *chi.Mux {
		h := NewChatHandler(svc)
		r := chi.NewRouter()
		r.Post("/api/v1/chat", h.Chat)
		return r
	}

	newRequest := func() *http.Request {
		body, _ := json.Marshal(service.ChatRequest{
			KotowazaID: testID,
			Messages:   []service.ChatMessage{{Role: "user", Content: "hi"}},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	t.Run("maps an unexpected error to 500", func(t *testing.T) {
		r := setup(&mockChatServicer{err: errors.New("something broke")})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, newRequest())

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("maps context canceled to 400", func(t *testing.T) {
		r := setup(&mockChatServicer{err: fmt.Errorf("llm chat: %w", context.Canceled)})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, newRequest())

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 200 with the service response", func(t *testing.T) {
		r := setup(&mockChatServicer{resp: &service.ChatResponse{
			Message: service.ChatMessage{Role: "assistant", Content: "やあ"},
		}})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, newRequest())

		assert.Equal(t, http.StatusOK, w.Code)

		var resp service.ChatResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, "やあ", resp.Message.Content)
	})
}
