package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (ud *UserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	err := ud.db.WithContext(ctx).Create(&user).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}

	return err
}

func (ud *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	Email    string `gorm:"unique"`
	Password string

	Ctime int64
	Utime int64
}
