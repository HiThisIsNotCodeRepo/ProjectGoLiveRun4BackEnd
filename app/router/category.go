package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/category"
)

func Category(s *mux.Router) *mux.Router {
	s.HandleFunc("/categories", category.Categories).Methods(http.MethodGet, http.MethodOptions)
	return s
}
