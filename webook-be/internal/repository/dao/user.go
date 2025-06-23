package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("邮箱已存在！")
	ErrUserNotFound       = gorm.ErrRecordNotFound
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
	Id       int64  `gorm:"primaryKey,autoIncrement"`      // 主键，自增
	Email    string `gorm:"type:varchar(255);uniqueIndex"` // 唯一索引，限制长度
	Password string

	// 毫秒数
	CreatedAt int64
	UpdatedAt int64
}

func (u *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (ud *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).First(&u, "id = ?", id).Error
	return u, err
}

func (u *UserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		switch mysqlErr := err.(type) {
		case *mysql.MySQLError:
			const uniqueConflictErrCode uint16 = 1062 // MySQL 错误码：唯一索引约束冲突
			if mysqlErr.Number == uniqueConflictErrCode {
				return ErrEmailAlreadyExists
			}
		default:
			// 其他错误直接返回
			return err
		}
	}
	return nil
}

// type Address struct {
// 	Id     int64
// 	UserId int64
// }
