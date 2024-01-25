package structures_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ziscky/toggle-test/api/structures"
	"github.com/ziscky/toggle-test/internal/models"
)

var (
	uuid1Str = "a251071b-662f-44b6-ba11-e24863039c59"
	uuid1    = uuid.MustParse("a251071b-662f-44b6-ba11-e24863039c59")
)

func TestDeckModelToApi(t *testing.T) {

	tests := []struct {
		name     string
		payload  *models.Deck
		expected *structures.Deck
	}{
		{
			name: "success",
			payload: &models.Deck{
				ID:       uuid1,
				Shuffled: true,
				Cards: []models.Card{
					{
						Rank:  1,
						Code:  "AS",
						Suit:  "SPADES",
						Value: "A",
					},
				},
			},
			expected: &structures.Deck{
				ID:       uuid1Str,
				Shuffled: true,
				Cards: []structures.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "A",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := structures.DeckModelToApi(tt.payload)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestDeckApiToModel(t *testing.T) {
	tests := []struct {
		name     string
		expected *models.Deck
		payload  *structures.Deck
		wantErr  bool
	}{
		{
			name: "success",
			expected: &models.Deck{
				ID:       uuid1,
				Shuffled: true,
				Cards: []models.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "A",
					},
				},
			},
			payload: &structures.Deck{
				ID:        uuid1Str,
				Shuffled:  true,
				Remaining: 1,
				Cards: []structures.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "A",
					},
				},
			},
		},
		{
			name:     "invalid deck id",
			expected: nil,
			payload: &structures.Deck{
				ID:        "invalid",
				Shuffled:  true,
				Remaining: 1,
				Cards: []structures.Card{
					{
						Code:  "AS",
						Suit:  "SPADES",
						Value: "A",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := structures.DeckApiToModel(tt.payload)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
