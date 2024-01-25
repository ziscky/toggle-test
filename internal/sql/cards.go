package sql

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ziscky/toggle-test/internal/models"
)

// CreateCards UPSERTS the provided array of cards.
func (p *Persist) CreateCards(ctx context.Context, cards []models.Card) error {
	err := p.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(&clause.OnConflict{
			Columns: []clause.Column{
				{Name: "code"},
			},
			UpdateAll: true,
		}).Create(cards).Error
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return nil
}

// GetCards returns all cards or if 'codes' are provided will filter for the provided codes.
// Results are ordered by 'rank' ASC
func (p *Persist) GetCards(ctx context.Context, codes []string) ([]models.Card, error) {
	var cards []models.Card
	q := p.orm.Model(&models.Card{})
	if len(codes) > 0 {
		q = q.Where("code in (?)", codes)
	}
	err := q.Order("rank ASC").Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// if codes has an invalid card_code, find the invalid code and add it to the error info.
	if len(codes) > 0 && len(cards) != len(codes) {
		for _, code := range codes {
			found := false
			for _, card := range cards {
				if code == card.Code {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("%w: %s", ErrNotFound, code)
			}
		}
	}

	return cards, nil
}
