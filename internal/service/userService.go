package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/trenchesdeveloper/go-store-app/config"
	db2 "github.com/trenchesdeveloper/go-store-app/internal/db/sqlc"
	"github.com/trenchesdeveloper/go-store-app/internal/dto"
	"github.com/trenchesdeveloper/go-store-app/internal/helper"
	"github.com/trenchesdeveloper/go-store-app/pkg/notification"
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

func (us *UserService) UpdateUser(ctx context.Context, id uint, input dto.UpdateUserRequest) (db2.UpdateUserRow,error) {
	user, err := us.Store.UpdateUser(ctx, db2.UpdateUserParams{
		ID:        int32(id),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone: pgtype.Text{
			String: input.Phone,
			Valid:  true,
		},
	})

	if err != nil {
		return db2.UpdateUserRow{}, fmt.Errorf("could not update user: %v", err)
	}

	return user, nil
}

func (us *UserService) CreateProfile(ctx context.Context, id uint, input dto.ProfileInput) error {
	_, err := us.Store.CreateAddress(ctx, db2.CreateAddressParams{
		UserID:       int32(id),
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: pgtype.Text{
			String: input.AddressInput.AddressLine2,
			Valid:  true,
		},
		City:         input.AddressInput.City,
		State:        input.AddressInput.State,
		Country:      input.AddressInput.Country,
		PostCode:     int32(input.AddressInput.PostCode),
	})

	if err != nil {
		return fmt.Errorf("could not create address: %v", err)
	}

	return nil
}

func (us *UserService) UpdateProfile(ctx context.Context, userID uint, input dto.ProfileInput) error {
	_, err := us.Store.UpdateAddress(ctx, db2.UpdateAddressParams{
		UserID:       int32(userID),
		UserID_2:    int32(userID),
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: pgtype.Text{
			String: input.AddressInput.AddressLine2,
			Valid:  true,
		},
		City:         input.AddressInput.City,
		State:        input.AddressInput.State,
		Country:      input.AddressInput.Country,
		PostCode:     int32(input.AddressInput.PostCode),
	})

	if err != nil {
		return fmt.Errorf("could not update address: %v", err)
	}

	return nil
}


func (us *UserService) DeleteUser(id uint) error {
	return nil
}

func (us *UserService) GetUserByID(id uint) (dto.ProfileInput,error) {
	// get the user and return the dto.ProfileInput
	user, err := us.Store.GetUser(context.Background(), int32(id))

	if err != nil {
		return dto.ProfileInput{}, errors.New("user not found")
	}

	address, err := us.Store.FindAddressByUser(context.Background(), int32(id))

	if err != nil {
		return dto.ProfileInput{}, errors.New("address not found")
	}

	return dto.ProfileInput{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AddressInput: dto.AddressInput{
			AddressLine1: address[0].AddressLine1,
			AddressLine2: address[0].AddressLine2.String,
			City:         address[0].City,
			State:        address[0].State,
			Country:      address[0].Country,
			PostCode:     uint(address[0].PostCode),
		},
	}, nil
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

func (us *UserService) CreateCart(ctx context.Context, userID uint, input dto.CreateCartRequest) ([]db2.Cart, error) {
	// check if the cart exists
	cart, _ := us.Store.FindCartItem(ctx, db2.FindCartItemParams{
		UserID:    int32(userID),
		ProductID: int32(input.ProductID),
	})
	if cart.ID > 0 {
		if input.ProductID == 0 {
			return nil, errors.New("please provide a valid product id")
		}
		if input.Quantity < 1 {
			err := us.Store.DeleteCartById(ctx, cart.ID)
			if err != nil {
				return nil, fmt.Errorf("could not delete cart: %v", err)
			}
		} else {
			// update cart item

			_, err := us.Store.UpdateCart(ctx, db2.UpdateCartParams{
				ID:       cart.ID,
				Quantity: int32(input.Quantity),
			})

			if err != nil {
				return nil, fmt.Errorf("could not update cart: %v", err)
			}

		}

	} else {

		// get the product
		product, err := us.Store.GetProductByID(ctx, int32(input.ProductID))

		if err != nil {
			return nil, errors.New("product not found")
		}

		// create cart
		_, err = us.Store.CreateCart(ctx, db2.CreateCartParams{
			UserID:    int32(userID),
			ProductID: int32(input.ProductID),
			Quantity:  int32(input.Quantity),
			Price:     product.Price,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl.String,
			SellerID:  product.UserID,
		})

		if err != nil {
			return nil, fmt.Errorf("could not create cart: %v", err)
		}
	}

	return us.Store.FindCartItems(ctx, int32(userID))
}

func (us *UserService) GetCart(ctx context.Context, userID uint) ([]db2.Cart, error) {
	return us.Store.FindCartItems(ctx, int32(userID))
}

func (us *UserService) GetOrders(id uint) error {
	return nil
}
