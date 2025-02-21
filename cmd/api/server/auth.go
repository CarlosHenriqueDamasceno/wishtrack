package server

import (
	"errors"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
)

var ErrLoggedUserNotFound = errors.New("could not get logged user id")

func (api *Api) GetLoggedUser(w http.ResponseWriter, r *http.Request) (*user.User, error) {
	contextValue := r.Context().Value(userCtx)
	if user, ok := contextValue.(*user.User); ok {
		return user, nil
	}

	api.logger.Warn("error getting logged user", "value", contextValue)
	return &user.User{}, ErrLoggedUserNotFound
}
