package test_data

import (
	"strings"
)

// CreateDeckRequest generates request params required to create a new deck
func CreateDeckRequest(shuffled string, cards ...string) map[string]string {
	m := map[string]string{
		"shuffled": shuffled,
	}
	if len(cards) > 0 {
		m["cards"] = strings.Join(cards, ",")
	}
	return m
}

// OpenDeckRequest generates request params required to open a deck
func OpenDeckRequest(id string) map[string]string {
	return map[string]string{
		"deck_id": id,
	}
}

// DrawCardRequest generates request params required to draw a card from a deck
func DrawCardRequest(deckID, count string) map[string]string {
	return map[string]string{
		"deck_id": deckID,
		"count":   count,
	}
}
