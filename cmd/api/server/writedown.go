package server

import (
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
)

// Write Down godoc
//
//	@Summary	Writes Down a content
//	@Tags		content
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		content.WriteDownInput		true	"Content info"
//	@Success	201		{object}	content.WriteDownOutput		"Content Details"
//	@Failure	422		{object}	validation.ErrorCollection	"Validation error"
//	@Router		/write-down [post]
//	@Security	ApiKeyAuth
func (api *Api) handleWriteDown(w http.ResponseWriter, r *http.Request) {
	input, err := utils.ReceiveJSON[*content.WriteDownInput](r)
	if err != nil {
		api.logger.Info("failed to decode JSON", "error", err)
		utils.RespondError(err, http.StatusBadRequest, w)
		return
	}

	user, err := api.GetLoggedUser(w, r)
	if err != nil {
		utils.RespondError(err, http.StatusUnauthorized, w)
		return
	}

	input.UserID = user.ID

	out, err := api.contentService.WriteDown(r.Context(), input)
	if err != nil {
		switch err.(type) {
		case validation.ErrorCollection:
			api.logger.Info("failed to decode JSON", "error", err)
			utils.RespondJSON(err, http.StatusUnprocessableEntity, w)
			return
		default:
			api.logger.Error("failed to write down content", "error", err)
			utils.RespondError(err, http.StatusInternalServerError, w)
			return
		}
	}

	utils.RespondJSON(out, http.StatusCreated, w)
}
