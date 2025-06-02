package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"user-svc/helpers/cache"
	"user-svc/helpers/fault"
	"user-svc/helpers/jwt"
	"user-svc/middlewares"
	"user-svc/model"
	repository "user-svc/repository/user"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// userUsecase mengimplementasikan logic bisnis terkait user
type userUsecase struct {
	user  repository.UserRepository // Interface ke DB
	redis *redis.Client             // Redis client untuk caching
}

// Constructor untuk userUsecase
func NewUserUsecase(repository repository.UserRepository, redis *redis.Client) *userUsecase {
	return &userUsecase{
		user:  repository,
		redis: redis,
	}
}

// Interface untuk semua fungsi yang tersedia di usecase user
type UserUsecases interface {
	Register(body model.RegisterUser) (*model.LoginResponse, error)
	Login(body model.UserLogin) (*model.LoginResponse, error)
}

// Register mendaftarkan user baru atau mengembalikan token jika sudah ada cache-nya
func (u *userUsecase) Register(body model.RegisterUser) (*model.LoginResponse, error) {
	ctx := context.TODO()
	cacheKey := fmt.Sprintf("login:%s", body.Email)

	// Cek apakah data user sudah ada di Redis
	cacheExist, err := cache.Exist(ctx, u.redis, cacheKey)
	if err != nil {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to check Redis key existence for '%s': %v", cacheKey, err),
		)
	}

	// Jika ada cache login, langsung ambil dari Redis
	if cacheExist {
		tokenValue, err := cache.Get(ctx, u.redis, cacheKey)
		if err != nil {
			return nil, fault.Custom(
				http.StatusInternalServerError,
				fault.ErrInternalServer,
				fmt.Sprintf("failed to retrieve access token from Redis for key '%s': %v", cacheKey, err),
			)
		}

		var res model.LoginResponse
		err = json.Unmarshal([]byte(tokenValue.(string)), &res)
		if err != nil {
			return nil, fault.Custom(
				http.StatusInternalServerError,
				fault.ErrInternalServer,
				fmt.Sprintf("failed to unmarshal cached access token for key '%s': %v", cacheKey, err),
			)
		}

		return &res, nil
	}

	// Cek apakah user sudah pernah terdaftar berdasarkan nama
	exist, err := u.user.ExistsByName(body.Name)
	if err != nil {
		return nil, err
	}

	var user *model.User
	var userId uuid.UUID

	// Jika belum ada, daftarkan user baru
	if !exist {
		body.Password = middlewares.GeneratePassword(body.Password)

		createdId, err := u.user.Insert(body)
		if err != nil {
			return nil, err
		}
		userId = *createdId

		// Ambil detail user baru berdasarkan ID
		user, err = u.user.Detail(model.GetUserDetailRequest{
			UserId: userId,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Kalau sudah ada user, ambil datanya berdasarkan email
		user, err = u.user.Detail(model.GetUserDetailRequest{
			Email: body.Email,
		})
		if err != nil {
			return nil, err
		}
		userId = user.Id
	}

	// Generate JWT access & refresh token
	accessToken, payload, err := jwt.CreateAccessToken(user.Name, user.Email, userId.String())
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := jwt.CreateRefreshToken(user.Name, user.Email, userId.String())
	if err != nil {
		return nil, err
	}

	// Bentuk response
	res := &model.LoginResponse{
		UserData:              *user,
		AccessToken:           *accessToken,
		AccessTokenExpiresAt:  &payload.ExpiresAt.Time,
		RefreshToken:          *refreshToken,
		RefreshTokenExpiresAt: &refreshPayload.ExpiresAt.Time,
	}

	// Simpan response ke Redis sebagai cache
	jsonValue, err := json.Marshal(res)
	if err != nil {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to marshal login response to JSON for key '%s': %v", cacheKey, err),
		)
	}

	err = cache.Set(ctx, u.redis, cacheKey, string(jsonValue), 10*time.Minute)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Login memverifikasi kredensial user dan mengembalikan token login
func (u *userUsecase) Login(body model.UserLogin) (*model.LoginResponse, error) {
	ctx := context.TODO()

	// Ambil detail user berdasarkan email
	user, err := u.user.Detail(model.GetUserDetailRequest{
		Email: body.Email,
	})
	if err != nil {
		return nil, err
	}

	// Jika user tidak ditemukan
	if user == nil {
		return nil, fault.Custom(
			http.StatusNotFound,
			fault.ErrNotFound,
			"user not found",
		)
	}

	log.Println(user.Password, body.Password)

	// Verifikasi password yang dimasukkan
	if !middlewares.VerifyPassword(user.Password, body.Password) {
		return nil, fault.Custom(
			http.StatusUnauthorized,
			fault.ErrUnauthorized,
			"Invalid password, does not match",
		)
	}

	// Bersihkan password sebelum dikembalikan
	user.Password = ""

	// Generate token JWT
	accessToken, payload, err := jwt.CreateAccessToken(user.Name, user.Email, string(user.Id.String()))
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := jwt.CreateRefreshToken(user.Name, user.Email, string(user.Id.String()))
	if err != nil {
		return nil, err
	}

	// Bentuk response
	res := &model.LoginResponse{
		UserData:              *user,
		AccessToken:           *accessToken,
		AccessTokenExpiresAt:  &payload.ExpiresAt.Time,
		RefreshToken:          *refreshToken,
		RefreshTokenExpiresAt: &refreshPayload.ExpiresAt.Time,
	}

	// Simpan response login ke Redis untuk caching
	jsonValue, err := json.Marshal(res)
	if err != nil {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to marshal login response to JSON: %v", err),
		)
	}

	cacheKey := fmt.Sprintf("login:%s", body.Email)
	err = cache.Set(ctx, u.redis, cacheKey, string(jsonValue), 10*time.Minute)
	if err != nil {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to set Redis key '%s': %v", cacheKey, err),
		)
	}

	return res, nil
}
