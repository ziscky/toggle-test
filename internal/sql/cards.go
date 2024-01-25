package sql

import (
	"context"
	"errors"
	"fmt"

	"github.com/glebarez/go-sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ziscky/toggle-test/internal/models"
)

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
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqliteForeignKeyViolationCode {
				return ErrForeignKeyViolation
			}
		}
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return nil
}

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
