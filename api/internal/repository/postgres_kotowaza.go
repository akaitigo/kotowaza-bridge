package repository

import (
	"context"
	"fmt"

	"github.com/akaitigo/kotowaza-bridge/api/internal/domain/kotowaza"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresKotowazaRepository implements KotowazaRepository using PostgreSQL.
type PostgresKotowazaRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresKotowazaRepository creates a new PostgresKotowazaRepository.
func NewPostgresKotowazaRepository(pool *pgxpool.Pool) *PostgresKotowazaRepository {
	return &PostgresKotowazaRepository{pool: pool}
}

func (r *PostgresKotowazaRepository) List(ctx context.Context, params kotowaza.ListParams) ([]kotowaza.Kotowaza, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM kotowaza").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count kotowaza: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		"SELECT id, japanese, reading, meaning, origin, usage_example, cultural_note, created_at FROM kotowaza ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		params.Limit, params.Offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list kotowaza: %w", err)
	}
	defer rows.Close()

	result, err := scanKotowazaRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *PostgresKotowazaRepository) GetByID(ctx context.Context, id uuid.UUID) (*kotowaza.Kotowaza, error) {
	row := r.pool.QueryRow(ctx,
		"SELECT id, japanese, reading, meaning, origin, usage_example, cultural_note, created_at FROM kotowaza WHERE id = $1",
		id,
	)

	var k kotowaza.Kotowaza
	err := row.Scan(&k.ID, &k.Japanese, &k.Reading, &k.Meaning, &k.Origin, &k.UsageExample, &k.CulturalNote, &k.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get kotowaza by id: %w", err)
	}

	eqRows, err := r.pool.Query(ctx,
		"SELECT id, kotowaza_id, language, expression, literal_meaning, explanation FROM equivalent WHERE kotowaza_id = $1 ORDER BY language",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("get equivalents: %w", err)
	}
	defer eqRows.Close()

	for eqRows.Next() {
		var eq kotowaza.Equivalent
		if err := eqRows.Scan(&eq.ID, &eq.KotowazaID, &eq.Language, &eq.Expression, &eq.LiteralMeaning, &eq.Explanation); err != nil {
			return nil, fmt.Errorf("scan equivalent: %w", err)
		}
		k.Equivalents = append(k.Equivalents, eq)
	}
	if err := eqRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate equivalents: %w", err)
	}

	return &k, nil
}

func (r *PostgresKotowazaRepository) Search(ctx context.Context, params kotowaza.SearchParams) ([]kotowaza.Kotowaza, int, error) {
	query := "%" + params.Query + "%"

	var total int
	err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM kotowaza WHERE japanese LIKE $1 OR reading LIKE $1 OR meaning LIKE $1",
		query,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count search: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		"SELECT id, japanese, reading, meaning, origin, usage_example, cultural_note, created_at FROM kotowaza WHERE japanese LIKE $1 OR reading LIKE $1 OR meaning LIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		query, params.Limit, params.Offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("search kotowaza: %w", err)
	}
	defer rows.Close()

	result, err := scanKotowazaRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func scanKotowazaRows(rows pgx.Rows) ([]kotowaza.Kotowaza, error) {
	var result []kotowaza.Kotowaza
	for rows.Next() {
		var k kotowaza.Kotowaza
		if err := rows.Scan(&k.ID, &k.Japanese, &k.Reading, &k.Meaning, &k.Origin, &k.UsageExample, &k.CulturalNote, &k.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan kotowaza: %w", err)
		}
		result = append(result, k)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate kotowaza: %w", err)
	}
	return result, nil
}
