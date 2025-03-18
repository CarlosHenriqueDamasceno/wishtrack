package content

import (
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
	"github.com/google/uuid"
)

type Name string

func (n *Name) validate() *validation.ValidationError {
	if string(*n) == "" {
		return &validation.ValidationError{
			Field:   "name",
			Message: "field \"name\" is required",
		}
	}

	if len(*n) < 3 {
		return &validation.ValidationError{
			Field:   "name",
			Message: "field \"name\" must be at least 3 characters long",
		}
	}

	return nil
}

type Genres []string

func (g Genres) validate() validation.ErrorCollection {
	var errors validation.ErrorCollection
	for _, v := range g {
		if string(v) == "" {
			err := &validation.ValidationError{
				Field:   "genre",
				Message: "field \"genre\" is required",
			}

			errors.Append(err)
			continue
		}

		if len(v) < 3 {
			err := &validation.ValidationError{
				Field:   "genre",
				Message: "field \"genre\" must be at least 3 characters long",
			}

			errors.Append(err)
		}
	}
	return errors
}

type WishLevel int

func (wl *WishLevel) validate() *validation.ValidationError {
	if int(*wl) < 1 || int(*wl) > 5 {
		return &validation.ValidationError{
			Field:   "wish_level",
			Message: "field \"wish_level\" must be between 1 and 5",
		}
	}

	return nil
}

type Rate int

func (wl *Rate) validate() *validation.ValidationError {
	if int(*wl) < 1 || int(*wl) > 5 {
		return &validation.ValidationError{
			Field:   "rate",
			Message: "field \"rate\" must be between 1 and 5",
		}
	}

	return nil
}

type Content struct {
	ID        uuid.UUID
	Name      Name
	Category  string
	Genres    Genres
	Summary   string
	WishLevel WishLevel
	UserID    uuid.UUID
	Rate      *Rate
	Comment   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewContent(name, category string, genres Genres, summary string, wishLevel int, userID uuid.UUID) *Content {
	return &Content{
		ID:        uuid.New(),
		Name:      Name(name),
		Category:  category,
		Genres:    genres,
		Summary:   summary,
		WishLevel: WishLevel(wishLevel),
		UserID:    userID,
	}
}
