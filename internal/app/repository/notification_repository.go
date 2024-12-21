package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/shuvo-paul/sitemonitor/internal/app/models"
)

type NotificationRepositoryInterface interface {
	Create(config *models.NotificationConfig) (*models.NotificationConfig, error)
	Update(config *models.NotificationConfig) (*models.NotificationConfig, error)
	GetByID(id int) (*models.NotificationConfig, error)
	GetByType(notificationType models.NotificationType) ([]*models.NotificationConfig, error)
	GetAll() ([]*models.NotificationConfig, error)
	Delete(id int) error
}

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// boolToInt converts a boolean to SQLite integer (0 or 1)
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// intToBool converts SQLite integer (0 or 1) to boolean
func intToBool(i int) bool {
	return i != 0
}

func (r *NotificationRepository) Create(config *models.NotificationConfig) (*models.NotificationConfig, error) {
	query := `
		INSERT INTO notification_configs (type, name, config, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now().UTC().Format(time.RFC3339)
	config.CreatedAt = now
	config.UpdatedAt = now

	configJSON, err := config.Config.Value()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	result, err := r.db.Exec(
		query,
		config.Type,
		config.Name,
		configJSON,
		boolToInt(config.Enabled),
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification config: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	config.ID = int(id)
	return config, nil
}

func (r *NotificationRepository) Update(config *models.NotificationConfig) (*models.NotificationConfig, error) {
	query := `
		UPDATE notification_configs
		SET type = ?, name = ?, config = ?, enabled = ?, updated_at = ?
		WHERE id = ?`

	now := time.Now().UTC().Format(time.RFC3339)
	config.UpdatedAt = now

	configJSON, err := config.Config.Value()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	result, err := r.db.Exec(
		query,
		config.Type,
		config.Name,
		configJSON,
		boolToInt(config.Enabled),
		now,
		config.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update notification config: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	return config, nil
}

func (r *NotificationRepository) GetByID(id int) (*models.NotificationConfig, error) {
	query := `
		SELECT id, type, name, config, enabled, created_at, updated_at
		FROM notification_configs
		WHERE id = ?`

	var (
		config    models.NotificationConfig
		enabled   int
		configStr string
	)

	err := r.db.QueryRow(query, id).Scan(
		&config.ID,
		&config.Type,
		&config.Name,
		&configStr,
		&enabled,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get notification config: %w", err)
	}

	config.Enabled = intToBool(enabled)
	if err := config.Config.Scan([]byte(configStr)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (r *NotificationRepository) GetByType(notificationType models.NotificationType) ([]*models.NotificationConfig, error) {
	query := `
		SELECT id, type, name, config, enabled, created_at, updated_at
		FROM notification_configs
		WHERE type = ? AND enabled = 1`

	rows, err := r.db.Query(query, notificationType)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.NotificationConfig
	for rows.Next() {
		var (
			config    models.NotificationConfig
			enabled   int
			configStr string
		)

		err := rows.Scan(
			&config.ID,
			&config.Type,
			&config.Name,
			&configStr,
			&enabled,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification config: %w", err)
		}

		config.Enabled = intToBool(enabled)
		if err := config.Config.Scan([]byte(configStr)); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

func (r *NotificationRepository) GetAll() ([]*models.NotificationConfig, error) {
	query := `
		SELECT id, type, name, config, enabled, created_at, updated_at
		FROM notification_configs`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.NotificationConfig
	for rows.Next() {
		var (
			config    models.NotificationConfig
			enabled   int
			configStr string
		)

		err := rows.Scan(
			&config.ID,
			&config.Type,
			&config.Name,
			&configStr,
			&enabled,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification config: %w", err)
		}

		config.Enabled = intToBool(enabled)
		if err := config.Config.Scan([]byte(configStr)); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

func (r *NotificationRepository) Delete(id int) error {
	query := `DELETE FROM notification_configs WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification config: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
