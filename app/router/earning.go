package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/my_info"
)

const (
	EarningURL = "/earning"
)

func EarningTaskRouter(s *mux.Router) *mux.Router {
	url := fmt.Sprintf("%s%s%s", EarningURL, "/tasks", userIDURL)
	fmt.Println(url)
	s.HandleFunc(url, my_info.GetEarningDataSource).Methods(http.MethodGet, http.MethodOptions)
	url = fmt.Sprintf("%s%s%s", EarningURL, "/past-days", userIDURL)
	fmt.Println(url)
	s.HandleFunc(url, my_info.GetEarningCard).Methods(http.MethodGet, http.MethodOptions)
	url = fmt.Sprintf("%s%s%s%s", EarningURL, "/radar","/last-week", userIDURL)
	fmt.Println(url)
	s.HandleFunc(url, my_info.GetEarningRadar).Methods(http.MethodGet, http.MethodOptions)
	url = fmt.Sprintf("%s%s%s%s", EarningURL,"/radar", "/this-week", userIDURL)
	fmt.Println(url)
	s.HandleFunc(url, my_info.GetEarningRadar).Methods(http.MethodGet, http.MethodOptions)
	return s
}
