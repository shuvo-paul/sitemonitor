package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shuvo-paul/sitemonitor/pkg/monitor"
)

var (
	ErrSiteNotFound = errors.New("site not found")
)

type SiteRepositoryInterface interface {
	Create(site *monitor.Site) (*monitor.Site, error)
	// GetByID(id int) (*monitor.Site, error)
	// GetAll() ([]*monitor.Site, error)
	// Update(site *monitor.Site) error
	// Delete(id int) error
	// UpdateStatus(id int, status string) error
}

var _ SiteRepositoryInterface = (*SiteRepository)(nil)

type SiteRepository struct {
	db *sql.DB
}

func NewSiteRepository(db *sql.DB) *SiteRepository {
	return &SiteRepository{db: db}
}

func (r *SiteRepository) Create(site *monitor.Site) (*monitor.Site, error) {
	query := `
		INSERT INTO sites (url, status, enabled, interval, status_changed_at)
		VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(
		query,
		site.URL,
		site.Status,
		site.Enabled,
		site.Interval.Seconds(),
		site.StatusChangedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	site.ID = int(id)
	return site, nil
}

func (r *SiteRepository) GetByID(id int) (*monitor.Site, error) {

	return nil, nil
}

func (r *SiteRepository) GetAll() ([]*monitor.Site, error) {

	return nil, nil
}

func (r *SiteRepository) Update(site *monitor.Site) error {

	return nil
}

func (r *SiteRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM sites WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSiteNotFound
	}

	return nil
}

func (r *SiteRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE sites
		SET status = $1, status_changed_at = NOW(), updated_at = NOW()
		WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSiteNotFound
	}

	return nil
}
