package handlers

import (
	"net/http"
	"truth/rest/server"
)

type StartHandler struct {
	server *server.Server
}

func NewStartHandler(s *server.Server) *StartHandler {
	return &StartHandler{s}
}

func (s *StartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.Run()

	w.WriteHeader(http.StatusOK)
}
