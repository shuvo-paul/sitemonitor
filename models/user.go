package models

import (
	"fmt"

	"github.com/shuvo-paul/sitemonitor/config"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string // This will store the hashed password
}

func (u *User) Save() error {
	if config.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	_, err := config.DB.Exec(query, u.Username, u.Email, u.Password)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func Register(username, email, password string) (*User, error) {
	// Check if the email already exists
	exists, err := EmailExists(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already in use")
	}

	user := &User{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return user, nil
}

func EmailExists(email string) (bool, error) {
	if config.DB == nil {
		return false, fmt.Errorf("database connection is not initialized")
	}

	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	var count int
	err := config.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to query email: %w", err)
	}

	return count > 0, nil
}
