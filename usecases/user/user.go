package usecases

import (
	"context"
	"encoding/json"
	"fmt"
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

type userUsecase struct {
	user  repository.UserRepository
	redis *redis.Client
}

func NewUserUsecase(repository repository.UserRepository, redis *redis.Client) *userUsecase {
	return &userUsecase{
		user:  repository,
		redis: redis,
	}
}

type UserUsecases interface {
	Register(body model.RegisterUser) (*model.LoginResponse, error)
}

func (u *userUsecase) Register(body model.RegisterUser) (*model.LoginResponse, error) {
	ctx := context.TODO()
	cacheKey := fmt.Sprintf("login:%s", body.Email)

	if cacheExist, err := cache.Exist(ctx, u.redis, cacheKey); err != nil {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to check Redis key existence for '%s': %v", cacheKey, err),
		)
	} else if cacheExist {
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

	exist, err := u.user.ExistsByName(body.Name)
	if err != nil {
		return nil, err
	}

	var user *model.User
	var userId uuid.UUID

	if !exist {
		body.Password = middlewares.GeneratePassword(body.Password)

		createdId, err := u.user.Insert(body)
		if err != nil {
			return nil, err
		}
		userId = *createdId

		user, err = u.user.Detail(model.GetUserDetailRequest{
			UserId: userId,
		})
		if err != nil {
			return nil, err
		}
	} else {
		user, err = u.user.Detail(model.GetUserDetailRequest{
			Email: body.Email,
		})
		if err != nil {
			return nil, err
		}
		userId = user.Id
	}

	accessToken, payload, err := jwt.CreateAccessToken(user.Name, user.Email, userId.String())
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := jwt.CreateRefreshToken(user.Name, user.Email, userId.String())
	if err != nil {
		return nil, err
	}

	res := &model.LoginResponse{
		UserData:              *user,
		AccessToken:           *accessToken,
		AccessTokenExpiresAt:  &payload.ExpiresAt.Time,
		RefreshToken:          *refreshToken,
		RefreshTokenExpiresAt: &refreshPayload.ExpiresAt.Time,
	}

	jsonValue, err := json.Marshal(res)
	if err != nil {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to marshal login response to JSON for key '%s': %v", cacheKey, err),
		)
	}

	if err := cache.Set(ctx, u.redis, cacheKey, string(jsonValue), 10*time.Minute); err != nil {
		return nil, err
	}

	return res, nil
}
