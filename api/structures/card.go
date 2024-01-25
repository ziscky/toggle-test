package structures

import "github.com/ziscky/toggle-test/internal/models"

// Card represents the API form of a playing card
type Card struct {
	Code  string `json:"code"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
}

// CardModelToApi converts an internal card model to the API format
func CardModelToApi(card *models.Card) *Card {
	return &Card{
		Code:  card.Code,
		Value: card.Value,
		Suit:  card.Suit,
	}
}

// CardApiToModel converts an API card to the internal model of Card
func CardApiToModel(card *Card) *models.Card {
	return &models.Card{
		Code:  card.Code,
		Value: card.Value,
		Suit:  card.Suit,
	}
}

// DrawCardResponse represents the response of drawing cards from the deck
type DrawCardResponse struct {
	Cards []Card `json:"cards"`
}
