package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/bid"
)

func TaskBid(s *mux.Router) *mux.Router {
	s.HandleFunc("/tasks/bid", bid.TaskBid).Methods(http.MethodPost, http.MethodOptions)
	return s
}
