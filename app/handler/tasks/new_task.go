package tasks

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"paotui.sg/app/handler/error_util"
	"strings"
	"time"
)

type NewTaskRequest struct {
	TaskTitle    string `json:"taskTitle"`
	TaskOwnerId  string `json:"taskOwnerId"`
	Description  string `json:"description"`
	From         string `json:"from"`
	To           string `json:"to"`
	Category     int    `json:"category"`
	ExpectedRate int    `json:"expectedRate"`
	Duration     int    `json:"duration"`
	Start        string `json:"start"`
}

type NewTaskResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func NewTask(w http.ResponseWriter, r *http.Request) {
	defer error_util.ErrorHandle(w)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var newTaskRequest = NewTaskRequest{}
	var newTasResponse NewTaskResponse
	var err error
	fmt.Printf("request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newTaskRequest)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	fmt.Printf("newTaskWizardRequest:%v\n", newTaskRequest)
	if strings.TrimSpace(newTaskRequest.TaskTitle) != "" && strings.TrimSpace(newTaskRequest.From) != "" && strings.TrimSpace(newTaskRequest.To) != "" && newTaskRequest.Category >= 0 && newTaskRequest.Category <= 3 && newTaskRequest.ExpectedRate >= 1 && newTaskRequest.ExpectedRate <= 1000 && newTaskRequest.Duration >= 10 && newTaskRequest.Duration <= 100 && strings.TrimSpace(newTaskRequest.Start) != "" {
		taskId := uuid.NewV4().String()
		taskTitle := newTaskRequest.TaskTitle
		taskDescription := newTaskRequest.Description
		taskCategory := newTaskRequest.Category
		fromLocation := newTaskRequest.From
		toLocation := newTaskRequest.To
		createTime := time.Now()
		startTime := strings.Replace(newTaskRequest.Start, "T", " ", 1)
		duration := newTaskRequest.Duration
		taskStep := 0
		taskOwnerId := newTaskRequest.TaskOwnerId
		ownerRate := newTaskRequest.ExpectedRate
		_, err = db.Db.Exec("INSERT INTO task (task_id,task_title,task_description,task_category_id,task_from,task_to,task_create,task_start,task_duration,task_step,task_owner_id,task_owner_rate) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)", taskId, taskTitle, taskDescription, taskCategory, fromLocation, toLocation, createTime, startTime, duration, taskStep, taskOwnerId, ownerRate)
		if err != nil {
			log.Println(err)
			goto Label0
		}
		newTasResponse.Status = "success"
		newTasResponse.Msg = "new task added success"
	} else {
		newTasResponse.Status = "error"
		newTasResponse.Msg = "new task wizard data error"
		goto Label1
	}

Label0:
	if newTasResponse.Status != "success" {
		newTasResponse.Status = "error"
		newTasResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(newTasResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
