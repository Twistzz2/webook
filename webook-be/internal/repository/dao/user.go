package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// User 直接对应数据库表结构，映射数据库
// 也称 entity, model, PO(Persistent Object)...
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"` // 主键，自动递增
	Email    string `gorm:"uniqueIndex"`              // 唯一索引，不能重复
	Password string

	// 毫秒数
	CreatedAt int64
	UpdatedAt int64
}

func (u *UserDAO) Insert(ctx context.Context, user User) error {
	// 这里调用 gorm 的方法
	// db.Create(&user)
	// return nil
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	return u.db.WithContext(ctx).Create(&user).Error
}

// type Address struct {
// 	Id     int64
// 	UserId int64
// }
