package sql_test

import (
	"context"
	"database/sql"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ziscky/toggle-test/internal/games"
	"github.com/ziscky/toggle-test/internal/models"
	"github.com/ziscky/toggle-test/internal/persist"
	msql "github.com/ziscky/toggle-test/internal/sql"
)

var p persist.PersistInterface
var log = logrus.New().WithField("go", runtime.Version())

func init() {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatal(err)
	}
	persist, err := msql.NewPersist(db)
	if err != nil {
		log.Fatal(err)
	}
	persist.Migrate(context.Background(), log)

	p = persist
	if err := createCards(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func TestCreateCards(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		payload []models.Card
		wantErr error
	}{
		{
			name: "success",
			payload: []models.Card{
				{
					Rank:  1,
					Code:  "AJ",
					Suit:  "SPADES",
					Value: "A",
				},
			},
			wantErr: nil,
		},
		{
			name: "error duplicate card",
			payload: []models.Card{
				{
					Rank:  1,
					Code:  "AS",
					Suit:  "SPADES",
					Value: "A",
				},
			},
			wantErr: msql.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.CreateCards(ctx, tt.payload)
			assert.ErrorIs(t, got, tt.wantErr)
		})
	}
}

func TestGetCards(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		payload  []string
		expected []models.Card
		wantErr  error
	}{
		{
			name:    "success",
			payload: []string{"AS", "2D", "JC", "QH"},
			expected: []models.Card{
				{
					Rank:  1,
					Code:  "AS",
					Suit:  "SPADES",
					Value: "A",
				},
				{
					Rank:  2,
					Code:  "2D",
					Suit:  "DIAMONDS",
					Value: "2",
				},
				{
					Rank:  11,
					Code:  "JC",
					Suit:  "CLUBS",
					Value: "JACK",
				},
				{
					Rank:  12,
					Code:  "QH",
					Suit:  "HEARTS",
					Value: "QUEEN",
				},
			},
			wantErr: nil,
		},
		{
			name:     "invalid card",
			payload:  []string{"KK"},
			expected: nil,
			wantErr:  msql.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.GetCards(ctx, tt.payload)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}

			for i := 0; i < len(tt.expected); i++ {
				assert.Equal(t, tt.expected[i].Rank, got[i].Rank)
				assert.Equal(t, tt.expected[i].Suit, got[i].Suit)
				assert.Equal(t, tt.expected[i].Code, got[i].Code)
				assert.Equal(t, tt.expected[i].Value, got[i].Value)
			}

		})
	}
}

func createCards(ctx context.Context) error {
	return games.InitializeGameRequirements(ctx, p)
}
