package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/my_info/my_info_earning"
	"paotui.sg/app/handler/my_info/my_info_spending"
)

const (
	userIDURL = "/{userID}"
)

func MyInfo(s *mux.Router) *mux.Router {
	s.HandleFunc("/earning/{userID}", my_info_earning.Earning).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/{userID}", my_info_spending.Spending).Methods(http.MethodGet, http.MethodOptions)
	return s
}
