package structures

import (
	"github.com/google/uuid"

	"github.com/ziscky/toggle-test/internal/models"
)

// Deck represents the API form of the card deck.
type Deck struct {
	ID        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
	Cards     []Card `json:"cards"`
}

// DeckModelToApi converts an internal Deck model to the API form
func DeckModelToApi(deck *models.Deck) *Deck {
	var cards []Card
	for i := range deck.Cards {
		c := CardModelToApi(&deck.Cards[i])
		cards = append(cards, *c)
	}
	return &Deck{
		ID:        deck.ID.String(),
		Shuffled:  deck.Shuffled,
		Cards:     cards,
		Remaining: deck.Remaining,
	}
}

// DeckApiToModel converts the Deck API form to the internal model.
func DeckApiToModel(deck *Deck) (*models.Deck, error) {
	var cards []models.Card
	for i := range deck.Cards {
		c := CardApiToModel(&deck.Cards[i])
		cards = append(cards, *c)
	}

	deckID, err := uuid.Parse(deck.ID)
	if err != nil {
		return nil, err
	}

	return &models.Deck{
		ID:       deckID,
		Shuffled: deck.Shuffled,
		Cards:    cards,
	}, nil
}

// CreateDeckResponse represents the API response of creating a deck.
type CreateDeckResponse struct {
	ID        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

// OpenDeckResponse represents the API response of opening a deck.
type OpenDeckResponse struct {
	Deck
}
