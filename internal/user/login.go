package user

import (
	"context"
	"errors"
)

var ErrIncorrectCredentials = errors.New("the provided e-mail or password are incorrect")

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

func (service *service) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, err := service.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			return nil, ErrIncorrectCredentials
		default:
			return nil, err
		}
	}

	if !user.VerifyPassword(input.Password) {
		return nil, ErrIncorrectCredentials
	}

	token, err := service.authenticator.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token: token,
	}, nil
}

func (service *service) ValidateToken(ctx context.Context, token string) (*User, error) {
	id, err := service.authenticator.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	user, err := service.repository.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
