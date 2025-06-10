package service

import (
	"context"

	"github.com/Twistzz2/webook/webook-be/internal/domain"
	"github.com/Twistzz2/webook/webook-be/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 加密
	// 存起来
	return svc.repo.Create(ctx, u)
}
