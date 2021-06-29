package task_update

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"strings"
	"time"
)

type ConfirmTaskDeliveredResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}
type ConfirmTaskDeliveredRequest struct {
	DeliverRate int    `json:"deliverRate"`
	DeliverId   string `json:"deliverId"`
}

func ConfirmTaskDeliver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var confirmTaskDeliverRequest ConfirmTaskDeliveredRequest
	var confirmTaskDeliverResponse ConfirmTaskDeliveredResponse
	var err error
	fmt.Printf("ConfirmTaskDeliver->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	taskId := mux.Vars(r)["taskID"]

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&confirmTaskDeliverRequest)
	if strings.TrimSpace(taskId) == "" {
		confirmTaskDeliverResponse.Status = "error"
		confirmTaskDeliverResponse.Msg = "No TaskId"
		goto Label1
	}
	if err != nil {
		log.Println(err)
		goto Label0
	}
	fmt.Printf("Confirm task, taskDeliverId:%v,taskDeliverRate:%v\n", confirmTaskDeliverRequest.DeliverId, confirmTaskDeliverRequest.DeliverRate)
	_, err = db.Db.Exec("UPDATE task SET task_deliver_id = ? ,task_deliver_rate = ?,task_step = 2,task_complete = ? WHERE task_id =?", confirmTaskDeliverRequest.DeliverId, confirmTaskDeliverRequest.DeliverRate, time.Now(), taskId)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	_, err = db.Db.Exec("DELETE FROM task_bid WHERE task_id = ? ", taskId)
	if err != nil {
		log.Println(err)
	}
	confirmTaskDeliverResponse.Status = "success"
	confirmTaskDeliverResponse.Msg = "owner expected rate updated success"

Label0:
	if confirmTaskDeliverResponse.Status != "success" {
		confirmTaskDeliverResponse.Status = "error"
		confirmTaskDeliverResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(confirmTaskDeliverResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
