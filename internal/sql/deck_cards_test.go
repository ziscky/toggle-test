package sql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ziscky/toggle-test/internal/models"
	msql "github.com/ziscky/toggle-test/internal/sql"
)

func TestUpdateDeckCardStatus(t *testing.T) {
	ctx := context.Background()

	deck, err := createDeck()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		cards   []models.Card
		status  models.CardStatus
		wantErr error
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
			},
			status:  models.CardStatusInHand,
			wantErr: nil,
		},
		{
			name: "error",
			cards: []models.Card{
				{
					Rank:  1,
					Code:  "AM",
					Suit:  "SPADES",
					Value: "A",
				},
			},
			status:  models.CardStatusInHand,
			wantErr: msql.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.UpdateDeckCardStatus(ctx, deck.ID, tt.cards, tt.status)
			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
			}
		})
	}
}
