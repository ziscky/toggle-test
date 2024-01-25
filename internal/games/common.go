package games

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ziscky/toggle-test/internal/models"
)

// SuitType is a string representing the playing cards suit type.
type SuitType string

const (
	// FrenchPlayingCards represents the standard 52-card deck playing suite
	FrenchPlayingCards SuitType = "French"
)

var (
	// ErrSuiteNotSupported returned when attempting to generate cards from an unsupported suit type.
	ErrSuiteNotSupported = errors.New("suit is not supported")
)

// ShuffleCards implements the Knuth shuffle algorithm to shuffle the provided cards in-place.
func ShuffleCards(cards []models.Card) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(cards); i++ {
		n := random.Intn(i + 1)
		cards[i], cards[n] = cards[n], cards[i]
	}
}

func generatePlayingCards(ctx context.Context, suitType SuitType) ([]models.Card, error) {
	switch suitType {
	case FrenchPlayingCards:
		return generateFrenchPlayingCards(), nil
	default:
		return nil, fmt.Errorf("%w :%s", ErrSuiteNotSupported, suitType)
	}
}

func generateFrenchPlayingCards() []models.Card {
	var cards []models.Card

	// ordered as per provided ranking requirements
	suits := []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
	courtCards := []string{"JACK", "QUEEN", "KING"}
	pips := []string{"A"}

	// add numbers 2...10 to the pips array
	for i := 2; i <= 10; i++ {
		pips = append(pips, fmt.Sprint(i))
	}

	// for each suit, generate the pips by concatenating the pip and the first letter of the suit e.g AS, 2S
	// the court cards are generated after the pips (due to ranking) by concateneting the first letter of the
	// court and the first letter of the suit e.g JS,QC
	rank := 1
	for _, suit := range suits {
		for _, pip := range pips {
			cards = append(cards, models.Card{
				Rank:  rank,
				Code:  fmt.Sprintf("%s%c", pip, suit[0]),
				Suit:  suit,
				Value: pip,
			})
			rank++
		}
		for _, court := range courtCards {
			cards = append(cards, models.Card{
				Rank:  rank,
				Code:  fmt.Sprintf("%c%c", court[0], suit[0]),
				Suit:  suit,
				Value: court,
			})
			rank++
		}
	}
	return cards
}
