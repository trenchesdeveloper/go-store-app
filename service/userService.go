package service

import (
	"context"
	db "github.com/trenchesdeveloper/go-store-app/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/domain"
	"log"
)

type UserService struct {
	Store db.Store
}

func (us *UserService) findUserByEmail(email string) (*domain.User, error) {
	return &domain.User{}, nil
}

func (us *UserService) Register(ctx context.Context, user db.CreateUserParams) (string, error) {
	_, err := us.Store.CreateUser(ctx, user)

	if err != nil {
		return "", err

	}
	return "", nil
}

func (us *UserService) GetVerificationCode(user *domain.User) error {
	log.Println("Registering user.sql", user)
	return nil
}

func (us *UserService) Login(email, password string) (*domain.User, error) {
	return &domain.User{}, nil
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
