package service

import (
	"fmt"
	"testing"

	"github.com/shuvo-paul/uptimebot/internal/auth/model"
	"github.com/stretchr/testify/assert"
)

// MockUserRepository is a mock implementation of UserRepository
type mockUserRepository struct {
	saveUserFunc       func(user *model.User) (*model.User, error)
	emailExistsFunc    func(email string) (bool, error)
	getUserByEmailFunc func(email string) (*model.User, error)
	getUserByIdFunc    func(id int) (*model.User, error)
}

func (m *mockUserRepository) SaveUser(user *model.User) (*model.User, error) {
	return m.saveUserFunc(user)
}

func (m *mockUserRepository) EmailExists(email string) (bool, error) {
	return m.emailExistsFunc(email)
}

func (m *mockUserRepository) GetUserByEmail(email string) (*model.User, error) {
	return m.getUserByEmailFunc(email)
}

func (m *mockUserRepository) GetUserByID(id int) (*model.User, error) {
	return m.getUserByIdFunc(id)
}

func TestCreateUser(t *testing.T) {
	mockRepo := &mockUserRepository{
		saveUserFunc: func(user *model.User) (*model.User, error) {
			user.ID = 1
			return user, nil
		},
	}
	userService := NewAuthService(mockRepo)
	t.Run("User created successfully", func(t *testing.T) {
		mockRepo.emailExistsFunc = func(email string) (bool, error) {
			return false, nil
		}

		user := &model.User{
			Name:     "testuser",
			Email:    "test@example.com",
			Password: "password123@",
		}
		savedUser, err := userService.CreateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, savedUser.ID, 1)
	})

	t.Run("Email Exist", func(t *testing.T) {
		mockRepo.emailExistsFunc = func(email string) (bool, error) {
			return true, fmt.Errorf("email already exists")
		}
		user := &model.User{
			Name:     "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		_, err := userService.CreateUser(user)
		assert.Error(t, err)
	})

}

func TestAuthentication(t *testing.T) {
	email := "test@example.com"
	password := "password123"
	wrongPassword := "wrongpassword123"

	user := &model.User{
		Name:     "testuser",
		Email:    email,
		Password: password,
	}

	user.HashPassword()

	mockRepo := &mockUserRepository{
		getUserByEmailFunc: func(email string) (*model.User, error) {

			return user, nil
		},
	}

	userService := NewAuthService(mockRepo)

	t.Run("Logged in succesfully", func(t *testing.T) {
		user, err := userService.Authenticate(email, password)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, email)
	})

	t.Run("Login failed", func(t *testing.T) {
		_, err := userService.Authenticate(email, wrongPassword)
		assert.Error(t, err)
	})
}
