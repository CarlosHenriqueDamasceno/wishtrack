package api

import (
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/handlers"
	_ "github.com/CarlosHenriqueDamasceno/wishtrack/etc/doc"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Api struct {
	router      *http.ServeMux
	userService user.Service
}

// Creates and starts a new Api.
// This method registers all handlers in the Api
func NewApi(router *http.ServeMux, userService user.Service) *Api {
	s := &Api{
		router:      router,
		userService: userService,
	}
	s.setupRoutes()
	return s
}

func (s *Api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Api) setupRoutes() {
	s.router.Handle("/swagger/*", httpSwagger.WrapHandler)
	s.router.Handle("POST /api/v1/register", handlers.NewRegisterHandler(s.userService))
}
