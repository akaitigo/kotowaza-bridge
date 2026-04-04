package handler

import (
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

func setupKotowazaHandler(mock *repository.MockKotowazaRepository) (*KotowazaHandler, *chi.Mux) {
	svc := service.NewKotowazaService(mock)
	h := NewKotowazaHandler(svc)
	r := chi.NewRouter()
	r.Get("/api/v1/kotowaza", h.List)
	r.Get("/api/v1/kotowaza/search", h.Search)
	r.Get("/api/v1/kotowaza/{id}", h.GetByID)
	return h, r
}

func TestKotowazaHandler_List(t *testing.T) {
	t.Run("returns 200 with kotowaza list", func(t *testing.T) {
		items := []kotowaza.Kotowaza{
			{ID: uuid.New(), Japanese: "猿も木から落ちる", Reading: "さるもきからおちる", Meaning: "テスト", CreatedAt: time.Now()},
		}
		mock := &repository.MockKotowazaRepository{
			ListFn: func(_ context.Context, _ kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error) {
				return items, 1, nil
			},
		}
		_, r := setupKotowazaHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza?limit=10&offset=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp service.ListResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, 1, len(resp.Items))
		assert.Equal(t, "猿も木から落ちる", resp.Items[0].Japanese)
	})
}

func TestKotowazaHandler_GetByID(t *testing.T) {
	t.Run("returns 200 with kotowaza detail", func(t *testing.T) {
		id := uuid.New()
		expected := &kotowaza.Kotowaza{
			ID:       id,
			Japanese: "七転び八起き",
			Reading:  "ななころびやおき",
			Meaning:  "テスト",
			Equivalents: []kotowaza.Equivalent{
				{ID: uuid.New(), KotowazaID: id, Language: "en", Expression: "Fall seven times, stand up eight"},
			},
			CreatedAt: time.Now(),
		}
		mock := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, reqID uuid.UUID) (*kotowaza.Kotowaza, error) {
				assert.Equal(t, id, reqID)
				return expected, nil
			},
		}
		_, r := setupKotowazaHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza/"+id.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp kotowaza.Kotowaza
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "七転び八起き", resp.Japanese)
		assert.Equal(t, 1, len(resp.Equivalents))
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		mock := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return nil, nil
			},
		}
		_, r := setupKotowazaHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza/"+uuid.New().String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for invalid UUID", func(t *testing.T) {
		mock := &repository.MockKotowazaRepository{}
		_, r := setupKotowazaHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza/not-a-uuid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestKotowazaHandler_Search(t *testing.T) {
	t.Run("returns 200 with search results", func(t *testing.T) {
		items := []kotowaza.Kotowaza{
			{ID: uuid.New(), Japanese: "猿も木から落ちる", Reading: "さるもきからおちる", Meaning: "テスト", CreatedAt: time.Now()},
		}
		mock := &repository.MockKotowazaRepository{
			SearchFn: func(_ context.Context, params kotowaza.SearchParams) ([]kotowaza.Kotowaza, int, error) {
				assert.Equal(t, "猿", params.Query)
				return items, 1, nil
			},
		}
		_, r := setupKotowazaHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza/search?q=猿", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("returns 400 when query is empty", func(t *testing.T) {
		mock := &repository.MockKotowazaRepository{}
		_, r := setupKotowazaHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/kotowaza/search?q=", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
