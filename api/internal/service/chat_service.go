package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/akaitigo/kotowaza-bridge/api/internal/domain/kotowaza"
	"github.com/akaitigo/kotowaza-bridge/api/internal/repository"
	"github.com/google/uuid"
)

// ErrValidation indicates a client input validation error.
var ErrValidation = errors.New("validation error")

// ErrNotFound indicates a resource was not found.
var ErrNotFound = errors.New("not found")

// ChatMessage represents a single chat message.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat request from the client.
type ChatRequest struct {
	KotowazaID uuid.UUID     `json:"kotowaza_id"`
	Messages   []ChatMessage `json:"messages"`
}

// ChatResponse represents a chat response to the client.
type ChatResponse struct {
	Message ChatMessage `json:"message"`
}

// LLMClient defines the interface for LLM API calls.
type LLMClient interface {
	Chat(ctx context.Context, systemPrompt string, messages []ChatMessage) (string, error)
}

// ChatService provides LLM-powered roleplay practice.
type ChatService struct {
	repo      repository.KotowazaRepository
	llmClient LLMClient
}

// NewChatService creates a new ChatService.
func NewChatService(repo repository.KotowazaRepository, llmClient LLMClient) *ChatService {
	return &ChatService{repo: repo, llmClient: llmClient}
}

const (
	maxMessages      = 50
	maxMessageLength = 2000
)

var validRoles = map[string]bool{"user": true, "assistant": true}

// Chat sends a chat message and returns the LLM response.
func (s *ChatService) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("%w: messages must not be empty", ErrValidation)
	}
	if len(req.Messages) > maxMessages {
		return nil, fmt.Errorf("%w: too many messages (max %d)", ErrValidation, maxMessages)
	}
	for _, m := range req.Messages {
		if !validRoles[m.Role] {
			return nil, fmt.Errorf("%w: invalid role: %s", ErrValidation, m.Role)
		}
		if len(m.Content) > maxMessageLength {
			return nil, fmt.Errorf("%w: message too long (max %d characters)", ErrValidation, maxMessageLength)
		}
	}

	k, err := s.repo.GetByID(ctx, req.KotowazaID)
	if err != nil {
		return nil, fmt.Errorf("get kotowaza for chat: %w", err)
	}
	if k == nil {
		return nil, fmt.Errorf("%w: kotowaza %s", ErrNotFound, req.KotowazaID)
	}

	systemPrompt := buildSystemPrompt(k)

	reply, err := s.llmClient.Chat(ctx, systemPrompt, req.Messages)
	if err != nil {
		return nil, fmt.Errorf("llm chat: %w", err)
	}

	return &ChatResponse{
		Message: ChatMessage{
			Role:    "assistant",
			Content: reply,
		},
	}, nil
}

func buildSystemPrompt(k *kotowaza.Kotowaza) string {
	return fmt.Sprintf(`あなたは日本語のことわざの使い方を教える先生です。
以下のことわざを日常会話で自然に使えるようになるための練習相手になってください。

ことわざ: %s（%s）
意味: %s
使用例: %s

ルール:
- 日常的なシナリオを提示して、ユーザーがことわざを使う機会を作ってください
- ユーザーがことわざを正しく使えたら褒めてください
- 不自然な使い方があれば、優しく訂正してください
- 会話は日本語で行ってください`,
		k.Japanese, k.Reading, k.Meaning, k.UsageExample)
}
