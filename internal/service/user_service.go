package service

import (
	"context"
	"errors"
	"time"

	"github.com/fevziatanoglu/test-go-project/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, req *models.User, password string) error {
	if err := req.Validate(); err != nil {
		return err
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("password could not be secured")
	}
	req.PasswordHash = string(hashedPwd)
	
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	if req.Role == "" {
		req.Role = "user" 
	}

	return s.repo.Create(ctx, req)
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found or invalid password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}