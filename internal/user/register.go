package user

import (
	"context"
	"log"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/validation"
	"github.com/google/uuid"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *service) Register(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	user, err := NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	err = s.validate(ctx, user)
	if err != nil {
		return nil, err
	}

	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return outputFromUser(user), nil
}

func (s *service) validate(ctx context.Context, user *User) error {
	var errors validation.ErrorCollection
	err := user.Email.Validate()
	if err != nil {
		errors.Append(err)
	}

	isTaken, databaseErr := s.repository.IsEmailAlreadyTaken(ctx, user.Email)
	if databaseErr != nil {
		return databaseErr
	}

	log.Println(isTaken)

	if isTaken {
		errors.WithMessage("email", "e-mail already in use")
	}

	if errors.HasError() {
		return errors
	}

	return nil
}

func outputFromUser(user *User) *RegisterOutput {
	return &RegisterOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     string(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
