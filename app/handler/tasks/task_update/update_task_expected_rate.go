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

type UpdateTaskExpectedRateResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}
type UpdateTaskExpectedRateRequest struct {
	Rate int `json:"rate"`
}

func UpdateTaskExpectedRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var updateTaskExpectedRateRequest UpdateTaskExpectedRateRequest
	var updateTaskExpectedRateResponse UpdateTaskExpectedRateResponse
	var err error
	fmt.Printf("request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	taskId := mux.Vars(r)["taskID"]
	decoder := json.NewDecoder(r.Body)
	if strings.TrimSpace(taskId) == "" {
		updateTaskExpectedRateResponse.Status = "error"
		updateTaskExpectedRateResponse.Msg = "No TaskId"
		goto Label1
	}
	err = decoder.Decode(&updateTaskExpectedRateRequest)
	if err != nil {
		log.Println(err)
		goto Label0
	}

	_, err = db.Db.Exec("UPDATE task SET task_owner_rate = ? WHERE task_id =?", updateTaskExpectedRateRequest.Rate, taskId)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	updateTaskExpectedRateResponse.Status = "success"
	updateTaskExpectedRateResponse.Msg = "owner expected rate updated success"

Label0:
	if updateTaskExpectedRateResponse.Status != "success" {
		updateTaskExpectedRateResponse.Status = "error"
		updateTaskExpectedRateResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(updateTaskExpectedRateResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
