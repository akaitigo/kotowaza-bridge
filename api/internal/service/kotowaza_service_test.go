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

func newTestKotowaza(japanese string) kotowaza.Kotowaza {
	return kotowaza.Kotowaza{
		ID:           uuid.New(),
		Japanese:     japanese,
		Reading:      "てすと",
		Meaning:      "テスト用の意味",
		Origin:       "テスト由来",
		UsageExample: "テスト使用例",
		CulturalNote: "テスト文化背景",
		CreatedAt:    time.Now(),
	}
}

func TestKotowazaService_List(t *testing.T) {
	t.Run("returns paginated results with default limit", func(t *testing.T) {
		items := []kotowaza.Kotowaza{newTestKotowaza("テスト1"), newTestKotowaza("テスト2")}
		mock := &repository.MockKotowazaRepository{
			ListFn: func(_ context.Context, params kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error) {
				assert.Equal(t, 20, params.Limit)
				assert.Equal(t, 0, params.Offset)
				return items, 2, nil
			},
		}
		svc := NewKotowazaService(mock)

		resp, err := svc.List(context.Background(), 0, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, len(resp.Items))
		assert.Equal(t, 2, resp.Total)
		assert.Equal(t, 20, resp.Limit)
	})

	t.Run("clamps limit to max 100", func(t *testing.T) {
		mock := &repository.MockKotowazaRepository{
			ListFn: func(_ context.Context, params kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error) {
				assert.Equal(t, 20, params.Limit)
				return nil, 0, nil
			},
		}
		svc := NewKotowazaService(mock)

		_, err := svc.List(context.Background(), 200, 0)
		require.NoError(t, err)
	})

	t.Run("normalizes negative offset to 0", func(t *testing.T) {
		mock := &repository.MockKotowazaRepository{
			ListFn: func(_ context.Context, params kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error) {
				assert.Equal(t, 0, params.Offset)
				return nil, 0, nil
			},
		}
		svc := NewKotowazaService(mock)

		_, err := svc.List(context.Background(), 10, -5)
		require.NoError(t, err)
	})

	t.Run("propagates repository error", func(t *testing.T) {
		mock := &repository.MockKotowazaRepository{
			ListFn: func(_ context.Context, _ kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error) {
				return nil, 0, errors.New("db error")
			},
		}
		svc := NewKotowazaService(mock)

		_, err := svc.List(context.Background(), 10, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestKotowazaService_GetByID(t *testing.T) {
	t.Run("returns kotowaza when found", func(t *testing.T) {
		expected := newTestKotowaza("猿も木から落ちる")
		mock := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, id uuid.UUID) (*kotowaza.Kotowaza, error) {
				assert.Equal(t, expected.ID, id)
				return &expected, nil
			},
		}
		svc := NewKotowazaService(mock)

		result, err := svc.GetByID(context.Background(), expected.ID)
		require.NoError(t, err)
		assert.Equal(t, expected.Japanese, result.Japanese)
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		mock := &repository.MockKotowazaRepository{
			GetByIDFn: func(_ context.Context, _ uuid.UUID) (*kotowaza.Kotowaza, error) {
				return nil, nil
			},
		}
		svc := NewKotowazaService(mock)

		result, err := svc.GetByID(context.Background(), uuid.New())
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestKotowazaService_Search(t *testing.T) {
	t.Run("passes query and params to repository", func(t *testing.T) {
		items := []kotowaza.Kotowaza{newTestKotowaza("猿も木から落ちる")}
		mock := &repository.MockKotowazaRepository{
			SearchFn: func(_ context.Context, params kotowaza.SearchParams) ([]kotowaza.Kotowaza, int, error) {
				assert.Equal(t, "猿", params.Query)
				assert.Equal(t, 10, params.Limit)
				assert.Equal(t, 0, params.Offset)
				return items, 1, nil
			},
		}
		svc := NewKotowazaService(mock)

		resp, err := svc.Search(context.Background(), "猿", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(resp.Items))
		assert.Equal(t, 1, resp.Total)
	})
}
