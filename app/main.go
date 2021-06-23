package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/router"
)

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s = router.SpendRouter(s)
	s.Use(mux.CORSMethodMiddleware(s))
	log.Fatalln(http.ListenAndServe(":5000", s))
}
