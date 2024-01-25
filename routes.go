package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ziscky/toggle-test/api"
	"github.com/ziscky/toggle-test/internal/persist"
)

func initRoutes(log *logrus.Entry, persist persist.PersistInterface) *mux.Router {
	h := api.NewRequestHandler(log, persist)

	router := mux.NewRouter()
	router.HandleFunc("/deck", h.CreateDeck).Methods(http.MethodPost)
	router.HandleFunc("/deck/open", h.OpenDeck).Methods(http.MethodGet)
	router.HandleFunc("/deck/draw", h.DrawCard).Methods(http.MethodPost)

	return router
}
