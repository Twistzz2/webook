package repository

import (
	"context"

	"github.com/Twistzz2/webook/webook-be/internal/domain"
	"github.com/Twistzz2/webook/webook-be/internal/repository/dao"
)

var (
	ErrEmailAlreadyExists = dao.ErrEmailAlreadyExists
	// repository 并不知道 ErrUserNotFound 用的是 gorm.ErrRecordNotFound
	ErrUserNotFound = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (ur *UserRepository) FindById(ctx context.Context,
	id int64) (domain.User, error) {
	u, err := ur.dao.FindById(ctx, id)
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
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
