package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shuvo-paul/sitemonitor/models"
	"github.com/shuvo-paul/sitemonitor/repository"
	"golang.org/x/crypto/bcrypt"
)

type SessionServiceInterface interface {
	CreateSession(userID int) (*models.Session, string, error)
	ValidateSession(token string) (*models.Session, error)
	DeleteSession(sessionID int) error
}

var _ SessionServiceInterface = (*SessionService)(nil)

type SessionService struct {
	sessionRepo repository.SessionRepositoryInterface
}

func NewSessionService(sessionRepo repository.SessionRepositoryInterface) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (s *SessionService) CreateSession(userID int) (*models.Session, string, error) {
	// Generate a unique token
	plainToken := uuid.New().String()

	// Hash the token
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(plainToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash token: %w", err)
	}

	session := &models.Session{
		UserID:    userID,
		Token:     string(hashedToken),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, "", err
	}

	return session, plainToken, nil
}

func (s *SessionService) ValidateSession(token string) (*models.Session, error) {
	session, err := s.sessionRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		s.sessionRepo.Delete(session.ID)
		return nil, fmt.Errorf("session has expired")
	}

	return session, nil
}

func (s *SessionService) DeleteSession(sessionID int) error {
	return s.sessionRepo.Delete(sessionID)
}
