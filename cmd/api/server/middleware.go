package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type userCtxType string

const userCtx userCtxType = "user"

func (api *Api) AuthTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.logger.Info("authorization error no header")
			RespondError(fmt.Errorf("authorization header is missing"), http.StatusUnauthorized, w)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			api.logger.Info("authorization error invalid header", "header", authHeader)
			RespondError(fmt.Errorf("authorization header is missing"), http.StatusUnauthorized, w)
			return
		}

		ctx := r.Context()

		token := parts[1]
		user, err := api.userService.ValidateToken(ctx, token)
		if err != nil {
			api.logger.Warn("authorization error", "error", err)
			RespondError(err, http.StatusUnauthorized, w)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
