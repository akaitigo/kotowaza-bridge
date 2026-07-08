package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/akaitigo/kotowaza-bridge/api/internal/service"
)

// ChatServicer is the subset of the chat service the handler depends on.
// Depending on the interface rather than the concrete *service.ChatService lets
// tests inject a mock and exercise error paths in isolation.
type ChatServicer interface {
	Chat(ctx context.Context, req service.ChatRequest) (*service.ChatResponse, error)
}

// ChatHandler handles HTTP requests for chat endpoints.
type ChatHandler struct {
	svc ChatServicer
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(svc ChatServicer) *ChatHandler {
	return &ChatHandler{svc: svc}
}

// Chat handles POST /api/v1/chat
func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	const maxBodySize = 1 << 16 // 64KB
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req service.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.svc.Chat(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrLLMRateLimit):
			writeError(w, http.StatusTooManyRequests, "LLM service is busy, please try again later")
		case errors.Is(err, service.ErrLLMUnavailable):
			writeError(w, http.StatusServiceUnavailable, "LLM service is temporarily unavailable")
		case errors.Is(err, context.DeadlineExceeded):
			writeError(w, http.StatusGatewayTimeout, "request timed out")
		case errors.Is(err, context.Canceled):
			// Client disconnected; no response needed but log it.
			writeError(w, http.StatusBadRequest, "request canceled")
		default:
			writeError(w, http.StatusInternalServerError, "chat failed")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
