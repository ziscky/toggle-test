package persist

import (
	"context"

	"github.com/google/uuid"

	"github.com/ziscky/toggle-test/internal/models"
)

// PersistInterface defines methods required to implement the required database operations.
type PersistInterface interface {
	CreateCards(ctx context.Context, cards []models.Card) error
	GetCards(ctx context.Context, codes []string) ([]models.Card, error)
	GetDeckByID(ctx context.Context, id uuid.UUID) (*models.Deck, error)
	CreateDeck(ctx context.Context, shuffled bool, cards []models.Card) (*models.Deck, error)
	UpdateDeckCardStatus(ctx context.Context, deckID uuid.UUID, cards []models.Card, status models.CardStatus) error
}
