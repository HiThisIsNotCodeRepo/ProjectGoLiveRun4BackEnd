package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/auth"
)

func Auth(s *mux.Router) *mux.Router {
	s.HandleFunc("/auth", auth.Auth).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
	return s
}
