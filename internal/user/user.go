package user

import (
	"regexp"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Email string

func (e Email) validate() *validation.ValidationError {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !pattern.MatchString(string(e)) {
		return &validation.ValidationError{
			Field:   "email",
			Message: "field \"email\" must be a valid e-mail address",
		}
	}
	return nil
}

type password struct {
	value string
	plain string
}

func newPassword(plainText string) (password, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return password{
		value: string(bytes),
		plain: plainText,
	}, err
}

func (p *password) validate() *validation.ValidationError {
	if len(p.plain) < 8 {
		return &validation.ValidationError{
			Field:   "password",
			Message: "field \"password\" must be at least 8 characters long",
		}
	}
	return nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password.value), []byte(password))
	return err == nil
}

type User struct {
	ID        uuid.UUID
	Name      string
	Email     Email
	password  password
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
		password: hashedPassword,
	}, nil
}
