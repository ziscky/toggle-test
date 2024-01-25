package structures

import "github.com/ziscky/toggle-test/internal/models"

type Card struct {
	Code  string `json:"code"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
}

func CardModelToApi(card *models.Card) *Card {
	return &Card{
		Code:  card.Code,
		Value: card.Value,
		Suit:  card.Suit,
	}
}

func CardApiToModel(card *Card) *models.Card {
	return &models.Card{
		Code:  card.Code,
		Value: card.Value,
		Suit:  card.Suit,
	}
}

type DrawCardRequest struct {
	Cards string
}

type DrawCardResponse struct {
	Cards []Card `json:"cards"`
}
