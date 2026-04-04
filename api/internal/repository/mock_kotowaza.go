package repository

import (
	"context"

	"github.com/akaitigo/kotowaza-bridge/api/internal/domain/kotowaza"
	"github.com/google/uuid"
)

// MockKotowazaRepository is a mock implementation of KotowazaRepository for testing.
type MockKotowazaRepository struct {
	ListFn   func(ctx context.Context, params kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error)
	GetByIDFn func(ctx context.Context, id uuid.UUID) (*kotowaza.Kotowaza, error)
	SearchFn func(ctx context.Context, params kotowaza.SearchParams) ([]kotowaza.Kotowaza, int, error)
}

func (m *MockKotowazaRepository) List(ctx context.Context, params kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error) {
	return m.ListFn(ctx, params)
}

func (m *MockKotowazaRepository) GetByID(ctx context.Context, id uuid.UUID) (*kotowaza.Kotowaza, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockKotowazaRepository) Search(ctx context.Context, params kotowaza.SearchParams) ([]kotowaza.Kotowaza, int, error) {
	return m.SearchFn(ctx, params)
}
