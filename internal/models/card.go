package models

import "time"

type Card struct {
	Rank      int
	Code      string
	Suit      string
	Value     string
	CreatedAt time.Time `gorm:"type:datetime"`
	UpdatedAt time.Time `gorm:"type:datetime"`
	DeletedAt time.Time `gorm:"type:datetime"`
}
