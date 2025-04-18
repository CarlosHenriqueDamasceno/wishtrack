package server

import (
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (a *Api) setupRoutes() {

	docsUrl := fmt.Sprintf("%s/swagger/doc.json", a.config.Address)
	a.router.Handle("GET /swagger/", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

	a.router.HandleFunc("POST /api/v1/users/register", a.CorsMiddleware(a.handleRegister))
	a.router.HandleFunc("POST /api/v1/users/login", a.CorsMiddleware(a.handleLogin))

	a.router.HandleFunc("GET /api/v1/contents", a.CorsMiddleware(a.AuthTokenMiddleware(a.handleListContents)))
	a.router.HandleFunc("POST /api/v1/contents/write-down", a.CorsMiddleware(a.AuthTokenMiddleware(a.handleWriteDown)))
	a.router.HandleFunc("GET /api/v1/contents/suggestions", a.CorsMiddleware(a.AuthTokenMiddleware(a.handleContentSuggestions)))
	a.router.HandleFunc("PUT /api/v1/contents/{id}", a.CorsMiddleware(a.AuthTokenMiddleware(a.handleContentEdit)))
	a.router.HandleFunc("GET /api/v1/contents/{id}", a.CorsMiddleware(a.AuthTokenMiddleware(a.handleFindContent)))
	a.router.HandleFunc("DELETE /api/v1/contents/{id}", a.CorsMiddleware(a.AuthTokenMiddleware(a.handleDeleteContent)))
	a.router.HandleFunc("POST /api/v1/contents/{id}/rate", a.CorsMiddleware(a.AuthTokenMiddleware(a.handleRateContent)))

	a.router.HandleFunc("GET /api/v1/suggestions", a.CorsMiddleware(a.AuthTokenMiddleware(a.handleSuggestions)))

	a.router.HandleFunc("OPTIONS /", a.CorsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
}
