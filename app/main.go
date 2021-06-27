package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"paotui.sg/app/router"
)

func main() {
	defer db.Db.Close()
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s = router.SpendCardRouter(s)
	s = router.SpendSummaryRouter(s)
	s = router.SpendTaskRouter(s)
	s = router.EarningTaskRouter(s)
	s = router.NewTaskRouter(s)
	s = router.AuthRouter(s)
	s.Use(mux.CORSMethodMiddleware(s))
	log.Fatalln(http.ListenAndServe(":5000", s))
}
