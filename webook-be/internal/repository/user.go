package repository

import (
	"context"

	"github.com/Twistzz2/webook/webook-be/internal/domain"
	"github.com/Twistzz2/webook/webook-be/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByID(int64) {
	// 先从 cache 找
	// 再从 dao 找
	// 找到了回写 cache
}
