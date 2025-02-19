package server

import (
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
)

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
func (api *Api) handleRegister(w http.ResponseWriter, req *http.Request) {
	input, err := utils.ReceiveJSON[*user.RegisterInput](req)
	if err != nil {
		api.logger.Info("Failed to decode JSON", "error", err)
		utils.RespondError(err, http.StatusBadRequest, w)
		return
	}

	out, err := api.userService.Register(req.Context(), input)
	if err != nil {
		switch err.(type) {
		case validation.ErrorCollection:
			api.logger.Info("Failed to decode JSON", "error", err)
			utils.RespondJSON(err, http.StatusUnprocessableEntity, w)
			return
		default:
			api.logger.Error("Failed to register user", "error", err)
			utils.RespondError(err, http.StatusInternalServerError, w)
			return
		}
	}

	utils.RespondJSON(out, http.StatusCreated, w)
}
