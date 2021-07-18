package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/avatar"
)

func Avatar(s *mux.Router) *mux.Router {
	s.HandleFunc("/avatar", avatar.NewAvatar).Methods(http.MethodPost, http.MethodOptions)
	return s
}
