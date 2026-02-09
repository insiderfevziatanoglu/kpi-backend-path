package models

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type BalanceRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*Balance, error)
	UpdateBalance(ctx context.Context, userID int64, amount float64) error
	TransferFunds(ctx context.Context, fromID, toID int64, amount float64) error
}

type UserService interface {
	Register(ctx context.Context, user *User, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}