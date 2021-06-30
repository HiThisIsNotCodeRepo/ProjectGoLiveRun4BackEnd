package my_info_earning

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"strings"
)

const SQLEarningDataSource = `select task_complete ,task_title,task_category_id,task_owner_id,task_deliver_id ,task_from,task_to,task_deliver_rate,datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),curdate()) from task
where task_deliver_id=?  AND task_step = 2 AND datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),curdate()) > -30 order by task_complete desc`

type EarningDataSourceResponse struct {
	Status string        `json:"status"`
	Msg    string        `json:"msg"`
	Tasks  []EarningTask `json:"tasks"`
}
type EarningTask struct {
	No               int    `json:"no"`
	CompleteDateTime string `json:"completeDateTime"`
	TaskTitle        string `json:"taskTitle"`
	TaskCategoryId   int    `json:"taskCategoryId"`
	TaskOwnerId      string `json:"taskOwnerId"`
	TaskDeliveredId  string `json:"taskDeliveredId"`
	TaskFrom         string `json:"taskFrom"`
	TaskTo           string `json:"taskTo"`
	TaskDeliverRate  int    `json:"taskDeliverRate"`
}

func DataSource(w http.ResponseWriter, r *http.Request) {
	var getSpendingTaskResponse EarningDataSourceResponse
	var err error
	var getAllRows *sql.Rows
	tasks := make([]EarningTask, 0)
	fmt.Printf("datasource->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	userID := strings.TrimSpace(mux.Vars(r)["userID"])
	if strings.TrimSpace(userID) == "" {
		getSpendingTaskResponse.Status = "error"
		getSpendingTaskResponse.Msg = "no userID"
		goto Label1
	}

	getAllRows, err = db.Db.Query(SQLEarningDataSource, userID)
	defer getAllRows.Close()
	if err != nil {
		log.Println(err)
		goto Label0
	}
	if getAllRows != nil {
		i := 1
		for getAllRows.Next() {
			var taskCompleteDate string
			var taskTitle string
			var taskCategoryId int
			var taskOwnerId string
			var taskDeliveredId string
			var taskFrom string
			var taskTo string
			var taskDeliveryRate int
			var diff int
			err = getAllRows.Scan(&taskCompleteDate, &taskTitle, &taskCategoryId, &taskOwnerId, &taskDeliveredId, &taskFrom, &taskTo, &taskDeliveryRate, &diff)
			if err != nil {
				log.Println(err)
				goto Label0
			}
			fmt.Printf("taskCompleteDate:%v,taskTitle:%v,taskCategoryId:%v,taskOwnerId:%v,taskDeliveredId:%v,taskFrom:%v,taskTo:%v,expense:%v,diff:%v\n", taskCompleteDate, taskTitle, taskCategoryId, taskOwnerId, taskDeliveredId, taskFrom, taskTo, taskDeliveryRate, diff)
			tasks = append(tasks, EarningTask{
				No:               i,
				CompleteDateTime: taskCompleteDate,
				TaskTitle:        taskTitle,
				TaskCategoryId:   taskCategoryId,
				TaskOwnerId:      getLast4Char(taskOwnerId),
				TaskDeliveredId:  getLast4Char(taskDeliveredId),
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
func getLast4Char(str string) string {
	strArr := strings.Split(str, "-")
	lastStrByteArr := []byte(strArr[len(strArr)-1])
	var last4ByteArr []byte
	for i := len(lastStrByteArr) - 4; i < len(lastStrByteArr); i++ {
		last4ByteArr = append(last4ByteArr, lastStrByteArr[i])
	}
	return string(last4ByteArr)
}