package service

import (
	"context"
	"errors"
	db "github.com/trenchesdeveloper/go-store-app/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/domain"
	"github.com/trenchesdeveloper/go-store-app/internal/helper"
	"log"
)

type UserService struct {
	Store db.Store
	Auth  helper.Auth
}

func (us *UserService) findUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := us.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err

	}

	return user, nil
}

func (us *UserService) Register(ctx context.Context, params db.CreateUserParams) (string, error) {
	hashPassword, err := us.Auth.HashPassword(params.Password)

	if err != nil {
		return "", err
	}

	params.Password = hashPassword

	createdUser, err := us.Store.CreateUser(ctx, params)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {

			return "", errors.New("email already exists")
		}
		return "", err
	}

	// generate token
	return us.Auth.GenerateToken(helper.TokenPayload{
		ID:    uint(createdUser.ID),
		Email: createdUser.Email,
		Role:  string(createdUser.UserType),
	})

}

func (us *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := us.findUserByEmail(ctx, email)

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = us.Auth.ComparePassword(user.Password, password)

	if err != nil {
		return "", errors.New("invalid credentials")

	}

	return us.Auth.GenerateToken(helper.TokenPayload{
		ID:    uint(user.ID),
		Email: user.Email,
		Role:  string(user.UserType),
	})

}

func (us *UserService) GetVerificationCode(user *domain.User) error {
	log.Println("Registering user.sql", user)
	return nil
}

func (us *UserService) VerifyUser(code string) error {
	return nil
}

func (us *UserService) UpdateUser(user *domain.User) error {
	return nil
}

func (us *UserService) DeleteUser(id uint) error {
	return nil
}

func (us *UserService) GetUserByID(id uint) (*domain.User, error) {
	return &domain.User{}, nil
}

func (us *UserService) BecomeSeller(id uint) (*domain.User, error) {
	return &domain.User{}, nil
}

func (us *UserService) CreateCart(id uint) (*domain.User, error) {
	return &domain.User{}, nil
}

func (us *UserService) GetOrders(id uint) (*domain.User, error) {
	return &domain.User{}, nil
}
