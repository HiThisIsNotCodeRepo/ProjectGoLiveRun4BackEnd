package my_info

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/db"
	"strings"
)

const SQLTask = `select task_complete ,task_title,task_category_id ,task_from,task_to,task_deliver_rate,datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),subdate(curdate(),date_format(curdate(),'%w')-7)) from task
where task_owner_id=? AND datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),subdate(curdate(),date_format(curdate(),'%w')-7)) > -30`

type SpendingTaskResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Tasks  []Task `json:"tasks"`
}
type Task struct {
	No               int    `json:"no"`
	CompleteDateTime string `json:"completeDateTime"`
	TaskTitle        string `json:"taskTitle"`
	TaskCategoryId   int    `json:"taskCategoryId"`
	TaskFrom         string `json:"taskFrom"`
	TaskTo           string `json:"taskTo"`
	TaskDeliverRate  int    `json:"taskDeliverRate"`
}

func GetSpendingTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var getSpendingTaskResponse SpendingTaskResponse
	var err error
	var getAllRows *sql.Rows
	tasks := make([]Task, 0)
	fmt.Printf("request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	userID := mux.Vars(r)["userID"]
	if strings.TrimSpace(userID) == "" {
		getSpendingTaskResponse.Status = "error"
		getSpendingTaskResponse.Msg = "no userID"
		goto Label1
	}

	getAllRows, err = db.Db.Query(SQLTask, userID)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	if getAllRows == nil {
		tasks = nil
	} else {
		i := 1
		for getAllRows.Next() {
			var taskCompleteDate string
			var taskTitle string
			var taskCategoryId int
			var taskFrom string
			var taskTo string
			var taskDeliveryRate int
			var diff int
			err = getAllRows.Scan(&taskCompleteDate, &taskTitle, &taskCategoryId, &taskFrom, &taskTo, &taskDeliveryRate, &diff)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("taskCompleteDate:%v,taskTitle:%v,taskCategoryId:%v,taskFrom:%v,taskTo:%v,expense:%v,diff:%v\n", taskCompleteDate, taskTitle, taskCategoryId, taskFrom, taskTo, taskDeliveryRate, diff)
			tasks = append(tasks, Task{
				No:               i,
				CompleteDateTime: taskCompleteDate,
				TaskTitle:        taskTitle,
				TaskCategoryId:   taskCategoryId,
				TaskFrom:         taskFrom,
				TaskTo:           taskTo,
				TaskDeliverRate:  taskDeliveryRate})
			i++
		}
	}

	getSpendingTaskResponse.Status = "success"
	getSpendingTaskResponse.Msg = fmt.Sprintf("tasks:%v", tasks)
	getSpendingTaskResponse.Tasks = tasks
Label0:
	if getSpendingTaskResponse.Status != "success" {
		getSpendingTaskResponse.Status = "error"
		getSpendingTaskResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(getSpendingTaskResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
