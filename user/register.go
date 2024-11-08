package user

import (
	"encoding/json"
	"net/http"
	"time"

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

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func InputFromRequest(req *http.Request) (*RegisterInput, error) {
	input := &RegisterInput{}
	err := json.NewDecoder(req.Body).Decode(input)
	return input, err
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
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

// Register godoc
// @Summary      Create's a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param registerRequest body user.RegisterInput true "User information for registration"
// @Success      201  {object}   user.RegisterOutput
// @Router       /register [post]
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	input, err := InputFromRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	out, err := h.register(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *RegisterHandler) register(input *RegisterInput) (*RegisterOutput, error) {
	usr, err := NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	h.repository.Create(usr)

	return outputFromUser(usr), nil
}
