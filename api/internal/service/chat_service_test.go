package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/akaitigo/kotowaza-bridge/api/internal/domain/kotowaza"
	"github.com/akaitigo/kotowaza-bridge/api/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLLMClient struct {
	ChatFn func(ctx context.Context, systemPrompt string, messages []ChatMessage) (string, error)
}

func (m *mockLLMClient) Chat(ctx context.Context, systemPrompt string, messages []ChatMessage) (string, error) {
	return m.ChatFn(ctx, systemPrompt, messages)
}

func TestChatService_Chat(t *testing.T) {
	testKotowaza := &kotowaza.Kotowaza{
		ID:           uuid.MustParse("a1b2c3d4-0001-4000-8000-000000000001"),
		Japanese:     "猿も木から落ちる",
		Reading:      "さるもきからおちる",
		Meaning:      "得意なことでも失敗することがある",
		UsageExample: "テスト使用例",
		CreatedAt:    time.Now(),
	}

	t.Run("returns LLM response for valid request", func(t *testing.T) {
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, id uuid.UUID) (*kotowaza.Kotowaza, error) {
				assert.Equal(t, testKotowaza.ID, id)
				return testKotowaza, nil
			},
		}
		mockLLM := &mockLLMClient{
			ChatFn: func(_ context.Context, systemPrompt string, messages []ChatMessage) (string, error) {
				assert.Contains(t, systemPrompt, "猿も木から落ちる")
				assert.Equal(t, 1, len(messages))
				return "素晴らしい質問ですね！", nil
			},
		}
		svc := NewChatService(mockRepo, mockLLM)

		resp, err := svc.Chat(context.Background(), ChatRequest{
			KotowazaID: testKotowaza.ID,
			Messages:   []ChatMessage{{Role: "user", Content: "このことわざはどう使いますか？"}},
		})

		require.NoError(t, err)
		assert.Equal(t, "assistant", resp.Message.Role)
		assert.Equal(t, "素晴らしい質問ですね！", resp.Message.Content)
	})

	t.Run("returns error when kotowaza not found", func(t *testing.T) {
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return nil, nil
			},
		}
		mockLLM := &mockLLMClient{}
		svc := NewChatService(mockRepo, mockLLM)

		_, err := svc.Chat(context.Background(), ChatRequest{
			KotowazaID: uuid.New(),
			Messages:   []ChatMessage{{Role: "user", Content: "test"}},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("returns error when LLM fails", func(t *testing.T) {
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return testKotowaza, nil
			},
		}
		mockLLM := &mockLLMClient{
			ChatFn: func(_ context.Context, _ string, _ []ChatMessage) (string, error) {
				return "", errors.New("API rate limited")
			},
		}
		svc := NewChatService(mockRepo, mockLLM)

		_, err := svc.Chat(context.Background(), ChatRequest{
			KotowazaID: testKotowaza.ID,
			Messages:   []ChatMessage{{Role: "user", Content: "test"}},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "llm chat")
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		mockRepo := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return nil, errors.New("db connection error")
			},
		}
		mockLLM := &mockLLMClient{}
		svc := NewChatService(mockRepo, mockLLM)

		_, err := svc.Chat(context.Background(), ChatRequest{
			KotowazaID: uuid.New(),
			Messages:   []ChatMessage{{Role: "user", Content: "test"}},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "get kotowaza for chat")
	})
}

func TestChatService_Validation(t *testing.T) {
	mockRepo := &repository.MockKotowazaRepository{}
	mockLLM := &mockLLMClient{}
	svc := NewChatService(mockRepo, mockLLM)

	t.Run("rejects empty messages", func(t *testing.T) {
		_, err := svc.Chat(context.Background(), ChatRequest{
			KotowazaID: uuid.New(),
			Messages:   []ChatMessage{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must not be empty")
	})

	t.Run("rejects too many messages", func(t *testing.T) {
		msgs := make([]ChatMessage, 51)
		for i := range msgs {
			msgs[i] = ChatMessage{Role: "user", Content: "test"}
		}
		_, err := svc.Chat(context.Background(), ChatRequest{
			KotowazaID: uuid.New(),
			Messages:   msgs,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "too many messages")
	})

	t.Run("rejects invalid role", func(t *testing.T) {
		_, err := svc.Chat(context.Background(), ChatRequest{
			KotowazaID: uuid.New(),
			Messages:   []ChatMessage{{Role: "system", Content: "hack"}},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
	})

	t.Run("rejects message too long", func(t *testing.T) {
		longMsg := make([]byte, 2001)
		for i := range longMsg {
			longMsg[i] = 'a'
		}
		_, err := svc.Chat(context.Background(), ChatRequest{
			KotowazaID: uuid.New(),
			Messages:   []ChatMessage{{Role: "user", Content: string(longMsg)}},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "message too long")
	})
}

func TestBuildSystemPrompt(t *testing.T) {
	k := &kotowaza.Kotowaza{
		Japanese:     "猿も木から落ちる",
		Reading:      "さるもきからおちる",
		Meaning:      "得意なことでも失敗することがある",
		UsageExample: "彼はベテランだが猿も木から落ちる",
	}

	prompt := buildSystemPrompt(k)
	assert.Contains(t, prompt, "猿も木から落ちる")
	assert.Contains(t, prompt, "さるもきからおちる")
	assert.Contains(t, prompt, "得意なことでも失敗することがある")
	assert.Contains(t, prompt, "ベテランだが")
}
