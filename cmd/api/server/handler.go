package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
	"github.com/google/uuid"
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

// Feed godoc
//
//	@Summary	Get suggestions for today
//	@Tags		content
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	content.FeedOutput	"Suggestions"
//	@Failure	401	{string}	string				"Unauthorized"
//	@Router		/contents/feed [get]
//	@Security	ApiKeyAuth
func (api *Api) handleFeed(w http.ResponseWriter, r *http.Request) {
	user, err := api.GetLoggedUser(w, r)
	if err != nil {
		api.handleError(w, r, "invalid token", err)
		return
	}

	out, err := api.contentService.Feed(r.Context(), user.ID)
	if err != nil {
		api.handleError(w, r, "error getting feed", err)
		return
	}

	utils.RespondJSON(out, http.StatusOK, w)
}

// Write Down godoc
//
//	@Summary	Edits a content
//	@Tags		content
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		content.EditContentInput	true	"Content info"
//	@Param		id		path		string						true	"Content ID"
//	@Success	201		{object}	content.EditContentOutput	"Content Details"
//	@Failure	422		{object}	validation.ErrorCollection	"Validation error"
//	@Router		/contents/{id} [put]
//	@Security	ApiKeyAuth
func (api *Api) handleContentEdit(w http.ResponseWriter, r *http.Request) {
	input, err := utils.ReceiveJSON[*content.EditContentInput](r)
	if err != nil {
		api.handleError(w, r, "failed to decode edit content input", err)
		return
	}

	user, err := api.GetLoggedUser(w, r)
	if err != nil {
		api.handleError(w, r, "invalid token", err)
		return
	}
	input.UserID = user.ID

	input.ID, err = uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.handleError(w, r, "invalid uuid for editing content", utils.NewParsingError(err))
	}

	out, err := api.contentService.Edit(r.Context(), input)
	if err != nil {
		api.handleError(w, r, "error editing content", err)
		return
	}

	utils.RespondJSON(out, http.StatusOK, w)
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
		case errors.Is(err, content.ErrContentNotFound):
			utils.RespondError(err, http.StatusNotFound, w)
		case errors.Is(err, utils.ErrParsingError):
			api.logger.Info(logMessage, "request", requestLog, "error", err.Error())
			utils.RespondError(errors.Unwrap(err), http.StatusBadRequest, w)
		default:
			api.logger.Error(logMessage, "request", requestLog, "error", err.Error())
			utils.RespondError(err, http.StatusInternalServerError, w)
		}
	}
}
