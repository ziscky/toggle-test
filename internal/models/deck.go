package models

import (
	"time"

	"github.com/google/uuid"
)

type Deck struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	Shuffled  bool
	Cards     []Card    `gorm:"-"`
	Remaining int       `gorm:"-"`
	CreatedAt time.Time `gorm:"type:datetime"`
	UpdatedAt time.Time `gorm:"type:datetime"`
	DeletedAt time.Time `gorm:"type:datetime"`
}
