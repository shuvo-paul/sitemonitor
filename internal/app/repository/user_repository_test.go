package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shuvo-paul/sitemonitor/internal/app/models"
	"github.com/shuvo-paul/sitemonitor/internal/app/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	userRepo := NewUserRepository(db)
	defer db.Close()

	user := &models.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	savedUser, err := userRepo.SaveUser(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, savedUser.ID)
	assert.Equal(t, user.Name, savedUser.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEmailExists(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	userRepo := NewUserRepository(db)
	defer db.Close()

	t.Run("email exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("existing@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := userRepo.EmailExists("existing@example.com")

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("email does not exist", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("nonexistent@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := userRepo.EmailExists("nonexistent@example.com")

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUser(t *testing.T) {
	db, mock := testutil.SetupTestDB(t)
	userRepo := NewUserRepository(db)
	defer db.Close()
	expectedUser := &models.User{
		ID:       1,
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	t.Run("By Email: user found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password"}).
			AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password)

		mock.ExpectQuery("SELECT id, name, email, password FROM users").
			WithArgs(expectedUser.Email).
			WillReturnRows(rows)

		user, err := userRepo.GetUserByEmail(expectedUser.Email)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("By ID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email)

		mock.ExpectQuery("SELECT id, name, email from users").
			WithArgs(expectedUser.ID).
			WillReturnRows(rows)

		repo := NewUserRepository(db)

		user, err := repo.GetUserByID(expectedUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Name, user.Name)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
