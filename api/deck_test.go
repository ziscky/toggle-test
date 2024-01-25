package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ziscky/toggle-test/api"
	"github.com/ziscky/toggle-test/internal/models"
	"github.com/ziscky/toggle-test/internal/sql"
	test_data "github.com/ziscky/toggle-test/test/data"
	"github.com/ziscky/toggle-test/test/mocks"
)

var log = logrus.New().WithField("go", runtime.Version())
var uuid1Str = "a251071b-662f-44b6-ba11-e24863039c59"
var uuid1 = uuid.MustParse("a251071b-662f-44b6-ba11-e24863039c59")

func TestCreateDeck(t *testing.T) {
	p := &mocks.PersistInterface{}
	h := api.NewRequestHandler(log, p)

	tests := []struct {
		name         string
		payload      map[string]string
		mockFn       func()
		expectedCode int
		expectedBody string
	}{
		{
			name:    "Create default shuffled deck success",
			payload: test_data.CreateDeckRequest("true"),
			mockFn: func() {
				p.On("GetCards", mock.Anything, mock.Anything).Once().Return(
					[]models.Card{}, nil)
				p.On("CreateDeck", mock.Anything, mock.Anything, mock.Anything).Once().Return(
					&models.Deck{
						ID:        uuid1,
						Shuffled:  true,
						Remaining: 52,
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
				"shuffled": true,
				"remaining": 52
			}`,
		},
		{
			name:    "Create default unshuffled deck success",
			payload: test_data.CreateDeckRequest("false"),
			mockFn: func() {
				p.On("GetCards", mock.Anything, mock.Anything).Once().Return(
					[]models.Card{}, nil)
				p.On("CreateDeck", mock.Anything, mock.Anything, mock.Anything).Once().Return(
					&models.Deck{
						ID:        uuid1,
						Shuffled:  false,
						Remaining: 52,
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
				"shuffled": false,
				"remaining": 52
			}`,
		},
		{
			name:    "Create 3 card shuffled deck success",
			payload: test_data.CreateDeckRequest("true", "AS", "KS", "JS"),
			mockFn: func() {
				p.On("GetCards", mock.Anything, mock.Anything).Once().Return(
					[]models.Card{}, nil)
				p.On("CreateDeck", mock.Anything, mock.Anything, mock.Anything).Once().Return(
					&models.Deck{
						ID:        uuid1,
						Shuffled:  true,
						Remaining: 3,
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
				"shuffled": true,
				"remaining": 3
			}`,
		},
		{
			name:    "Create 3 card unshuffled deck success",
			payload: test_data.CreateDeckRequest("false", "AS", "KS", "JS"),
			mockFn: func() {
				p.On("GetCards", mock.Anything, mock.Anything).Once().Return(
					[]models.Card{}, nil)
				p.On("CreateDeck", mock.Anything, mock.Anything, mock.Anything).Once().Return(
					&models.Deck{
						ID:        uuid1,
						Shuffled:  false,
						Remaining: 3,
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
				"shuffled": false,
				"remaining": 3
			}`,
		},
		{
			name:    "Invalid shuffle parameter",
			payload: test_data.CreateDeckRequest("invalid", "AS", "KS", "JS"),
			mockFn: func() {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `invalid value for shuffle`,
		},
		{
			name:    "GetCards failed",
			payload: test_data.CreateDeckRequest("true", "AS", "KS", "JS"),
			mockFn: func() {
				p.On("GetCards", mock.Anything, mock.Anything).Once().Return(
					nil, sql.ErrInternal,
				)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `internal error`,
		},
		{
			name:    "CreateDeck failed",
			payload: test_data.CreateDeckRequest("true", "AS", "KS", "JS"),
			mockFn: func() {
				p.On("GetCards", mock.Anything, mock.Anything).Once().Return(
					[]models.Card{}, nil)
				p.On("CreateDeck", mock.Anything, mock.Anything, mock.Anything).Once().Return(
					nil, sql.ErrInternal)

			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `internal error`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			gotCode, gotBody := doRequest(t, h.CreateDeck, http.MethodPost, "/deck", tt.payload, nil)
			assert.Equal(t, tt.expectedCode, gotCode)
			assertBodyEqual(t, tt.expectedBody, gotBody)
		})
	}
}

func TestOpenDeck(t *testing.T) {
	log.Logger.SetReportCaller(true)
	p := &mocks.PersistInterface{}
	h := api.NewRequestHandler(log, p)

	tests := []struct {
		name         string
		payload      map[string]string
		mockFn       func()
		expectedCode int
		expectedBody string
	}{
		{
			name:    "Open deck success",
			payload: test_data.OpenDeckRequest(uuid1Str),
			mockFn: func() {
				p.On("GetDeckByID", mock.Anything, mock.Anything).Once().Return(
					&models.Deck{
						ID:        uuid1,
						Remaining: 3,
						Cards: []models.Card{
							{
								Value: "ACE",
								Suit:  "SPADES",
								Code:  "AS",
							},
							{
								Value: "KING",
								Suit:  "HEARTS",
								Code:  "KH",
							},
							{
								Value: "8",
								Suit:  "CLUBS",
								Code:  "8C",
							},
						},
					}, nil,
				)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
				"shuffled": false,
				"remaining": 3,
				"cards": [
					{
						"code": "AS",
						"value": "ACE",
						"suit": "SPADES"
					},
					{
						"code": "KH",
						"value": "KING",
						"suit": "HEARTS"
					},
					{
						"code": "8C",
						"value": "8",
						"suit": "CLUBS"
					}
				]
			}`,
		},
		{
			name:    "invalid deck uuid",
			payload: test_data.OpenDeckRequest("invalid"),
			mockFn: func() {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `invalid deck_id`,
		},
		{
			name:    "Get deck internal error",
			payload: test_data.OpenDeckRequest(uuid1Str),
			mockFn: func() {
				p.On("GetDeckByID", mock.Anything, mock.Anything).Once().Return(
					nil, sql.ErrInternal)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `internal error`,
		},
		{
			name:    "Get deck id does not exist",
			payload: test_data.OpenDeckRequest(uuid1Str),
			mockFn: func() {
				p.On("GetDeckByID", mock.Anything, mock.Anything).Once().Return(
					nil, sql.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `resource not found`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			gotCode, gotBody := doRequest(t, h.OpenDeck, http.MethodGet, "deck/open", tt.payload, nil)
			assert.Equal(t, tt.expectedCode, gotCode)
			assertBodyEqual(t, tt.expectedBody, gotBody)
		})
	}
}

func TestDrawCard(t *testing.T) {
	log.Logger.SetReportCaller(true)
	p := &mocks.PersistInterface{}
	h := api.NewRequestHandler(log, p)

	tests := []struct {
		name         string
		payload      map[string]string
		mockFn       func()
		expectedCode int
		expectedBody string
	}{
		{
			name:    "Draw card success",
			payload: test_data.DrawCardRequest(uuid1Str, "3"),
			mockFn: func() {
				p.On("GetDeckByID", mock.Anything, mock.Anything).Once().Return(
					&models.Deck{
						ID: uuid1,
						Cards: []models.Card{
							{
								Value: "ACE",
								Suit:  "SPADES",
								Code:  "AS",
							},
							{
								Value: "KING",
								Suit:  "HEARTS",
								Code:  "KH",
							},
							{
								Value: "8",
								Suit:  "CLUBS",
								Code:  "8C",
							},
						},
					}, nil)
				p.On("UpdateDeckCardStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(
					nil, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
				"cards": [
					{
						"code": "AS",
						"value": "ACE",
						"suit": "SPADES"
					},
					{
						"code": "KH",
						"value": "KING",
						"suit": "HEARTS"
					},
					{
						"code": "8C",
						"value": "8",
						"suit": "CLUBS"
					}
				]
			}`,
		},
		{
			name:    "invalid deck_id",
			payload: test_data.DrawCardRequest("invalid", "2"),
			mockFn: func() {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `invalid deck_id`,
		},
		{
			name:    "invalid count",
			payload: test_data.DrawCardRequest(uuid1Str, "three"),
			mockFn: func() {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `invalid count`,
		},
		{
			name:    "Get deck internal error",
			payload: test_data.DrawCardRequest(uuid1Str, "3"),
			mockFn: func() {
				p.On("GetDeckByID", mock.Anything, mock.Anything).Once().Return(
					nil, sql.ErrInternal)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `internal error`,
		},
		{
			name:    "Get deck id does not exist",
			payload: test_data.DrawCardRequest(uuid1Str, "3"),
			mockFn: func() {
				p.On("GetDeckByID", mock.Anything, mock.Anything).Once().Return(
					nil, sql.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `resource not found`,
		},
		{
			name:    "Update deck card status failed",
			payload: test_data.DrawCardRequest(uuid1Str, "3"),
			mockFn: func() {
				p.On("GetDeckByID", mock.Anything, mock.Anything).Once().Return(
					&models.Deck{
						ID: uuid1,
						Cards: []models.Card{
							{
								Value: "ACE",
								Suit:  "SPADES",
								Code:  "AS",
							},
							{
								Value: "KING",
								Suit:  "HEARTS",
								Code:  "KH",
							},
							{
								Value: "8",
								Suit:  "CLUBS",
								Code:  "8C",
							},
						},
					}, nil)
				p.On("UpdateDeckCardStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(
					sql.ErrInternal)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `resource not found`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			gotCode, gotBody := doRequest(t, h.DrawCard, http.MethodPost, "/deck/draw", tt.payload, nil)
			assert.Equal(t, tt.expectedCode, gotCode)
			if tt.expectedCode == http.StatusOK {
				expected := unmarshalBody(t, []byte(tt.expectedBody))
				got := unmarshalBody(t, []byte(gotBody))
				assert.ElementsMatch(t, expected["cards"], got["cards"])
			}
		})
	}
}

func doRequest(t *testing.T, handler http.HandlerFunc, method, endpoint string, queryParams map[string]string, body any) (int, string) {
	var payload io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		payload = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	if len(queryParams) > 0 {
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr.Code, rr.Body.String()
}

func assertBodyEqual(t *testing.T, rq, resp string) {
	r := strings.NewReplacer("\n", "", "\t", "", " ", "")

	assert.Equal(t, r.Replace(rq), r.Replace(resp))
}

func unmarshalBody(t *testing.T, data []byte) map[string]any {
	p := map[string]any{}
	assert.NoError(t, json.Unmarshal(data, &p))
	return p
}
