package games

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ziscky/toggle-test/internal/models"
)

func TestShuffleCards(t *testing.T) {
	cards := []models.Card{
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
	}
	var shuffled []models.Card
	for i := 0; i < 5; i++ {

		ShuffleCards(cards)
		assert.NotEqual(t, cards, shuffled)
		copy(shuffled, cards)
	}

}
func Test_generateFrenchPlayingCards(t *testing.T) {
	valid := validCardOrdered()

	tests := []struct {
		name     string
		payload  SuitType
		expected []string
		wantErr  bool
	}{
		{
			name:     "success",
			payload:  FrenchPlayingCards,
			expected: valid,
			wantErr:  false,
		},
		{
			name:     "unimplemented suite type",
			payload:  "British",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generatePlayingCards(context.Background(), tt.payload)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}

			var codes []string
			for _, card := range got {
				codes = append(codes, card.Code)
			}
			assert.Equal(t, tt.expected, codes)
		})
	}
}

func validCardOrdered() []string {
	return []string{"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "JS", "QS", "KS", "AD", "2D",
		"3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "JD", "QD", "KD", "AC", "2C", "3C", "4C", "5C", "6C", "7C",
		"8C", "9C", "10C", "JC", "QC", "KC", "AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "JH", "QH", "KH"}
}
