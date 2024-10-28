package repository

import (
	"database/sql"
	"fmt"

	"github.com/shuvo-paul/sitemonitor/models"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	query := `INSERT INTO sessions (user_id, token, created_at, expires_at) 
			  VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, session.UserID, session.Token,
		session.CreatedAt, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *SessionRepository) GetByToken(token string) (*models.Session, error) {
	var session models.Session
	query := `SELECT id, user_id, token, created_at, expires_at 
			  FROM sessions WHERE token = ?`
	err := r.db.QueryRow(query, token).Scan(
		&session.ID, &session.UserID, &session.Token,
		&session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

func (r *SessionRepository) Delete(sessionID int) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := r.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

type SessionRepositoryInterface interface {
	Create(session *models.Session) error
	GetByToken(token string) (*models.Session, error)
	Delete(sessionID int) error
}

// Ensure SessionRepository implements the interface
var _ SessionRepositoryInterface = (*SessionRepository)(nil)
