-- +goose Up
CREATE TABLE IF NOT EXISTS cards (
		code TEXT PRIMARY KEY,
        rank integer,
        suit TEXT,
        "value" TEXT,
		created_at DATETIME,
        updated_at DATETIME,
        deleted_at DATETIME);