-- +migrate Up
CREATE TABLE sites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL,
    status VARCHAR(50),
    enabled BOOLEAN DEFAULT true,
    interval FLOAT NOT NULL,
    status_changed_at TIMESTAMP WITH TIME ZONE
);

-- +migrate Down
DROP TABLE sites;