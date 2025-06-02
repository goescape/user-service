package jwt

import (
	"fmt"
	"net/http"
	"time"
	"user-svc/helpers/fault"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	tokenExpiry        = 1024 * time.Minute // Durasi berlaku token akses
	refreshTokenExpiry = 72 * time.Hour   // Durasi berlaku refresh token
)

var signedKey = []byte("secret") // Kunci rahasia untuk tanda tangan JWT

// JWTPayload Struktur payload JWT berisi info user dan klaim standar JWT
type JWTPayload struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

// CreateAccessToken Buat token akses dengan durasi tokenExpiry
func CreateAccessToken(name, email, userId string) (*string, *JWTPayload, error) {
	return generateToken(name, email, userId, tokenExpiry)
}

// CreateRefreshToken Buat refresh token dengan durasi refreshTokenExpiry
func CreateRefreshToken(name, email, userId string) (*string, *JWTPayload, error) {
	return generateToken(name, email, userId, refreshTokenExpiry)
}

// Generate token JWT dengan payload dan tanda tangan menggunakan HS256
func generateToken(name, email, userId string, duration time.Duration) (*string, *JWTPayload, error) {
	payload, err := newJWTPayload(name, email, userId, duration)
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

// Buat payload JWT baru dengan klaim yang lengkap dan expired time
func newJWTPayload(name, email, userId string, duration time.Duration) (*JWTPayload, error) {
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
		Name:   name,
		Email:  email,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user_login",            // Pengeluarnya token
			Subject:   "go-escape",             // Subjek token
			ID:        tokenID.String(),        // ID unik token
			IssuedAt:  jwt.NewNumericDate(now), // Waktu dibuat token
			NotBefore: jwt.NewNumericDate(now), // Token berlaku sejak waktu ini
			ExpiresAt: jwt.NewNumericDate(exp), // Waktu kadaluarsa token
		},
	}, nil
}

// GetClaims Ambil klaim dari token string, cek validitas dan signature
func GetClaims(token string) (*JWTPayload, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		return signedKey, nil
	})
	if err != nil {
		return nil, fault.Custom(
			http.StatusUnauthorized,
			fault.ErrUnauthorized,
			fmt.Sprintf("failed to parse token: %v", err.Error()),
		)
	}

	if claims, ok := parsedToken.Claims.(*JWTPayload); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fault.Custom(
		http.StatusUnauthorized,
		fault.ErrUnauthorized,
		"invalid token",
	)
}

// VerifyToken Verifikasi token dan kembalikan klaim jika valid
func VerifyToken(tokenString string) (*JWTPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		return signedKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTPayload)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
