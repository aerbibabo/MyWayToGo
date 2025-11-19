package repostory

import (
	"context"
	"github.com/aerbibabo/MyWayToGo/webook/internal/domain"
	"github.com/aerbibabo/MyWayToGo/webook/internal/repostory/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
