package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeLikePattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "猿も木から落ちる",
			expected: "猿も木から落ちる",
		},
		{
			name:     "percent sign is escaped",
			input:    "100%達成",
			expected: `100\%達成`,
		},
		{
			name:     "underscore is escaped",
			input:    "a_b",
			expected: `a\_b`,
		},
		{
			name:     "backslash is escaped",
			input:    `path\to\file`,
			expected: `path\\to\\file`,
		},
		{
			name:     "multiple special characters",
			input:    `50%_off\sale`,
			expected: `50\%\_off\\sale`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    `%_\`,
			expected: `\%\_\\`,
		},
		{
			name:     "regular ascii",
			input:    "hello world",
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, escapeLikePattern(tt.input))
		})
	}
}
