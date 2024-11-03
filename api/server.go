package api

import (
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/register"
)

type ApiServer struct {
	router *http.ServeMux
}

// Creates and boots a new ApiServer.
// This method registers all handlers in the ApiServer
func NewApiServer(router *http.ServeMux) *ApiServer {
	s := &ApiServer{
		router: router,
	}
	s.start()
	return s
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *ApiServer) start() error {
	s.router.Handle("/api/v1/register", &register.RegisterHandler{})
	return nil
}
