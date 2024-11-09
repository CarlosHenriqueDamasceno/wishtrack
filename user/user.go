package user

import (
	"regexp"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Email string

func (e Email) Validate() *validation.ValidationError {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !pattern.MatchString(string(e)) {
		return &validation.ValidationError{
			Field:   "email",
			Message: "field \"email\" must be a valid e-mail address",
		}
	}
	return nil
}

type User struct {
	ID        uuid.UUID
	Name      string
	Email     Email
	password  string
	CreatedAt time.Time
}

func NewUser(name string, email Email, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     Email(email),
		password:  hashedPassword,
		CreatedAt: time.Now(),
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err == nil
}
