package user

import (
	"net/http"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/validation"
	"github.com/google/uuid"
)

type RegisterHandler struct {
	repository Repository
}

func NewRegisterHandler(repository Repository) *RegisterHandler {
	return &RegisterHandler{
		repository: repository,
	}
}

// Register godoc
// @Summary      Creates a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        registerRequest body user.RegisterInput true "User information for registration"
// @Success      201  {object}   user.RegisterOutput
// @Failure      422  {object}   map[string][]string "Validation Errors"
// @Router       /register [post]
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	input, err := utils.ReceiveJSON[*RegisterInput](req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	out, err := h.register(input)
	if err != nil {
		utils.RespondJSON(err, http.StatusUnprocessableEntity, w)
		return
	}

	utils.RespondJSON(out, http.StatusCreated, w)
}

type RegisterInput struct {
	Name     string
	Email    Email
	Password string
}

func (i *RegisterInput) validate() error {
	var errors validation.ErrorCollection
	err := i.Email.Validate()
	if err != nil {
		errors.Append(err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}

type RegisterOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func outputFromUser(user *User) *RegisterOutput {
	return &RegisterOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     string(user.Email),
		CreatedAt: user.CreatedAt,
	}
}

func (h *RegisterHandler) register(input *RegisterInput) (*RegisterOutput, error) {
	err := input.validate()
	if err != nil {
		return nil, err
	}

	usr, err := NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	h.repository.Create(usr)

	return outputFromUser(usr), nil
}
