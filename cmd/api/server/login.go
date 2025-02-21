package server

import (
	"errors"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
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
		switch {
		case errors.Is(err, user.ErrIncorrectCredentials):
			utils.RespondError(err, http.StatusUnauthorized, w)
		default:
			api.logger.Error("failed login user", "error", err)
			utils.RespondError(err, http.StatusInternalServerError, w)
			return
		}
	}

	utils.RespondJSON(out, http.StatusOK, w)
}
