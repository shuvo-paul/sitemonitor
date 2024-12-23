package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Notifier represents a row in the notifiers table
type Notifier struct {
	ID        int64          `json:"id"`
	Config    NotifierConfig `json:"config"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// NotifierRepository handles database operations for notifiers
type NotifierRepository struct {
	db *sql.DB
}

// NewNotifierRepository creates a new notifier repository
func NewNotifierRepository(db *sql.DB) *NotifierRepository {
	return &NotifierRepository{db: db}
}

// Create inserts a new notifier into the database
func (r *NotifierRepository) Create(config NotifierConfig) (*Notifier, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	now := time.Now().UTC()
	query := `
		INSERT INTO notifiers (config, created_at, updated_at)
		VALUES (?, ?, ?)
		RETURNING id, config, created_at, updated_at
	`

	var notifier Notifier
	var configStr string
	err = r.db.QueryRow(query, string(configJSON), now, now).Scan(
		&notifier.ID,
		&configStr,
		&notifier.CreatedAt,
		&notifier.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create notifier: %w", err)
	}

	err = json.Unmarshal([]byte(configStr), &notifier.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &notifier, nil
}

// Get retrieves a notifier by ID
func (r *NotifierRepository) Get(id int64) (*Notifier, error) {
	query := `
		SELECT id, config, created_at, updated_at
		FROM notifiers
		WHERE id = ?
	`

	var notifier Notifier
	var configStr string
	err := r.db.QueryRow(query, id).Scan(
		&notifier.ID,
		&configStr,
		&notifier.CreatedAt,
		&notifier.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notifier: %w", err)
	}

	err = json.Unmarshal([]byte(configStr), &notifier.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &notifier, nil
}

// List retrieves all notifiers
func (r *NotifierRepository) List() ([]Notifier, error) {
	query := `
		SELECT id, config, created_at, updated_at
		FROM notifiers
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifiers: %w", err)
	}
	defer rows.Close()

	var notifiers []Notifier
	for rows.Next() {
		var notifier Notifier
		var configStr string
		err := rows.Scan(
			&notifier.ID,
			&configStr,
			&notifier.CreatedAt,
			&notifier.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notifier: %w", err)
		}

		err = json.Unmarshal([]byte(configStr), &notifier.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		notifiers = append(notifiers, notifier)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifiers: %w", err)
	}

	return notifiers, nil
}

// Update updates a notifier's configuration
func (r *NotifierRepository) Update(id int64, config NotifierConfig) (*Notifier, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	now := time.Now().UTC()
	query := `
		UPDATE notifiers
		SET config = ?, updated_at = ?
		WHERE id = ?
		RETURNING id, config, created_at, updated_at
	`

	var notifier Notifier
	var configStr string
	err = r.db.QueryRow(query, string(configJSON), now, id).Scan(
		&notifier.ID,
		&configStr,
		&notifier.CreatedAt,
		&notifier.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update notifier: %w", err)
	}

	err = json.Unmarshal([]byte(configStr), &notifier.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &notifier, nil
}

// Delete removes a notifier from the database
func (r *NotifierRepository) Delete(id int64) error {
	query := `DELETE FROM notifiers WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notifier: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return nil
	}

	return nil
}
