package usecases

import (
	"time"
	"user-svc/helpers/jwt"
	"user-svc/middlewares"
	"user-svc/model"
	repository "user-svc/repository/user"

	"github.com/google/uuid"
)

type userUsecase struct {
	user repository.UserRepository
}

func NewUserUsecase(repository repository.UserRepository) *userUsecase {
	return &userUsecase{
		user: repository,
	}
}

type UserUsecases interface {
	UserRegister(body model.RegisterUser) (*model.LoginResponse, error)
}

func (u *userUsecase) UserRegister(body model.RegisterUser) (*model.LoginResponse, error) {
	exist, err := u.user.UserExistsByName(body.Name)
	if err != nil {
		return nil, err
	}

	var user *model.User
	var userId uuid.UUID

	if !exist {
		body.Password = middlewares.GenerateHashed(body.Password)

		createdId, err := u.user.InsertUser(body)
		if err != nil {
			return nil, err
		}
		userId = *createdId

		user, err = u.user.GetUserDetail(model.GetUserDetailRequest{
			UserId: userId,
		})
		if err != nil {
			return nil, err
		}
	} else {
		user, err = u.user.GetUserDetail(model.GetUserDetailRequest{
			Email: body.Email,
		})
		if err != nil {
			return nil, err
		}
		userId = user.Id
	}

	var (
		tokenExpiry        = 30 * time.Minute
		refreshTokenExpiry = 72 * time.Hour
	)

	accessToken, payload, err := jwt.CreateAccessToken(user.Name, user.Email, userId.String(), tokenExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := jwt.CreateRefreshToken(user.Name, user.Email, userId.String(), refreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		UserData:              *user,
		AccessToken:           *accessToken,
		AccessTokenExpiresAt:  &payload.ExpiresAt.Time,
		RefreshToken:          *refreshToken,
		RefreshTokenExpiresAt: &refreshPayload.ExpiresAt.Time,
	}, nil
}
