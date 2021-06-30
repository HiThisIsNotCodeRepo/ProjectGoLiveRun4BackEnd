package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/user"
)

func User(s *mux.Router) *mux.Router {
	s.HandleFunc("/users/user", user.Register).Methods(http.MethodPost, http.MethodOptions)
	return s
}
