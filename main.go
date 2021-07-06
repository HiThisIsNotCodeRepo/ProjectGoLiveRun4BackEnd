package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"paotui.sg/app/middleware"
	"paotui.sg/app/router"
)

func main() {
	defer db.Db.Close()
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s = router.Task(s)
	s = router.Auth(s)
	s = router.MyInfo(s)
	s = router.Category(s)
	s = router.TaskBid(s)
	s = router.User(s)
	s.Use(middleware.AuthMiddleware)
	s.Use(mux.CORSMethodMiddleware(s))
	log.Fatalln(http.ListenAndServe(":5000", s))
}
