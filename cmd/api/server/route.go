package server

import (
	"fmt"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (a *Api) setupRoutes() {
	docsUrl := fmt.Sprintf("%s/swagger/doc.json", a.config.Address)
	a.router.Handle("GET /swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

	a.router.HandleFunc("POST /api/v1/users/register", a.handleRegister)
	a.router.HandleFunc("POST /api/v1/users/login", a.handleLogin)

	a.router.HandleFunc("POST /api/v1/contents/write-down", a.AuthTokenMiddleware(a.handleWriteDown))
	a.router.HandleFunc("GET /api/v1/contents/feed", a.AuthTokenMiddleware(a.handleFeed))
	a.router.HandleFunc("PUT /api/v1/contents/{id}", a.AuthTokenMiddleware(a.handleContentEdit))
	a.router.HandleFunc("GET /api/v1/contents/{id}", a.AuthTokenMiddleware(a.handleFindContent))
	a.router.HandleFunc("POST /api/v1/contents/{id}/rate", a.AuthTokenMiddleware(a.handleRateContent))
}
