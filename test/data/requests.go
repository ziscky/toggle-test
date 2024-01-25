package test_data

import (
	"strings"
)

func CreateDeckRequest(shuffled string, cards ...string) map[string]string {
	m := map[string]string{
		"shuffled": shuffled,
	}
	if len(cards) > 0 {
		m["cards"] = strings.Join(cards, ",")
	}
	return m
}

func OpenDeckRequest(id string) map[string]string {
	return map[string]string{
		"deck_id": id,
	}
}

func DrawCardRequest(deckID, count string) map[string]string {
	return map[string]string{
		"deck_id": deckID,
		"count":   count,
	}
}
