package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/tasks"
)

const (
	NewTaskURL = "/tasks"
)

func Task(s *mux.Router) *mux.Router {
	url := fmt.Sprintf("%s%s", NewTaskURL, "/task")
	fmt.Println(url)
	s.HandleFunc(url, tasks.GetNewTaskWizardRequest).Methods(http.MethodPost, http.MethodOptions)
	url = fmt.Sprintf("%s%s", NewTaskURL, userIDURL)
	fmt.Println(url)
	s.HandleFunc(url, tasks.GetOnGoingNewTask).Methods(http.MethodGet, http.MethodOptions)
	return s
}
