package user

import (
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	Find(id uuid.UUID) (*User, error)
	Create(user *User) error
}

type InMemoryRepository struct {
	users map[string]*User
}

func NewInMemoryRepository() Repository {
	return &InMemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryRepository) Find(id uuid.UUID) (*User, error) {
	user := r.users[id.String()]
	if user == nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}

func (r *InMemoryRepository) Create(user *User) error {
	r.users[user.ID.String()] = user
	return nil
}
