package server

import (
	"net/http"

	"github.com/qwuiemme/ellipsespace-server/config"
	"github.com/qwuiemme/ellipsespace-server/internal/handler"
)

type Server struct {
	core *http.Server
}

func (s *Server) Run(addr string) error {
	s.core = &http.Server{
		Addr:    addr,
		Handler: handler.InitHandler(),
	}

	return s.core.ListenAndServe()
}

func (s Server) MakeAddr() string {
	conf := config.New()

	return conf.AppIP + ":" + conf.AppPort
}
