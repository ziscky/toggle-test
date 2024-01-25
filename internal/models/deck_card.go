package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	CardStatusOnDeck CardStatus = "ON-DECK"
	CardStatusInHand CardStatus = "IN-HAND"
	CardStatusPlayed CardStatus = "PLAYED"
)

// CardStatus represents the status of the card in the deck.
type CardStatus string

// DeckCard represents an intermediate table for storing the cards present in a deck along with status info.
type DeckCard struct {
	// Position is the position of the card in the deck. Results are ordered by this field ASC
	Position int
	DeckID   uuid.UUID `gorm:"type:uuid"`
	CardCode string
	// Status represents the current status of the card in the deck
	Status    CardStatus
	CreatedAt time.Time `gorm:"type:datetime"`
	UpdatedAt time.Time `gorm:"type:datetime"`
	DeletedAt time.Time `gorm:"type:datetime"`
}
