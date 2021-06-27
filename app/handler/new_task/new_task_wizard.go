package new_task

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"paotui.sg/app/db"
	"strings"
	"time"
)

type NewTaskWizardRequest struct {
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

type NewTaskWizardResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func GetNewTaskWizardRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var newTaskWizardRequest = NewTaskWizardRequest{}
	var getNewTaskWizardRequestResponse NewTaskWizardResponse
	var err error
	fmt.Printf("request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newTaskWizardRequest)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	fmt.Printf("newTaskWizardRequest:%v\n", newTaskWizardRequest)
	if strings.TrimSpace(newTaskWizardRequest.TaskTitle) != "" && strings.TrimSpace(newTaskWizardRequest.From) != "" && strings.TrimSpace(newTaskWizardRequest.To) != "" && newTaskWizardRequest.Category >= 0 && newTaskWizardRequest.Category <= 3 && newTaskWizardRequest.ExpectedRate >= 1 && newTaskWizardRequest.ExpectedRate <= 1000 && newTaskWizardRequest.Duration >= 10 && newTaskWizardRequest.Duration <= 100 && strings.TrimSpace(newTaskWizardRequest.Start) != "" {
		taskId := uuid.NewV4().String()
		taskTitle := newTaskWizardRequest.TaskTitle
		taskDescription := newTaskWizardRequest.Description
		taskCategory := newTaskWizardRequest.Category
		fromLocation := newTaskWizardRequest.From
		toLocation := newTaskWizardRequest.To
		createTime := time.Now()
		startTime := strings.Replace(newTaskWizardRequest.Start, "T", " ", 1)
		duration := newTaskWizardRequest.Duration
		taskStep := 0
		taskOwnerId := newTaskWizardRequest.TaskOwnerId
		ownerRate := newTaskWizardRequest.ExpectedRate
		_, err = db.Db.Exec("INSERT INTO task (task_id,task_title,task_description,task_category_id,task_from,task_to,task_create,task_start,task_duration,task_step,task_owner_id,task_owner_rate) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)", taskId, taskTitle, taskDescription, taskCategory, fromLocation, toLocation, createTime, startTime, duration, taskStep, taskOwnerId, ownerRate)
		if err != nil {
			fmt.Println(err)
			goto Label0
		}
		getNewTaskWizardRequestResponse.Status = "success"
		getNewTaskWizardRequestResponse.Msg = "new task added success"
	} else {
		getNewTaskWizardRequestResponse.Status = "error"
		getNewTaskWizardRequestResponse.Msg = "new task wizard data error"
		goto Label1
	}

Label0:
	if getNewTaskWizardRequestResponse.Status != "success" {
		getNewTaskWizardRequestResponse.Status = "error"
		getNewTaskWizardRequestResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(getNewTaskWizardRequestResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
