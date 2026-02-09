package models

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(u.Email) {
		return errors.New("not a valid email address")
	}
	return nil
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type Transaction struct {
	ID         int64     `json:"id"`
	FromUserID *int64    `json:"from_user_id,omitempty"`
	ToUserID   *int64    `json:"to_user_id,omitempty"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`   
	Status     string    `json:"status"` 
	CreatedAt  time.Time `json:"created_at"`
}

func (t *Transaction) Complete() {
	t.Status = "COMPLETED"
}

func (t *Transaction) Fail() {
	t.Status = "FAILED"
}

type Balance struct {
	UserID        int64     `json:"user_id"`
	Amount        float64   `json:"amount"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}