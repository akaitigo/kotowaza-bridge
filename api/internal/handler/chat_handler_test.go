package handler

import (
	"bytes"
	"context"
	"encoding/json"
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
}
