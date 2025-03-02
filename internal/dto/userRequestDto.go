package dto

import "github.com/jackc/pgx/v5/pgtype"

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegister struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Phone     string `json:"phone"`
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
}

type UserResponse struct {
	ID        int32       `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Phone     pgtype.Text `json:"phone"`
	Verified  bool        `json:"verified"`
}

type VerificationCodeInput struct {
	Code int `json:"code"`
}

type SellerInput struct {
	FirstName         string `json:"first_name" validate:"required"`
	LastName          string `json:"last_name" validate:"required"`
	PhoneNumber       string `json:"phone_number" validate:"required"`
	BankAccountNumber uint64 `json:"bank_account_number" validate:"required"`
	SwiftCode         string `json:"swift_code" validate:"required"`
	PaymentType       string `json:"payment_type" validate:"required"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type AddressInput struct {
	AddressLine1 string `json:"address_line1" validate:"required"`
	AddressLine2 string `json:"address_line2"`
	City		 string `json:"city" validate:"required"`
	State		 string `json:"state" validate:"required"`
	Country		 string `json:"country" validate:"required"`
	PostCode	 uint `json:"post_code" validate:"required"`
}

type ProfileInput struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	AddressInput AddressInput `json:"address" validate:"required"`
}
