package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
)

// Register godoc
//
//	@Summary	Creates a new user
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		user.RegisterInput	true	"User information for registration"
//	@Success	201		{object}	user.RegisterOutput
//	@Failure	422		{object}	validation.ErrorCollection	"Validation Errors"
//	@Router		/users/register [post]
func (api *Api) handleRegister(w http.ResponseWriter, r *http.Request) {
	input, err := utils.ReceiveJSON[*user.RegisterInput](r)
	if err != nil {
		api.handleError(w, r, "failed to decode register input", err)
		return
	}

	out, err := api.userService.Register(r.Context(), input)
	if err != nil {
		api.handleError(w, r, "error registering user", err)
		return
	}

	utils.RespondJSON(out, http.StatusCreated, w)
}

// Register godoc
//
//	@Summary	Login
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		user.LoginInput		true	"User credentials"
//	@Success	200		{object}	user.LoginOutput	"token"
//	@Failure	401		{string}	string				"Invalid credentials"
//	@Router		/users/login [post]
func (api *Api) handleLogin(w http.ResponseWriter, r *http.Request) {
	input, err := utils.ReceiveJSON[*user.LoginInput](r)
	if err != nil {
		api.handleError(w, r, "failed to decode login input", err)
		return
	}

	out, err := api.userService.Login(r.Context(), input)
	if err != nil {
		api.handleError(w, r, "login error", err)
		return
	}

	utils.RespondJSON(out, http.StatusOK, w)
}

// Write Down godoc
//
//	@Summary	Writes Down a content
//	@Tags		content
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		content.WriteDownInput		true	"Content info"
//	@Success	201		{object}	content.WriteDownOutput		"Content Details"
//	@Failure	422		{object}	validation.ErrorCollection	"Validation error"
//	@Router		/contents/write-down [post]
//	@Security	ApiKeyAuth
func (api *Api) handleWriteDown(w http.ResponseWriter, r *http.Request) {
	input, err := utils.ReceiveJSON[*content.WriteDownInput](r)
	if err != nil {
		api.handleError(w, r, "failed to decode write down input", err)
		return
	}

	user, err := api.GetLoggedUser(w, r)
	if err != nil {
		api.handleError(w, r, "invalid token", err)
		return
	}

	input.UserID = user.ID

	out, err := api.contentService.WriteDown(r.Context(), input)
	if err != nil {
		api.handleError(w, r, "error writing down content", err)
		return
	}

	utils.RespondJSON(out, http.StatusCreated, w)
}

func (api *Api) handleError(w http.ResponseWriter, r *http.Request, logMessage string, err error) {
	if _, ok := err.(validation.ErrorCollection); ok {
		api.logger.Info("validation error", "error", err)
		utils.RespondJSON(err, http.StatusUnprocessableEntity, w)
		return
	}

	if err != nil {
		requestLog := fmt.Sprintf("[Path]: %s [Method]: %s", r.URL.Path, r.Method)
		switch {
		case errors.Is(err, user.ErrIncorrectCredentials):
			utils.RespondError(err, http.StatusUnauthorized, w)
		case errors.Is(err, utils.ErrParsingError):
			api.logger.Info(logMessage, "request", requestLog, "error", err.Error())
			utils.RespondError(errors.Unwrap(err), http.StatusBadRequest, w)
		default:
			api.logger.Error(logMessage, "request", requestLog, "error", err.Error())
			utils.RespondError(err, http.StatusInternalServerError, w)
		}
	}

}
