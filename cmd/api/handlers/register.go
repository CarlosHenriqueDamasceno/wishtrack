package handlers

import (
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
)

type RegisterHandler struct {
	service user.Service
}

func NewRegisterHandler(service user.Service) *RegisterHandler {
	return &RegisterHandler{
		service: service,
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
	input, err := utils.ReceiveJSON[*user.RegisterInput](req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	out, err := h.service.Register(req.Context(), input)
	if err != nil {
		utils.RespondJSON(err, http.StatusUnprocessableEntity, w)
		return
	}

	utils.RespondJSON(out, http.StatusCreated, w)
}
