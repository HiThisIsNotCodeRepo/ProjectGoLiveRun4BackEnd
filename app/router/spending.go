package router

import (
	"github.com/gorilla/mux"
	"net/http"
	my_info "paotui.sg/app/handler/my-info"
)

func SpendRouter(s *mux.Router) *mux.Router {
	s.HandleFunc("/spending/buy-necessity/yesterday/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/buy-necessity/two-days-ago/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/buy-necessity/three-days-ago/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/food-delivery/yesterday/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/food-delivery/two-days-ago/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/food-delivery/three-days-ago/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/send-document/yesterday/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/send-document/two-days-ago/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/send-document/three-days-ago/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/other/yesterday/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/other/two-days-ago/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/spending/other/three-days-ago/{userID}", my_info.GetSpending).Methods(http.MethodGet, http.MethodOptions)
	return s
}
