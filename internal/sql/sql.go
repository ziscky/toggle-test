package sql

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/glebarez/sqlite"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var (
	ErrNotFound                   = errors.New("resource not found")
	ErrInternal                   = errors.New("internal error")
	ErrDuplicateKey               = errors.New("duplicate key")
	ErrForeignKeyViolation        = errors.New("foreign key violation")
	sqliteForeignKeyViolationCode = 787
)

type Persist struct {
	db  *sql.DB
	orm *gorm.DB
}

func NewPersist(db *sql.DB) (*Persist, error) {
	gormDB, err := gorm.Open(sqlite.Dialector{Conn: db})
	if err != nil {
		return nil, err
	}

	return &Persist{
		db:  db,
		orm: gormDB,
	}, nil
}

func (p *Persist) Migrate(ctx context.Context, log *logrus.Entry) error {
	goose.SetLogger(log)
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(p.db, "migrations"); err != nil {
		return err
	}
	return nil
}
