-- +migrate Up
CREATE TABLE notifiers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS notifiers;
