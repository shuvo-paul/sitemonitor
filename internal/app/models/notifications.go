package models

import (
	"database/sql"
	"encoding/json"
)

type NotificationType string

const (
	NotificationTypeSlack NotificationType = "slack"
	NotificationTypeEmail NotificationType = "email"
)

// NotificationConfig stores notification settings
type NotificationConfig struct {
	ID        int              `db:"id"`
	Type      NotificationType `db:"type"`
	Name      string           `db:"name"`
	Config    ConfigData       `db:"config"`
	Enabled   bool             `db:"enabled"`
	CreatedAt string           `db:"created_at"`
	UpdatedAt string           `db:"updated_at"`
}

// ConfigData is a flexible JSON structure for different notification configs
type ConfigData map[string]interface{}

// Scan implements sql.Scanner interface
func (c *ConfigData) Scan(value interface{}) error {
	if value == nil {
		*c = make(ConfigData)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return sql.ErrNoRows
	}

	return json.Unmarshal(bytes, c)
}

// Value implements driver.Valuer interface
func (c ConfigData) Value() (interface{}, error) {
	return json.Marshal(c)
}

// SlackConfig returns config data for Slack
func (c ConfigData) SlackConfig() string {
	if webhook, ok := c["webhook_url"].(string); ok {
		return webhook
	}
	return ""
}

// EmailRecipients returns config data for email
func (c ConfigData) EmailRecipients() []string {
	if recipients, ok := c["recipients"].([]interface{}); ok {
		emails := make([]string, 0, len(recipients))
		for _, r := range recipients {
			if email, ok := r.(string); ok {
				emails = append(emails, email)
			}
		}
		return emails
	}
	return []string{}
}
