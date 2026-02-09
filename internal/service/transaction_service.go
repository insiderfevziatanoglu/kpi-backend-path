package service

import (
	"context"
	"errors"
	"time"

	"github.com/fevziatanoglu/test-go-project/internal/models"
)

type TransactionManager interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
}

type TransactionService struct {
	repo       TransactionRepository
	balanceRepo BalanceRepository
	txManager  TransactionManager
}

func NewTransactionService(repo TransactionRepository, balRepo BalanceRepository, tm TransactionManager) *TransactionService {
	return &TransactionService{
		repo:        repo,
		balanceRepo: balRepo,
		txManager:   tm,
	}
}

func (s *TransactionService) Transfer(ctx context.Context, fromID, toID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	return s.txManager.ExecTx(ctx, func(ctx context.Context) error {
		
		fromBalance, err := s.balanceRepo.GetByUserID(ctx, fromID)
		if err != nil {
			return err
		}

		if !fromBalance.Withdraw(amount) {
			return errors.New("insufficient balance")
		}

		toBalance, err := s.balanceRepo.GetByUserID(ctx, toID)
		if err != nil {
			return err
		}

		toBalance.Deposit(amount)

		if err := s.balanceRepo.Update(ctx, fromBalance); err != nil {
			return err
		}
		if err := s.balanceRepo.Update(ctx, toBalance); err != nil {
			return err
		}

		txRecord := &models.Transaction{
			FromUserID: &fromID,
			ToUserID:   &toID,
			Amount:     amount,
			Type:       "TRANSFER",
			Status:     "COMPLETED",
			CreatedAt:  time.Now(),
		}
		
		if err := s.repo.Create(ctx, txRecord); err != nil {
			return err
		}

		return nil
	})
}