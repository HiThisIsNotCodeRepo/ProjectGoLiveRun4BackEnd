package task_update

import "net/http"

func TaskUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	option := r.URL.Query().Get("option")
	if option == "confirm-task-deliver" {
		ConfirmTaskDeliver(w, r)
	} else if option == "update-expected-rate" {
		UpdateTaskExpectedRate(w, r)
	} else if option == "delete" {
		DeleteTask(w, r)
	}
}
