package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ziscky/toggle-test/internal/models"
)

// createDeckCards adds the array of deckCards to the database.
func (p *Persist) createDeckCards(tx *gorm.DB, deckCards []models.DeckCard) error {
	// inherints context from the provided transaction
	return tx.Create(deckCards).Error
}

// UpdateDeckCardStatus will update the cards provided in the deckID provided to the status provided.
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
