package jwt

import (
	"net/http"
	"time"
	"user-svc/helpers/fault"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var signedKey = []byte("secret")

type JWTPayload struct {
	Email  string `json:"email"`
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateAccessToken(email, userId string, tokenExpiry time.Duration) (*string, *JWTPayload, error) {
	return generateToken(email, userId, tokenExpiry)
}

func CreateRefreshToken(email, userId string, tokenExpiry time.Duration) (*string, *JWTPayload, error) {
	return generateToken(email, userId, tokenExpiry)
}

func generateToken(email, userId string, duration time.Duration) (*string, *JWTPayload, error) {
	payload, err := newJWTPayload(email, userId, duration)
	if err != nil {
		return nil, nil, err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(signedKey)
	if err != nil {
		return nil, nil, fault.Custom(
			http.StatusUnprocessableEntity,
			fault.ErrUnprocessable,
			"failed signing JWT token: "+err.Error(),
		)
	}

	return &token, payload, nil
}

func newJWTPayload(email, userId string, duration time.Duration) (*JWTPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fault.Custom(
			http.StatusUnprocessableEntity,
			fault.ErrUnprocessable,
			"failed to generate token ID: "+err.Error(),
		)
	}

	now := time.Now()
	exp := now.Add(duration)

	return &JWTPayload{
		Email:  email,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user_login",
			Subject:   "go-escape",
			ID:        tokenID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}, nil
}
