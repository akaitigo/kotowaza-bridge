package kotowaza

import (
	"time"

	"github.com/google/uuid"
)

// Kotowaza represents a Japanese proverb with its metadata.
type Kotowaza struct {
	ID           uuid.UUID    `json:"id"`
	Japanese     string       `json:"japanese"`
	Reading      string       `json:"reading"`
	Meaning      string       `json:"meaning"`
	Origin       string       `json:"origin"`
	UsageExample string       `json:"usage_example"`
	CulturalNote string       `json:"cultural_note"`
	Equivalents  []Equivalent `json:"equivalents,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
}

// Equivalent represents a foreign language equivalent of a kotowaza.
type Equivalent struct {
	ID             uuid.UUID `json:"id"`
	KotowazaID     uuid.UUID `json:"kotowaza_id"`
	Language       string    `json:"language"`
	Expression     string    `json:"expression"`
	LiteralMeaning string    `json:"literal_meaning"`
	Explanation    string    `json:"explanation"`
}

// ListParams holds parameters for listing kotowaza.
type ListParams struct {
	Limit  int
	Offset int
}

// SearchParams holds parameters for searching kotowaza.
type SearchParams struct {
	Query  string
	Limit  int
	Offset int
}
