package task_update

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"strings"
)

type DeleteTaskResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	var deleteTaskResponse DeleteTaskResponse
	var err error
	fmt.Printf("DeleteTask->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	taskId := mux.Vars(r)["taskID"]
	if strings.TrimSpace(taskId) == "" {
		deleteTaskResponse.Status = "error"
		deleteTaskResponse.Msg = "No TaskId"
		goto Label1
	}
	_, err = db.Db.Exec("DELETE FROM task_bid WHERE task_id =?", taskId)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	_, err = db.Db.Exec("DELETE FROM task WHERE task_id = ? ", taskId)
	if err != nil {
		log.Println(err)
	}
	deleteTaskResponse.Status = "success"
	deleteTaskResponse.Msg = "delete task success"

Label0:
	if deleteTaskResponse.Status != "success" {
		deleteTaskResponse.Status = "error"
		deleteTaskResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(deleteTaskResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
