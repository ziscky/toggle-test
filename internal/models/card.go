package models

import "time"

// Card represents the database model of a playing card.
type Card struct {
	Rank      int
	Code      string
	Suit      string
	Value     string
	CreatedAt time.Time `gorm:"type:datetime"`
	UpdatedAt time.Time `gorm:"type:datetime"`
	DeletedAt time.Time `gorm:"type:datetime"`
}
