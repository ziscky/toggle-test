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

type CardStatus string

type DeckCard struct {
	Position  int
	DeckID    uuid.UUID `gorm:"type:uuid"`
	CardCode  string
	Status    CardStatus
	CreatedAt time.Time `gorm:"type:datetime"`
	UpdatedAt time.Time `gorm:"type:datetime"`
	DeletedAt time.Time `gorm:"type:datetime"`
}
