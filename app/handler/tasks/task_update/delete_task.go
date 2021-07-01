package task_update

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"paotui.sg/app/handler/error_util"
	"strings"
)

type DeleteTaskResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	defer error_util.ErrorHandle(w)
	var deleteTaskResponse DeleteTaskResponse
	var err error
	var tx *sql.Tx
	fmt.Printf("DeleteTask->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	taskId := mux.Vars(r)["taskID"]
	if strings.TrimSpace(taskId) == "" {
		deleteTaskResponse.Status = "error"
		deleteTaskResponse.Msg = "No TaskId"
		goto Label1
	}
	tx, err = db.Db.Begin()
	if err != nil {
		log.Println(err)
		goto Label0
	}
	_, err = tx.Exec("DELETE FROM task_bid WHERE task_id =?", taskId)
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		goto Label0
	}
	_, err = tx.Exec("DELETE FROM task WHERE task_id = ? ", taskId)
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		goto Label0
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		goto Label0
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
	}
}
