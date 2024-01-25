-- +goose Up
CREATE TABLE IF NOT EXISTS deck_cards (
        position INTEGER NOT NULL,
		deck_id TEXT,
        card_code TEXT,
        status TEXT CHECK(status IN ('ON-DECK','IN-HAND','PLAYED')) DEFAULT 'ON-DECK',
		created_at DATETIME,
        updated_at DATETIME,
        deleted_at DATETIME,
        FOREIGN KEY (deck_id) REFERENCES decks(id),
        FOREIGN KEY (card_code) REFERENCES cards(code)
        );

CREATE TABLE IF NOT EXISTS decks (
		id TEXT PRIMARY KEY,
        shuffled boolean,
		created_at DATETIME,
        updated_at DATETIME,
        deleted_at DATETIME);