package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ziscky/toggle-test/internal/models"
)

func (p *Persist) createDeckCards(tx *gorm.DB, deckCards []models.DeckCard) error {
	return tx.Create(deckCards).Error
}

func (p *Persist) UpdateDeckCardStatus(ctx context.Context, deckID uuid.UUID, cards []models.Card, status models.CardStatus) error {
	codes := make([]string, len(cards))
	for _, card := range cards {
		codes = append(codes, card.Code)
	}

	err := p.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		q := tx.Model(&models.DeckCard{}).Where("deck_id = ? AND card_code in (?)", deckID, codes).Update("status", status)

		if err := q.Error; err != nil {
			return fmt.Errorf("%w: %w", ErrInternal, err)
		}

		if q.RowsAffected != int64(len(cards)) {
			return ErrNotFound
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}