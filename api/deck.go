package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ziscky/toggle-test/api/structures"
	"github.com/ziscky/toggle-test/internal/games"
	"github.com/ziscky/toggle-test/internal/models"
	"github.com/ziscky/toggle-test/internal/persist"
	msql "github.com/ziscky/toggle-test/internal/sql"
)

type RequestHandler struct {
	persist persist.PersistInterface
	log     *logrus.Entry
}

func NewRequestHandler(log *logrus.Entry, persist persist.PersistInterface) *RequestHandler {
	return &RequestHandler{
		persist: persist,
		log:     log,
	}
}

func (h *RequestHandler) CreateDeck(rw http.ResponseWriter, r *http.Request) {
	var (
		err         error
		shuffle     bool
		customCards []string
	)

	ctx := r.Context()

	cardsString := r.URL.Query().Get("cards")
	if cardsString != "" {
		customCards = strings.Split(cardsString, ",")
	}

	shuffleString := r.URL.Query().Get("shuffled")
	if shuffleString != "" {
		shuffle, err = strconv.ParseBool(shuffleString)
		if err != nil {
			h.log.Errorf("strconv.ParseBool() failed: %s", err)
			writeResponseString(rw, http.StatusBadRequest, "invalid value for shuffle")
			return
		}
	}

	cards, err := h.persist.GetCards(ctx, customCards)
	if err != nil {
		h.log.Errorf("persist.GetCards() failed: %s", err)
		writeResponseString(rw, http.StatusInternalServerError, err.Error())
		return
	}

	if shuffle {
		games.ShuffleCards(cards)
	}

	deck, err := h.persist.CreateDeck(ctx, shuffle, cards)
	if err != nil {
		h.log.Errorf("persist.CreateDeck() failed: %s", err)
		writeResponseString(rw, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &structures.CreateDeckResponse{
		ID:        deck.ID.String(),
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
	}
	payload, err := marshalResponse(resp)
	if err != nil {
		h.log.Errorf("marshalResponse() failed: %s", err)
		writeResponseString(rw, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponseBytes(rw, http.StatusOK, payload)
}

func (h *RequestHandler) OpenDeck(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	deckIDStr := r.URL.Query().Get("deck_id")

	deckID, err := uuid.Parse(deckIDStr)
	if err != nil {
		writeResponseString(rw, http.StatusBadRequest, "invalid deck_id")
		return
	}

	deck, err := h.persist.GetDeckByID(ctx, deckID)
	if err != nil {
		if errors.Is(err, msql.ErrNotFound) {
			writeResponseString(rw, http.StatusNotFound, err.Error())
			return
		}
		h.log.Errorf("persist.GetDeckByID() failed: %s", err)
		writeResponseString(rw, http.StatusInternalServerError, err.Error())
		return
	}

	deckApi := structures.DeckModelToApi(deck)
	resp := &structures.OpenDeckResponse{
		Deck: *deckApi,
	}

	payload, err := marshalResponse(resp)
	if err != nil {
		h.log.Errorf("marshalResponse() failed: %s", err)
		writeResponseString(rw, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponseBytes(rw, http.StatusOK, payload)
}

func (h *RequestHandler) DrawCard(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	countStr := r.URL.Query().Get("count")

	deckIDStr := r.URL.Query().Get("deck_id")
	deckID, err := uuid.Parse(deckIDStr)
	if err != nil {
		writeResponseString(rw, http.StatusBadRequest, "invalid deck_id")
		return
	}

	var count = 1
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			writeResponseString(rw, http.StatusBadRequest, "invalid count")
			return
		}
	}

	deck, err := h.persist.GetDeckByID(ctx, deckID)
	if err != nil {
		if errors.Is(err, msql.ErrNotFound) {
			writeResponseString(rw, http.StatusNotFound, err.Error())
			return
		}
		h.log.Errorf("persist.GetDeckByID() failed: %s", err)
		writeResponseString(rw, http.StatusInternalServerError, err.Error())
		return
	}

	drawn := deck.Cards[:count]

	if err := h.persist.UpdateDeckCardStatus(ctx, deckID, drawn, models.CardStatusInHand); err != nil {
		h.log.Errorf("persist.UpdateDeckCardStatus() failed: %s", err)
		writeResponseString(rw, http.StatusInternalServerError, err.Error())
		return
	}

	drawnApi := make([]structures.Card, len(drawn))
	for i, tmp := range drawn {
		drawnApi[i] = *structures.CardModelToApi(&tmp)
	}

	payload, err := marshalResponse(&structures.DrawCardResponse{
		Cards: drawnApi,
	})
	if err != nil {
		h.log.Errorf("marshalResponse() failed: %s", err)
		writeResponseString(rw, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponseBytes(rw, http.StatusOK, payload)
}
