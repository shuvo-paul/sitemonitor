package services

import (
	"fmt"

	"github.com/shuvo-paul/sitemonitor/internal/app/models"
	"github.com/shuvo-paul/sitemonitor/internal/app/repository"
	"github.com/shuvo-paul/sitemonitor/pkg/notification"
)

type NotificationServiceInterface interface {
	CreateSlackConfig(name string, webhookURL string, enabled bool) (*models.NotificationConfig, error)
	CreateEmailConfig(name string, recipients []string, enabled bool) (*models.NotificationConfig, error)
	UpdateConfig(config *models.NotificationConfig) (*models.NotificationConfig, error)
	GetConfig(id int) (*models.NotificationConfig, error)
	GetAllConfigs() ([]*models.NotificationConfig, error)
	DeleteConfig(id int) error
	InitializeNotifiers() error
}

type NotificationService struct {
	repo            repository.NotificationRepositoryInterface
	notificationHub *notification.NotificationHub
	smtpConfig      notification.SMTPConfig
}

func NewNotificationService(repo repository.NotificationRepositoryInterface, smtpConfig notification.SMTPConfig) *NotificationService {
	return &NotificationService{
		repo:            repo,
		notificationHub: notification.NewNotificationHub(),
		smtpConfig:      smtpConfig,
	}
}

func (s *NotificationService) CreateSlackConfig(name string, webhookURL string, enabled bool) (*models.NotificationConfig, error) {
	config := &models.NotificationConfig{
		Type:    models.NotificationTypeSlack,
		Name:    name,
		Enabled: enabled,
		Config: models.ConfigData{
			"webhook_url": webhookURL,
		},
	}

	config, err := s.repo.Create(config)
	if err != nil {
		return nil, err
	}

	if enabled {
		s.notificationHub.RegisterNotifier(notification.NewSlackNotifier(webhookURL))
	}

	return config, nil
}

func (s *NotificationService) CreateEmailConfig(name string, recipients []string, enabled bool) (*models.NotificationConfig, error) {
	config := &models.NotificationConfig{
		Type:    models.NotificationTypeEmail,
		Name:    name,
		Enabled: enabled,
		Config: models.ConfigData{
			"recipients": recipients,
		},
	}

	config, err := s.repo.Create(config)
	if err != nil {
		return nil, err
	}

	if enabled {
		s.notificationHub.RegisterNotifier(notification.NewEmailNotifier(s.smtpConfig, recipients))
	}

	return config, nil
}

func (s *NotificationService) UpdateConfig(config *models.NotificationConfig) (*models.NotificationConfig, error) {
	updatedConfig, err := s.repo.Update(config)
	if err != nil {
		return nil, err
	}

	// Reinitialize notifiers after update
	return updatedConfig, s.InitializeNotifiers()
}

func (s *NotificationService) GetConfig(id int) (*models.NotificationConfig, error) {
	return s.repo.GetByID(id)
}

func (s *NotificationService) GetAllConfigs() ([]*models.NotificationConfig, error) {
	return s.repo.GetAll()
}

func (s *NotificationService) DeleteConfig(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Reinitialize notifiers after deletion
	return s.InitializeNotifiers()
}

func (s *NotificationService) InitializeNotifiers() error {
	// Clear existing notifiers
	s.notificationHub = notification.NewNotificationHub()

	// Initialize Slack notifiers
	slackConfigs, err := s.repo.GetByType(models.NotificationTypeSlack)
	if err != nil {
		return fmt.Errorf("failed to get slack configs: %w", err)
	}

	for _, config := range slackConfigs {
		if webhookURL := config.Config.SlackConfig(); webhookURL != "" {
			s.notificationHub.RegisterNotifier(notification.NewSlackNotifier(webhookURL))
		}
	}

	// Initialize Email notifiers
	emailConfigs, err := s.repo.GetByType(models.NotificationTypeEmail)
	if err != nil {
		return fmt.Errorf("failed to get email configs: %w", err)
	}

	for _, config := range emailConfigs {
		if recipients := config.Config.EmailRecipients(); len(recipients) > 0 {
			s.notificationHub.RegisterNotifier(notification.NewEmailNotifier(s.smtpConfig, recipients))
		}
	}

	return nil
}

// GetNotificationHub returns the notification hub for use by other services
func (s *NotificationService) GetNotificationHub() *notification.NotificationHub {
	return s.notificationHub
}
