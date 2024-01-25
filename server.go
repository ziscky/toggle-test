package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/sirupsen/logrus"

	"github.com/ziscky/toggle-test/internal/games"
	msql "github.com/ziscky/toggle-test/internal/sql"
)

const defaultHTTPTimeout = 10 * time.Second

type Server struct {
	db  *sql.DB
	srv *http.Server
	log *logrus.Entry
}

func NewServer(log *logrus.Entry, options *options) (*Server, error) {
	db, err := connectDB(options.dbFilePath)
	if err != nil {
		return nil, err
	}

	persist, err := msql.NewPersist(db)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := persist.Migrate(ctx, log); err != nil {
		return nil, err
	}

	router := initRoutes(log, persist)

	if err := games.InitializeGameRequirements(ctx, persist); err != nil {
		return nil, err
	}

	return &Server{
		db: db,
		srv: &http.Server{
			Addr:         options.addr,
			Handler:      router,
			ReadTimeout:  defaultHTTPTimeout,
			WriteTimeout: defaultHTTPTimeout,
		},
		log: log,
	}, nil
}

func (s *Server) Start() error {
	s.log.Infof("starting server on %s", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return s.db.Close()
}

func connectDB(dbFilePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", fmt.Sprintf("%s?_pragma=foreign_keys(1)", dbFilePath))
	if err != nil {
		return nil, err
	}

	return db, nil
}
