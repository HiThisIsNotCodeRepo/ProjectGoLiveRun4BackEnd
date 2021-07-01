package task_update

import (
	"net/http"
	"strings"
)

func TaskUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	option := strings.TrimSpace(r.URL.Query().Get("option"))
	if option == "confirm-task-deliver" {
		ConfirmTaskDeliver(w, r)
	} else if option == "update-expected-rate" {
		UpdateTaskExpectedRate(w, r)
	} else if option == "delete" {
		DeleteTask(w, r)
	}
}
