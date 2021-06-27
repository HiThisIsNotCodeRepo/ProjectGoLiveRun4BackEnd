package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/handler/new_task"
)

const (
	NewTaskURL = "/tasks"
)

func NewTaskRouter(s *mux.Router) *mux.Router {
	url := fmt.Sprintf("%s%s", NewTaskURL, "/task")
	fmt.Println(url)
	s.HandleFunc(url, new_task.GetNewTaskWizardRequest).Methods(http.MethodPost, http.MethodOptions)
	return s
}
