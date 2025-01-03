package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/shuvo-paul/sitemonitor/internal/app/models"
)

type NotifierRepositoryInterface interface {
	Create(*models.Notifier) error
	Get(int64) (*models.Notifier, error)
	Update(int, *models.NotifierConfig) (*models.Notifier, error)
	Delete(int64) error
	GetBySiteID(int) ([]*models.Notifier, error)
}

var _ NotifierRepositoryInterface = (*NotifierRepository)(nil)

// NotifierRepository handles database operations for notifiers
type NotifierRepository struct {
	db *sql.DB
}

// NewNotifierRepository creates a new notifier repository
func NewNotifierRepository(db *sql.DB) *NotifierRepository {
	return &NotifierRepository{db: db}
}

// Create inserts a new notifier into the database
func (r *NotifierRepository) Create(notifier *models.Notifier) error {
	query := `
		INSERT INTO notifiers (site_id, config)
		VALUES (?, ?)
		RETURNING id
	`

	configBytes, err := json.Marshal(notifier.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = r.db.QueryRow(query, notifier.SiteId, configBytes).Scan(&notifier.ID)
	if err != nil {
		return fmt.Errorf("failed to create notifier: %w", err)
	}

	return nil
}

// Get retrieves a notifier by ID
func (r *NotifierRepository) Get(id int64) (*models.Notifier, error) {
	query := `
		SELECT *
		FROM notifiers
		WHERE id = ?
	`

	notifier := &models.Notifier{}

	err := r.db.QueryRow(query, id).Scan(
		&notifier.ID,
		&notifier.SiteId,
		&notifier.Config,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notifier: %w", err)
	}

	return notifier, nil
}

// Update updates a notifier's configuration
func (r *NotifierRepository) Update(id int, config *models.NotifierConfig) (*models.Notifier, error) {
	query := `
		UPDATE notifiers
		SET config = ?
		WHERE id = ?
		RETURNING *
	`

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	notifier := &models.Notifier{}

	err = r.db.QueryRow(query, configBytes, id).Scan(&notifier.ID, &notifier.SiteId, &notifier.Config)

	if err != nil {
		return nil, fmt.Errorf("failed to update: %w", err)
	}

	return notifier, nil
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

// GetBySiteID retrieves all notifiers for a specific site
func (r *NotifierRepository) GetBySiteID(siteID int) ([]*models.Notifier, error) {
	query := `
		SELECT id, site_id, config
		FROM notifiers
		WHERE site_id = ?
	`

	rows, err := r.db.Query(query, siteID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifiers: %w", err)
	}
	defer rows.Close()

	var notifiers []*models.Notifier
	for rows.Next() {
		notifier := &models.Notifier{}
		err := rows.Scan(&notifier.ID, &notifier.SiteId, &notifier.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notifier: %w", err)
		}
		notifiers = append(notifiers, notifier)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifiers: %w", err)
	}

	return notifiers, nil
}
