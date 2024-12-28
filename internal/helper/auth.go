package helper

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	db "github.com/trenchesdeveloper/go-store-app/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

const (
	authorizationPayloadKey = "authorization_payload"
)

type TokenPayload struct {
	ID    uint
	Email string
	Role  string
}

type Auth struct {
	Secret string
}

func NewAuth(secret string) Auth {
	return Auth{
		Secret: secret,
	}

}

func (a *Auth) HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password must be at least 6 characters")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		// log actual error and report a generic error to the user
		return "", errors.New("something went wrong")
	}

	return string(hashPassword), nil
}

func (a *Auth) ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}

func (a *Auth) GenerateToken(payload TokenPayload) (string, error) {
	if payload.ID == 0 || payload.Email == "" || payload.Role == "" {
		return "", errors.New("invalid payload")

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": payload.ID,
		"email":   payload.Email,
		"role":    payload.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.Secret))

	if err != nil {
		log.Println("could not sign the token")
		return "", errors.New("")
	}

	return tokenString, nil
}

func (a *Auth) VerifyToken(token string) (TokenPayload, error) {
	if token == "" {
		return TokenPayload{}, errors.New("token is empty")

	}

	tokenArr := strings.Split(token, " ")

	if len(tokenArr) != 2 {
		return TokenPayload{}, errors.New("invalid token")
	}

	if tokenArr[0] != "Bearer" {
		return TokenPayload{}, errors.New("invalid token")
	}

	token = tokenArr[1]

	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")

		}

		return []byte(a.Secret), nil

	})

	if err != nil {
		return TokenPayload{}, errors.New("invalid token format")

	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		if claims["exp"].(float64) < float64(time.Now().Unix()) {
			return TokenPayload{}, errors.New("please login again")
		}
		payload := TokenPayload{}

		payload.ID = uint(claims["user_id"].(float64))
		payload.Email = claims["email"].(string)
		payload.Role = claims["role"].(string)

		return payload, nil

	}

	return TokenPayload{}, errors.New("invalid token")
}

func (a *Auth) Authorize(ctx *fiber.Ctx) error {
	authHeader := ctx.GetReqHeaders()["Authorization"]

	if len(authHeader) == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	payload, err := a.VerifyToken(authHeader[0])

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	ctx.Locals(authorizationPayloadKey, payload)

	return ctx.Next()
}

func (a *Auth) GetCurrentUser(ctx *fiber.Ctx) (TokenPayload, error) {
	user := ctx.Locals(authorizationPayloadKey)

	return user.(TokenPayload), nil
}

func (a Auth) GenerateCode() (int, error) {
	return RandomNumbers(6)
}


func (a Auth) AuthorizeSeller(ctx *fiber.Ctx) error {
	authHeader := ctx.GetReqHeaders()["Authorization"]

	if len(authHeader) == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message": "unauthorized",
		})
	}

	payload, err := a.VerifyToken(authHeader[0])

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	if payload.Role != string(db.UserTypeSeller) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"message": "unauthorized",
		})
	}

	ctx.Locals(authorizationPayloadKey, payload)

	return ctx.Next()
}