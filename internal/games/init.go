package games

import (
	"context"
	"fmt"

	"github.com/ziscky/toggle-test/internal/persist"
)

// InitializeGameRequirements generates game requirements when the server is started.
// Generates french playing cards and adds them to the db.
func InitializeGameRequirements(ctx context.Context, persist persist.PersistInterface) error {
	playingCards, err := generatePlayingCards(ctx, FrenchPlayingCards)
	if err != nil {
		return fmt.Errorf("generatePlayingCards() failed: %w", err)
	}
	if err := persist.CreateCards(ctx, playingCards); err != nil {
		return fmt.Errorf("persist.CreateCards() failed: %w", err)
	}
	return nil
}
