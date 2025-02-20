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
		return nil, err
	}

	if !user.VerifyPassword(input.Password) {
		return nil, ErrIncorrectCredentials
	}

	//TODO: Generate token

	return &LoginOutput{
		Token: "asdasdadasda",
	}, nil
}
