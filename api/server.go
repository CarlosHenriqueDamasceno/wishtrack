package api

import (
	"net/http"

	_ "github.com/CarlosHenriqueDamasceno/wishtrack/etc/doc"
	"github.com/CarlosHenriqueDamasceno/wishtrack/user"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type ApiServer struct {
	router         *http.ServeMux
	userRepository user.Repository
}

// Creates and starts a new ApiServer.
// This method registers all handlers in the ApiServer
func NewApiServer(router *http.ServeMux, userRepository user.Repository) *ApiServer {
	s := &ApiServer{
		router:         router,
		userRepository: userRepository,
	}
	s.setupRoutes()
	return s
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *ApiServer) setupRoutes() {
	s.router.Handle("/swagger/*", httpSwagger.WrapHandler)
	s.router.Handle("POST /api/v1/register", user.NewRegisterHandler(s.userRepository))
}
