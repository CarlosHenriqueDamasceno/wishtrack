package user

import (
	"regexp"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/validation"
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

type Password string

func newPassword(plainText string) (Password, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return Password(bytes), err
}

func (p *Password) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*p), []byte(password))
	return err == nil
}

type User struct {
	ID        uuid.UUID
	Name      string
	Email     Email
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name string, email string, password string) (*User, error) {
	hashedPassword, err := newPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New(),
		Name:     name,
		Email:    Email(email),
		Password: hashedPassword,
	}, nil
}
