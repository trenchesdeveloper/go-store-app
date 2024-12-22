package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/trenchesdeveloper/go-store-app/config"
	db2 "github.com/trenchesdeveloper/go-store-app/internal/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/dto"
	"github.com/trenchesdeveloper/go-store-app/internal/helper"
	"github.com/trenchesdeveloper/go-store-app/pkg/notification"
	"strconv"
	"time"
)

type UserService struct {
	Store  db2.Store
	Auth   helper.Auth
	Config config.AppConfig
}

func (us *UserService) findUserByEmail(ctx context.Context, email string) (db2.User, error) {
	user, err := us.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return db2.User{}, err

	}

	return user, nil
}

func (us *UserService) Register(ctx context.Context, params db2.CreateUserParams) (string, error) {
	hashPassword, err := us.Auth.HashPassword(params.Password)

	if err != nil {
		return "", err
	}

	params.Password = hashPassword

	createdUser, err := us.Store.CreateUser(ctx, params)

	if err != nil {
		if db2.ErrorCode(err) == db2.UniqueViolation {

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

func (us *UserService) isVerifiedUser(ctx *fiber.Ctx, id uint) bool {
	currentUser, err := us.Store.GetUser(ctx.Context(), int32(id))

	if err != nil {
		return false
	}

	return currentUser.Verified
}

func (us *UserService) GetVerificationCode(ctx *fiber.Ctx, currentUser helper.TokenPayload) error {

	// check if user is already verified
	if us.isVerifiedUser(ctx, currentUser.ID) {
		return errors.New("user already verified ")
	}

	// generate verification code
	code, err := us.Auth.GenerateCode()

	if err != nil {
		return errors.New("could not generate code")
	}

	// update user with verification code
	userParams := db2.UpdateUserCodeAndExpiryParams{
		ID: int32(currentUser.ID),
		Expiry: pgtype.Timestamp{
			Time:  time.Now().Add(30 * time.Minute),
			Valid: true,
		},

		Code: pgtype.Text{
			String: strconv.Itoa(code),
			Valid:  true,
		},
	}

	updatedUser, err := us.Store.UpdateUserCodeAndExpiry(ctx.Context(), userParams)

	if err != nil {
		return fmt.Errorf("could not update user: %v", err)

	}

	// send sms/verification code to user
	smsClient := notification.NewNotificationClient(us.Config)

	err = smsClient.SendSMS(updatedUser.Phone.String, fmt.Sprintf("Your verification code is %d", code))

	if err != nil {
		return fmt.Errorf("could not send sms: %v", err)
	}

	return nil
}

func (us *UserService) VerifyUser(ctx *fiber.Ctx, userID uint, code int) error {
	// check if user is already verified
	if us.isVerifiedUser(ctx, userID) {
		return errors.New("user already verified ")
	}

	// get user
	user, err := us.Store.GetUser(ctx.Context(), int32(userID))

	if err != nil {
		return errors.New("user not found")

	}

	if user.Code.String != strconv.Itoa(code) {
		return errors.New("invalid code")
	}

	if user.Expiry.Time.Before(time.Now()) {
		return errors.New("verification code expired")
	}

	// update user
	userParams := db2.UpdateUserVerifiedParams{
		ID:       int32(userID),
		Verified: true,
	}

	_, err = us.Store.UpdateUserVerified(ctx.Context(), userParams)

	if err != nil {
		return fmt.Errorf("could not update user: %v", err)

	}

	return nil
}

func (us *UserService) UpdateUser() error {
	return nil
}

func (us *UserService) DeleteUser(id uint) error {
	return nil
}

func (us *UserService) GetUserByID(id uint) error {
	return nil
}

func (us *UserService) BecomeSeller(ctx *fiber.Ctx, userid uint, input dto.SellerInput) (string, error) {
	// find existing user
	user, err := us.Store.GetUser(ctx.Context(), int32(userid))

	if err != nil {
		return "", errors.New("user not found")

	}
	// check if user is already a seller
	if user.UserType == db2.UserTypeSeller {
		return "", errors.New("user is already a seller")
	}

	// update user to seller
	seller, err := us.Store.UpdateUserToSeller(ctx.Context(), db2.UpdateUserToSellerParams{
		ID:        int32(userid),
		UserType:  db2.UserTypeSeller,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone: pgtype.Text{
			String: input.PhoneNumber,
			Valid:  true,
		},
	})

	if err != nil {
		return "", fmt.Errorf("could not update user: %v", err)
	}

	// create bank account
	_, err = us.Store.CreateBankAccount(ctx.Context(), db2.CreateBankAccountParams{
		UserID:      int64(seller.ID),
		BankAccount: int64(input.BankAccountNumber),
		SwiftCode: pgtype.Text{
			String: input.SwiftCode,
			Valid:  true,
		},
		PaymentType: pgtype.Text{
			String: input.PaymentType,
			Valid:  true,
		},
	})

	if err != nil {
		return "", fmt.Errorf("could not create bank account: %v", err)
	}

	//generate token
	tokenPayload := helper.TokenPayload{
		ID:    uint(seller.ID),
		Email: seller.Email,
		Role:  string(seller.UserType),
	}

	return us.Auth.GenerateToken(tokenPayload)

}

func (us *UserService) CreateCart(id uint) error {
	return nil
}

func (us *UserService) GetOrders(id uint) error {
	return nil
}
