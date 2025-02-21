package server

import (
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/CarlosHenriqueDamasceno/wishtrack/etc/doc"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Api struct {
	router      *http.ServeMux
	userService user.Service
	config      *Config
	logger      *slog.Logger
}

// Creates and starts a new Api.
// This method registers all handlers in the Api
func NewApi(router *http.ServeMux, config *Config, logger *slog.Logger, userService user.Service) *Api {
	s := &Api{
		router:      router,
		userService: userService,
		config:      config,
		logger:      logger,
	}
	s.setupRoutes()
	return s
}

func (a *Api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a *Api) setupRoutes() {
	docsUrl := fmt.Sprintf("%s/swagger/doc.json", a.config.Address)
	a.router.Handle("GET /swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

	a.router.HandleFunc("POST /api/v1/register", a.handleRegister)
	a.router.HandleFunc("POST /api/v1/login", a.handleLogin)
}
