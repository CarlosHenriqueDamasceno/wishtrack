package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/utils"
)

const userCtx = "user"

func (api *Api) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(fmt.Errorf("authorization header is missing"), http.StatusUnauthorized, w)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondError(fmt.Errorf("authorization header is missing"), http.StatusUnauthorized, w)
			return
		}

		ctx := r.Context()

		token := parts[1]
		user, err := api.userService.ValidateToken(ctx, token)
		if err != nil {
			utils.RespondError(err, http.StatusUnauthorized, w)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
