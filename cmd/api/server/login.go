package server

import (
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
)

// Register godoc
//
//	@Summary	Login
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		user.RegisterInput	true	"User credentials"
//	@Success	200		{string}	"token"
//	@Failure	422		{object}	validation.ErrorCollection	"Validation Errors"
//	@Router		/register [post]
func (api *Api) handleLogin(w http.ResponseWriter, r *http.Request) {
	input, err := utils.ReceiveJSON[*user.LoginInput](r)
	if err != nil {
		api.logger.Info("failed to decode JSON", "error", err)
		utils.RespondError(err, http.StatusBadRequest, w)
		return
	}

	out, err := api.userService.Login(r.Context(), input)
	if err != nil {
		switch err.(type) {
		case validation.ErrorCollection:
			api.logger.Info("failed to decode JSON", "error", err)
			utils.RespondJSON(err, http.StatusUnprocessableEntity, w)
			return
		default:
			api.logger.Error("failed login user", "error", err)
			utils.RespondError(err, http.StatusInternalServerError, w)
			return
		}
	}

	utils.RespondJSON(out, http.StatusOK, w)
}
