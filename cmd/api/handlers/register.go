package handlers

import (
	"log/slog"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/validation"
)

type RegisterHandler struct {
	service user.Service
	logger  *slog.Logger
}

func NewRegisterHandler(service user.Service, logger *slog.Logger) *RegisterHandler {
	return &RegisterHandler{
		service: service,
		logger:  logger,
	}
}

// Register godoc
//
//	@Summary	Creates a new user
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		registerRequest	body		user.RegisterInput	true	"User information for registration"
//	@Success	201				{object}	user.RegisterOutput
//	@Failure	422				{object}	validation.ErrorCollection	"Validation Errors"
//	@Router		/register [post]
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	input, err := utils.ReceiveJSON[*user.RegisterInput](req)
	if err != nil {
		h.logger.Info("Failed to decode JSON", "error", err)
		utils.RespondError(err, http.StatusBadRequest, w)
		return
	}

	out, err := h.service.Register(req.Context(), input)
	if err != nil {
		switch err.(type) {
		case validation.ErrorCollection:
			h.logger.Info("Failed to decode JSON", "error", err)
			utils.RespondJSON(err, http.StatusUnprocessableEntity, w)
			return
		default:
			h.logger.Error("Failed to register user", "error", err)
			utils.RespondError(err, http.StatusInternalServerError, w)
			return
		}
	}

	utils.RespondJSON(out, http.StatusCreated, w)
}
