package service

import (
	"context"
	"fmt"

	"github.com/akaitigo/kotowaza-bridge/api/internal/domain/kotowaza"
	"github.com/akaitigo/kotowaza-bridge/api/internal/repository"
	"github.com/google/uuid"
)

// KotowazaService provides business logic for kotowaza operations.
type KotowazaService struct {
	repo repository.KotowazaRepository
}

// NewKotowazaService creates a new KotowazaService.
func NewKotowazaService(repo repository.KotowazaRepository) *KotowazaService {
	return &KotowazaService{repo: repo}
}

// ListResponse holds paginated list results.
type ListResponse struct {
	Items  []kotowaza.Kotowaza `json:"items"`
	Total  int                 `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

// List returns a paginated list of kotowaza.
func (s *KotowazaService) List(ctx context.Context, limit, offset int) (*ListResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := s.repo.List(ctx, kotowaza.ListParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, fmt.Errorf("list kotowaza: %w", err)
	}

	return &ListResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// GetByID returns a single kotowaza with its equivalents.
func (s *KotowazaService) GetByID(ctx context.Context, id uuid.UUID) (*kotowaza.Kotowaza, error) {
	k, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get kotowaza: %w", err)
	}
	return k, nil
}

// Search finds kotowaza matching the query.
func (s *KotowazaService) Search(ctx context.Context, query string, limit, offset int) (*ListResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := s.repo.Search(ctx, kotowaza.SearchParams{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("search kotowaza: %w", err)
	}

	return &ListResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
