package repository

import (
	"context"

	"github.com/akaitigo/kotowaza-bridge/api/internal/domain/kotowaza"
	"github.com/google/uuid"
)

// KotowazaRepository defines the interface for kotowaza data access.
type KotowazaRepository interface {
	List(ctx context.Context, params kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*kotowaza.Kotowaza, error)
	Search(ctx context.Context, params kotowaza.SearchParams) ([]kotowaza.Kotowaza, int, error)
}
