-- +goose Up
CREATE TABLE notification_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    config TEXT NOT NULL DEFAULT '{}',
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_notification_configs_type ON notification_configs(type);
CREATE INDEX idx_notification_configs_enabled ON notification_configs(enabled);

-- +goose Down
DROP TABLE IF EXISTS notification_configs; 