package server

import (
	"errors"
	"fmt"
	"net/http"

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
	input, err := ReceiveJSON[*user.RegisterInput](r)
	if err != nil {
		api.handleError(w, r, "failed to decode register input", err)
		return
	}

	out, err := api.userService.Register(r.Context(), input)
	if err != nil {
		api.handleError(w, r, "error registering user", err)
		return
	}

	RespondJSON(out, http.StatusCreated, w)
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
	input, err := ReceiveJSON[*user.LoginInput](r)
	if err != nil {
		api.handleError(w, r, "failed to decode login input", err)
		return
	}

	out, err := api.userService.Login(r.Context(), input)
	if err != nil {
		api.handleError(w, r, "login error", err)
		return
	}

	RespondJSON(out, http.StatusOK, w)
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
	input, err := ReceiveJSON[*content.WriteDownInput](r)
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

	RespondJSON(out, http.StatusCreated, w)
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

	RespondJSON(out, http.StatusOK, w)
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
	input, err := ReceiveJSON[*content.EditContentInput](r)
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
		api.handleError(w, r, "invalid uuid for editing content", NewParsingError(err))
	}

	out, err := api.contentService.Edit(r.Context(), input)
	if err != nil {
		api.handleError(w, r, "error editing content", err)
		return
	}

	RespondJSON(out, http.StatusOK, w)
}

// Rate content godoc
//
//	@Summary	Rate a content
//	@Tags		content
//	@Accept		json
//	@Produce	json
//	@Param		payload	body	content.RateContentInput	true	"Rate"
//	@Param		id		path	string						true	"Content ID"
//	@Success	204
//	@Failure	422	{object}	validation.ErrorCollection	"Validation error"
//	@Router		/contents/{id}/rate [post]
//	@Security	ApiKeyAuth
func (api *Api) handleRateContent(w http.ResponseWriter, r *http.Request) {
	input, err := ReceiveJSON[*content.RateContentInput](r)
	if err != nil {
		api.handleError(w, r, "failed to decode rate content input", err)
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
		api.handleError(w, r, "invalid uuid for rating content", NewParsingError(err))
	}

	err = api.contentService.Rate(r.Context(), input)
	if err != nil {
		api.handleError(w, r, "error rating content", err)
		return
	}

	RespondJSON(nil, http.StatusNoContent, w)
}

// Find content godoc
//
//	@Summary	finds a content
//	@Tags		content
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string						true	"Content ID"
//	@Success	200	{object}	content.FindContentOutput	"Content info"
//	@Failure	404	{string}	string						"Content not found"
//	@Router		/contents/{id} [get]
//	@Security	ApiKeyAuth
func (api *Api) handleFindContent(w http.ResponseWriter, r *http.Request) {
	user, err := api.GetLoggedUser(w, r)
	if err != nil {
		api.handleError(w, r, "invalid token", err)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.handleError(w, r, "invalid uuid for finding a content", NewParsingError(err))
	}

	out, err := api.contentService.Find(r.Context(), id, user.ID)
	if err != nil {
		api.handleError(w, r, "error finding a content", err)
		return
	}

	RespondJSON(out, http.StatusOK, w)
}

// Delete a content godoc
//
//	@Summary	deletes a content
//	@Tags		content
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Content ID"
//	@Success	204	{object}	string
//	@Failure	404	{string}	string	"Content not found"
//	@Router		/contents/{id} [delete]
//	@Security	ApiKeyAuth
func (api *Api) handleDeleteContent(w http.ResponseWriter, r *http.Request) {
	user, err := api.GetLoggedUser(w, r)
	if err != nil {
		api.handleError(w, r, "invalid token", err)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.handleError(w, r, "invalid uuid for deleting a content", NewParsingError(err))
	}

	err = api.contentService.Delete(r.Context(), id, user.ID)
	if err != nil {
		api.handleError(w, r, "error deleting a content", err)
		return
	}

	RespondJSON(nil, http.StatusNoContent, w)
}

func (api *Api) handleError(w http.ResponseWriter, r *http.Request, logMessage string, err error) {
	if _, ok := err.(validation.ErrorCollection); ok {
		api.logger.Info("validation error", "error", err)
		RespondJSON(err, http.StatusUnprocessableEntity, w)
		return
	}

	if err != nil {
		requestLog := fmt.Sprintf("[Path]: %s [Method]: %s", r.URL.Path, r.Method)
		switch {
		case errors.Is(err, user.ErrIncorrectCredentials):
			RespondError(err, http.StatusUnauthorized, w)
		case errors.Is(err, content.ErrContentNotFound):
			RespondError(err, http.StatusNotFound, w)
		case errors.Is(err, ErrParsingError):
			api.logger.Info(logMessage, "request", requestLog, "error", err.Error())
			RespondError(errors.Unwrap(err), http.StatusBadRequest, w)
		default:
			api.logger.Error(logMessage, "request", requestLog, "error", err.Error())
			RespondError(err, http.StatusInternalServerError, w)
		}
	}
}
