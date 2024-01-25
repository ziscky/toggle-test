package games

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ziscky/toggle-test/internal/models"
)

type SuitType string

const (
	FrenchPlayingCards SuitType = "French"
)

var (
	ErrSuiteNotSupported = errors.New("suit is not supported")
)

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

	suits := []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
	courtCards := []string{"JACK", "QUEEN", "KING"}
	pips := []string{"A"}

	for i := 2; i <= 10; i++ {
		pips = append(pips, fmt.Sprint(i))
	}

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
