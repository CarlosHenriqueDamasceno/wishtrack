package server

import (
	"log/slog"
	"net/http"

	_ "github.com/CarlosHenriqueDamasceno/wishtrack/etc/doc"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/suggestion"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
)

type Api struct {
	router            *http.ServeMux
	userService       user.Service
	contentService    content.Service
	suggestionService suggestion.Service
	config            *Config
	logger            *slog.Logger
}

// Creates and starts a new Api.
// This method registers all handlers in the Api
func NewApi(
	router *http.ServeMux,
	config *Config,
	logger *slog.Logger,
	userService user.Service,
	contentService content.Service,
	suggestionService suggestion.Service,
) *Api {
	s := &Api{
		router:            router,
		userService:       userService,
		contentService:    contentService,
		suggestionService: suggestionService,
		config:            config,
		logger:            logger,
	}
	s.setupRoutes()
	return s
}

func (a *Api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}
