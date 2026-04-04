package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/kotowaza-bridge/api/internal/service"
)

// ChatHandler handles HTTP requests for chat endpoints.
type ChatHandler struct {
	svc *service.ChatService
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(svc *service.ChatService) *ChatHandler {
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
		writeError(w, http.StatusInternalServerError, "chat failed")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
