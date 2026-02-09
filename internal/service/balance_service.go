package service

import (
	"context"
	"errors"

	"github.com/fevziatanoglu/test-go-project/internal/models"
)

type BalanceRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*models.Balance, error)
	Update(ctx context.Context, balance *models.Balance) error
}

type BalanceService struct {
	repo BalanceRepository
}

func NewBalanceService(repo BalanceRepository) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *BalanceService) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	balance, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if amount > 0 {
		balance.Deposit(amount)
	} else {
		if !balance.Withdraw(-amount) {
			return errors.New("insufficient balance")
		}
	}

	return s.repo.Update(ctx, balance)
}