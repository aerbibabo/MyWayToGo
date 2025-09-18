package service

import (
	"context"
	"errors"
	"github.com/aerbibabo/MyWayToGo/webbook/internal/domain"
	"github.com/aerbibabo/MyWayToGo/webbook/internal/repostory"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repostory.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("invalid user/password")

type UserService struct {
	repo *repostory.UserRepository
}

func NewUserService(repo *repostory.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}
