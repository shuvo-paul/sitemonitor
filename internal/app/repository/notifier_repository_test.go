package repository

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shuvo-paul/sitemonitor/internal/app/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNotifierRepository_Create(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewNotifierRepository(db)

	slackConfig := json.RawMessage(`{"webhook_url": "https://hooks.slack.com/test"}`)
	config := NotifierConfig{
		Type:   "slack",
		Config: slackConfig,
	}
	configJSON, err := json.Marshal(config)
	assert.NoError(t, err)

	mock.ExpectQuery("INSERT INTO notifiers").
		WithArgs(string(configJSON), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "config", "created_at", "updated_at"}).
			AddRow(1, string(configJSON), time.Now(), time.Now()))

	notifier, err := repo.Create(config)
	assert.NoError(t, err)
	assert.NotNil(t, notifier)
	assert.Equal(t, int64(1), notifier.ID)
	assert.Equal(t, config.Type, notifier.Config.Type)
	assert.JSONEq(t, string(slackConfig), string(notifier.Config.Config))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotifierRepository_Get(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewNotifierRepository(db)

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notifiers").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{}))

		notifier, err := repo.Get(1)
		assert.NoError(t, err)
		assert.Nil(t, notifier)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Found", func(t *testing.T) {
		config := NotifierConfig{
			Type:   "slack",
			Config: json.RawMessage(`{"webhook_url": "https://hooks.slack.com/test"}`),
		}
		configJSON, err := json.Marshal(config)
		assert.NoError(t, err)

		now := time.Now()
		mock.ExpectQuery("SELECT (.+) FROM notifiers").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "config", "created_at", "updated_at"}).
				AddRow(1, string(configJSON), now, now))

		notifier, err := repo.Get(1)
		assert.NoError(t, err)
		assert.NotNil(t, notifier)
		assert.Equal(t, int64(1), notifier.ID)
		assert.Equal(t, config.Type, notifier.Config.Type)
		assert.JSONEq(t, string(config.Config), string(notifier.Config.Config))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotifierRepository_List(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewNotifierRepository(db)

	t.Run("EmptyList", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notifiers").
			WillReturnRows(sqlmock.NewRows([]string{}))

		notifiers, err := repo.List()
		assert.NoError(t, err)
		assert.Empty(t, notifiers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MultipleNotifiers", func(t *testing.T) {
		configs := []NotifierConfig{
			{
				Type:   "slack",
				Config: json.RawMessage(`{"webhook_url": "https://hooks.slack.com/test1"}`),
			},
			{
				Type:   "email",
				Config: json.RawMessage(`{"smtp_host": "smtp.example.com"}`),
			},
		}

		rows := sqlmock.NewRows([]string{"id", "config", "created_at", "updated_at"})
		now := time.Now()

		for i, config := range configs {
			configJSON, err := json.Marshal(config)
			assert.NoError(t, err)
			rows.AddRow(i+1, string(configJSON), now, now)
		}

		mock.ExpectQuery("SELECT (.+) FROM notifiers").
			WillReturnRows(rows)

		notifiers, err := repo.List()
		assert.NoError(t, err)
		assert.Len(t, notifiers, len(configs))
		for i, notifier := range notifiers {
			assert.Equal(t, int64(i+1), notifier.ID)
			assert.Equal(t, configs[i].Type, notifier.Config.Type)
			assert.JSONEq(t, string(configs[i].Config), string(notifier.Config.Config))
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotifierRepository_Update(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewNotifierRepository(db)

	t.Run("NotFound", func(t *testing.T) {
		config := NotifierConfig{
			Type:   "slack",
			Config: json.RawMessage(`{"webhook_url": "https://hooks.slack.com/test"}`),
		}
		configJSON, err := json.Marshal(config)
		assert.NoError(t, err)

		mock.ExpectQuery("UPDATE notifiers").
			WithArgs(string(configJSON), sqlmock.AnyArg(), 999).
			WillReturnRows(sqlmock.NewRows([]string{}))

		notifier, err := repo.Update(999, config)
		assert.NoError(t, err)
		assert.Nil(t, notifier)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success", func(t *testing.T) {
		config := NotifierConfig{
			Type:   "slack",
			Config: json.RawMessage(`{"webhook_url": "https://hooks.slack.com/test"}`),
		}
		configJSON, err := json.Marshal(config)
		assert.NoError(t, err)

		now := time.Now()
		mock.ExpectQuery("UPDATE notifiers").
			WithArgs(string(configJSON), sqlmock.AnyArg(), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "config", "created_at", "updated_at"}).
				AddRow(1, string(configJSON), now, now))

		notifier, err := repo.Update(1, config)
		assert.NoError(t, err)
		assert.NotNil(t, notifier)
		assert.Equal(t, int64(1), notifier.ID)
		assert.Equal(t, config.Type, notifier.Config.Type)
		assert.JSONEq(t, string(config.Config), string(notifier.Config.Config))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotifierRepository_Delete(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewNotifierRepository(db)

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM notifiers").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM notifiers").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNotifierRepository_Errors(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	defer db.Close()

	repo := NewNotifierRepository(db)

	t.Run("Create Error", func(t *testing.T) {
		config := NotifierConfig{
			Type:   "slack",
			Config: json.RawMessage(`{"webhook_url": "https://hooks.slack.com/test"}`),
		}
		configJSON, err := json.Marshal(config)
		assert.NoError(t, err)

		mock.ExpectQuery("INSERT INTO notifiers").
			WithArgs(string(configJSON), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sqlmock.ErrCancelled)

		notifier, err := repo.Create(config)
		assert.Error(t, err)
		assert.Nil(t, notifier)
	})

	t.Run("Get Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notifiers").
			WithArgs(1).
			WillReturnError(sqlmock.ErrCancelled)

		notifier, err := repo.Get(1)
		assert.Error(t, err)
		assert.Nil(t, notifier)
	})

	t.Run("List Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notifiers").
			WillReturnError(sqlmock.ErrCancelled)

		notifiers, err := repo.List()
		assert.Error(t, err)
		assert.Nil(t, notifiers)
	})

	t.Run("Update Error", func(t *testing.T) {
		config := NotifierConfig{
			Type:   "slack",
			Config: json.RawMessage(`{"webhook_url": "https://hooks.slack.com/test"}`),
		}
		configJSON, err := json.Marshal(config)
		assert.NoError(t, err)

		mock.ExpectQuery("UPDATE notifiers").
			WithArgs(string(configJSON), sqlmock.AnyArg(), 1).
			WillReturnError(sqlmock.ErrCancelled)

		notifier, err := repo.Update(1, config)
		assert.Error(t, err)
		assert.Nil(t, notifier)
	})

	t.Run("Delete Error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM notifiers").
			WithArgs(1).
			WillReturnError(sqlmock.ErrCancelled)

		err := repo.Delete(1)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
