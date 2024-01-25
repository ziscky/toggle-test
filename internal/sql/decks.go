package sql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ziscky/toggle-test/internal/models"
)

// GetDeckByID returns a populated Deck model from the database by ID.
// The Deck.Cards array is ordered by DeckCard.Position.
func (p *Persist) GetDeckByID(ctx context.Context, id uuid.UUID) (*models.Deck, error) {
	var deck models.Deck
	err := p.orm.Debug().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rows, err := tx.Model(&models.Deck{}).
			Select("decks.id, decks.shuffled, cards.rank, cards.code, cards.suit, cards.value, deck_cards.status").
			Joins("left join deck_cards on deck_cards.deck_id = decks.id").
			Joins("left join cards on deck_cards.card_code = cards.code").
			Where("decks.id = ? and deck_cards.status = ?", id, models.CardStatusOnDeck).
			Order("deck_cards.position ASC").
			Rows()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInternal, err)
		}
		defer rows.Close()

		var deckIDStr string

		for rows.Next() {
			card := models.Card{}
			var status models.CardStatus

			if err := rows.Scan(&deckIDStr, &deck.Shuffled, &card.Rank, &card.Code, &card.Suit, &card.Value, &status); err != nil {
				return fmt.Errorf("%w: %w", ErrInternal, err)
			}

			// if deck.ID is already set, skip parsing UUID again
			if deck.ID == (uuid.UUID{}) {
				deck.ID, err = uuid.Parse(deckIDStr)
				if err != nil {
					return fmt.Errorf("%w: error parsing uuid", ErrInternal)
				}
			}

			// the remaining cards on deck are identified by the status: models.CardStatusOnDeck
			if status == models.CardStatusOnDeck {
				deck.Remaining++
			}
			deck.Cards = append(deck.Cards, card)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if deck.ID == (uuid.UUID{}) {
		return nil, ErrNotFound
	}
	return &deck, nil
}

// CreateDeck adds the deck to the database and also fills in the intermediary table: deck_cards.
func (p *Persist) CreateDeck(ctx context.Context, shuffled bool, cards []models.Card) (*models.Deck, error) {
	deckID := uuid.New()
	var deckCards []models.DeckCard
	for i, card := range cards {
		deckCards = append(deckCards, models.DeckCard{
			Position:  i + 1,
			DeckID:    deckID,
			CardCode:  card.Code,
			Status:    models.CardStatusOnDeck,
			CreatedAt: time.Now(),
		})
	}

	deck := &models.Deck{
		ID:        deckID,
		Shuffled:  shuffled,
		CreatedAt: time.Now(),
		Cards:     cards,
		Remaining: len(cards),
	}
	err := p.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(deck).Error
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInternal, err)
		}

		if err := p.createDeckCards(tx, deckCards); err != nil {
			return fmt.Errorf("p.createDeckCards() failed: %w", err)
		}
		return nil
	})
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqliteForeignKeyViolationCode {
				return nil, ErrForeignKeyViolation
			}
		}
		return nil, err
	}
	return deck, nil
}
