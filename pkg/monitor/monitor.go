package monitor

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const (
	statusUp     = "up"
	statusError  = "error"
	statusDown   = "down"
	statusPaused = "paused"
)

// ClientConfig holds HTTP client configuration
type ClientConfig struct {
	Timeout         time.Duration
	MaxIdleConns    int
	IdleConnTimeout time.Duration
}

// DefaultClientConfig provides sensible defaults
var DefaultClientConfig = ClientConfig{
	Timeout:         10 * time.Second,
	MaxIdleConns:    100,
	IdleConnTimeout: 90 * time.Second,
}

type Site struct {
	ID              int
	URL             string
	Status          string
	Enabled         bool
	Interval        time.Duration
	StatusChangedAt time.Time
	mu              sync.RWMutex
	cancelFunc      context.CancelFunc
	client          *http.Client // Add dedicated client per site
}

// NewSite creates a new Site with configured HTTP client
func NewSite(id int, url string, interval time.Duration, config ClientConfig) *Site {
	transport := &http.Transport{
		MaxIdleConns:    config.MaxIdleConns,
		IdleConnTimeout: config.IdleConnTimeout,
	}

	return &Site{
		ID:       id,
		URL:      url,
		Interval: interval,
		Enabled:  true,
		client: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
	}
}

func (s *Site) Check() error {
	r, err := s.client.Get(s.URL) // Use site-specific client

	if err != nil {
		s.updateStatus(statusError)
		return fmt.Errorf("connection error: %w", err)
	}

	defer r.Body.Close()

	if r.StatusCode >= 400 {
		s.updateStatus(statusDown)
		return fmt.Errorf("HTTP error: %d", r.StatusCode)
	}

	s.updateStatus(statusUp)

	return nil
}

func (s *Site) updateStatus(status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Status != status {
		s.Status = status
		s.StatusChangedAt = time.Now()
	}
}

type Manager struct {
	mu    sync.Mutex
	sites map[int]*Site
}

func NewManager() *Manager {
	return &Manager{
		sites: make(map[int]*Site),
	}
}

func (m *Manager) RegisterSite(site *Site) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sites[site.ID]; ok {
		return fmt.Errorf("site %s already being monitored", site.URL)
	}

	ctx, cancel := context.WithCancel(context.Background())
	site.cancelFunc = cancel

	m.sites[site.ID] = site

	go func() {
		ticker := time.NewTicker(site.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				slog.Info("Monitoring stopperd", "site", site.URL)
				m.mu.Lock()
				delete(m.sites, site.ID)
				m.mu.Unlock()
				return
			case <-ticker.C:
				if !site.Enabled {
					continue
				}
				if err := site.Check(); err != nil {
					slog.Error("Site check failed", "site", site.URL, "error", err)
				}
			}
		}

	}()

	slog.Info("Monitoring started", "site", site.URL)
	return nil
}

func (m *Manager) RevokeSite(siteID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if site, exist := m.sites[siteID]; exist {
		site.cancelFunc()
		delete(m.sites, siteID)
		slog.Info("Monitoring Stopped", "Site", site.URL)
	} else {
		slog.Info("Site removed, but no monitoring was active", "siteID", site.URL)
	}
}

func (m *Manager) EnableSite(siteID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if site, exists := m.sites[siteID]; exists {
		site.Enabled = true
		slog.Info("Site monitoring enabled", "site", site.URL)
	}
}

func (m *Manager) DisableSite(siteID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if site, exists := m.sites[siteID]; exists {
		site.Enabled = false
		slog.Info("Site monitoring disabled", "site", site.URL)
	}
}
