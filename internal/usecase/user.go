package service

import (
	"context"
	"errors"
	"github.com/Enthreeka/tg-posting-bot/internal/entity"
	"github.com/Enthreeka/tg-posting-bot/internal/repo"
	"github.com/Enthreeka/tg-posting-bot/pkg/logger"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetAllAdmin(ctx context.Context) ([]entity.User, error)

	CreateUserIFNotExist(ctx context.Context, user *entity.User) error

	UpdateRoleByUsername(ctx context.Context, role entity.UserRole, username string) error
}

type userService struct {
	userRepo repo.UserRepo
	log      *logger.Logger
}

func NewUserService(userRepo repo.UserRepo, log *logger.Logger) (UserService, error) {
	if userRepo == nil {
		return nil, errors.New("userRepo is nil")
	}
	if log == nil {
		return nil, errors.New("log is nil")
	}

	return &userService{
		userRepo: userRepo,
		log:      log,
	}, nil
}

func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *userService) GetAllAdmin(ctx context.Context) ([]entity.User, error) {
	return u.userRepo.GetAllAdmin(ctx)
}

func (u *userService) CreateUserIFNotExist(ctx context.Context, user *entity.User) error {
	isExist, err := u.userRepo.IsUserExistByUserID(ctx, user.ID)
	if err != nil {
		u.log.Error("userRepo.IsUserExistByUsernameTg: failed to check user: %v", err)
		return err
	}

	if !isExist {
		u.log.Info("GetPub user: %s", user.String())
		err := u.userRepo.CreateUser(ctx, user)
		if err != nil {
			u.log.Error("userRepo.CreateUser: failed to create user: %v", err)
			return err
		}
	}

	return nil
}

func (u *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *userService) UpdateRoleByUsername(ctx context.Context, role entity.UserRole, username string) error {
	return u.userRepo.UpdateRoleByUsername(ctx, role, username)
}
