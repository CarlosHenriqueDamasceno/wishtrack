package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Find(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	IsEmailAlreadyTaken(ctx context.Context, email Email) (bool, error)
}

type InMemoryRepository struct {
	users map[string]*User
}

func NewInMemoryRepository() Repository {
	return &InMemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryRepository) Find(ctx context.Context, id uuid.UUID) (*User, error) {
	user := r.users[id.String()]
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryRepository) Create(ctx context.Context, user *User) error {
	r.users[user.ID.String()] = user
	return nil
}

func (r *InMemoryRepository) IsEmailAlreadyTaken(ctx context.Context, email Email) (bool, error) {
	for _, user := range r.users {
		if user.Email == email {
			return true, nil
		}
	}

	return false, nil
}

func (r *InMemoryRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range r.users {
		if string(user.Email) == email {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}
