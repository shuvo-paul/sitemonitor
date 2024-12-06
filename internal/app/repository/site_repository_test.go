package repository

import (
	"testing"
	"time"

	"github.com/shuvo-paul/sitemonitor/pkg/monitor"
	"github.com/stretchr/testify/assert"
)

func TestSiteRepository(t *testing.T) {
	siteRepo := NewSiteRepository(db)

	site := &monitor.Site{
		URL:             "example.org",
		Status:          "up",
		Enabled:         false,
		Interval:        30 * time.Second,
		StatusChangedAt: time.Now(),
	}
	t.Run("create", func(t *testing.T) {
		newSite, err := siteRepo.Create(site)
		assert.NoError(t, err)
		assert.Equal(t, site.URL, newSite.URL)
		assert.NotZero(t, newSite.ID)
	})
}
