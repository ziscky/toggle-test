package structures_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ziscky/toggle-test/api/structures"
	"github.com/ziscky/toggle-test/internal/models"
)

func TestCardModelToApi(t *testing.T) {

	tests := []struct {
		name     string
		payload  *models.Card
		expected *structures.Card
	}{
		{
			name: "success",
			payload: &models.Card{
				Rank:  1,
				Code:  "AS",
				Suit:  "SPADES",
				Value: "A",
			},
			expected: &structures.Card{
				Code:  "AS",
				Suit:  "SPADES",
				Value: "A",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := structures.CardModelToApi(tt.payload)

			assert.Equal(t, tt.expected, got)

		})
	}
}

func TestCardApiToModel(t *testing.T) {

	tests := []struct {
		name     string
		expected *models.Card
		payload  *structures.Card
	}{
		{
			name: "success",
			expected: &models.Card{
				Rank:  0,
				Code:  "AS",
				Suit:  "SPADES",
				Value: "A",
			},
			payload: &structures.Card{
				Code:  "AS",
				Suit:  "SPADES",
				Value: "A",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := structures.CardApiToModel(tt.payload)

			assert.Equal(t, tt.expected, got)

		})
	}
}
