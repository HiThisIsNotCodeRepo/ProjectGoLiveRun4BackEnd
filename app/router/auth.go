package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/auth"
)

const (
	AuthURL        = "/auth"
	LoginURL       = "/login"
	TokenVerifyURL = "/token-verify"
)

func AuthRouter(s *mux.Router) *mux.Router {
	url := fmt.Sprintf("%s%s", AuthURL, LoginURL)
	fmt.Println(url)
	s.HandleFunc(url, auth.UserLogin).Methods(http.MethodPost, http.MethodOptions)
	url = fmt.Sprintf("%s%s%s", AuthURL, TokenVerifyURL, userIDURL)
	fmt.Println(url)
	s.HandleFunc(url, auth.VerifyToken).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
	return s
}
