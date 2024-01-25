package sql_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ziscky/toggle-test/internal/models"
	msql "github.com/ziscky/toggle-test/internal/sql"
)

func TestGetDeckByID(t *testing.T) {
	ctx := context.Background()

	deck, err := createDeck()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		payload  uuid.UUID
		expected *models.Deck
		wantErr  error
	}{
		{
			name:     "success",
			payload:  deck.ID,
			expected: deck,
			wantErr:  nil,
		},
		{
			name:     "not found",
			payload:  uuid.New(),
			expected: nil,
			wantErr:  msql.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.GetDeckByID(ctx, tt.payload)
			if tt.wantErr != nil {
				log.Error(err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expected.ID, got.ID)
			assert.Equal(t, tt.expected.Shuffled, got.Shuffled)
			for i := 0; i < len(tt.expected.Cards); i++ {
				assert.Equal(t, tt.expected.Cards[i].Code, got.Cards[i].Code)
				assert.Equal(t, tt.expected.Cards[i].Rank, got.Cards[i].Rank)
				assert.Equal(t, tt.expected.Cards[i].Value, got.Cards[i].Value)
				assert.Equal(t, tt.expected.Cards[i].Suit, got.Cards[i].Suit)
			}
		})
	}
}

func TestCreateDeck(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		shuffled bool
		cards    []models.Card
		expected *models.Deck
		wantErr  error
	}{
		{
			name: "success",
			cards: []models.Card{
				{
					Rank:  1,
					Code:  "AS",
					Suit:  "SPADES",
					Value: "A",
				},
				{
					Rank:  2,
					Code:  "2S",
					Suit:  "SPADES",
					Value: "2",
				},
			},
			expected: &models.Deck{
				Shuffled: true,
				Cards: []models.Card{
					{
						Rank:  1,
						Code:  "AS",
						Suit:  "SPADES",
						Value: "A",
					},
					{
						Rank:  2,
						Code:  "2S",
						Suit:  "SPADES",
						Value: "2",
					},
				},
			},
			shuffled: true,
			wantErr:  nil,
		},
		{
			name: "invslif cards",
			cards: []models.Card{
				{
					Rank:  1,
					Code:  "AM",
					Suit:  "SPADES",
					Value: "A",
				},
				{
					Rank:  2,
					Code:  "2M",
					Suit:  "SPADES",
					Value: "2",
				},
			},
			expected: nil,
			wantErr:  msql.ErrForeignKeyViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.CreateDeck(ctx, tt.shuffled, tt.cards)
			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.expected.Shuffled, got.Shuffled)
				for i := 0; i < len(tt.expected.Cards); i++ {
					assert.Equal(t, tt.expected.Cards[i].Code, got.Cards[i].Code)
					assert.Equal(t, tt.expected.Cards[i].Rank, got.Cards[i].Rank)
					assert.Equal(t, tt.expected.Cards[i].Value, got.Cards[i].Value)
					assert.Equal(t, tt.expected.Cards[i].Suit, got.Cards[i].Suit)
				}
			}
		})
	}
}

func createDeck() (*models.Deck, error) {
	return p.CreateDeck(context.Background(), true, []models.Card{
		{
			Rank:  1,
			Code:  "AS",
			Suit:  "SPADES",
			Value: "A",
		},
		{
			Rank:  2,
			Code:  "2S",
			Suit:  "SPADES",
			Value: "2",
		},
		{
			Rank:  3,
			Code:  "3S",
			Suit:  "SPADES",
			Value: "3",
		},
		{
			Rank:  4,
			Code:  "4S",
			Suit:  "SPADES",
			Value: "4",
		},
		{
			Rank:  5,
			Code:  "5S",
			Suit:  "SPADES",
			Value: "5",
		},
	})

}
