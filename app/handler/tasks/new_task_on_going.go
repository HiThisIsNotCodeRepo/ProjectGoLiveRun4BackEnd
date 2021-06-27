package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/db"
	"strings"
)

const SQLOngoingNewTask = `select task_id,task_title,task_start ,task_category_id,task_owner_id,task_from,task_to,task_duration,task_owner_rate from task
where task_owner_id=? AND task_step = 0 order by task_start`

const SQLOngoingNewTaskBidder = `select task_bidder_id,task_bidder_rate from task_bid where task_id = ? `

type OnGoingNewTaskResponse struct {
	Status string           `json:"status"`
	Msg    string           `json:"msg"`
	Tasks  []OnGoingNewTask `json:"tasks"`
}
type OnGoingNewTask struct {
	No              int          `json:"no"`
	TaskTitle       string       `json:"taskTitle"`
	TaskStart       string       `json:"taskStart"`
	TaskCategoryId  int          `json:"taskCategoryId"`
	TaskOwnerId     string       `json:"taskOwnerId"`
	TaskFrom        string       `json:"taskFrom"`
	TaskTo          string       `json:"taskTo"`
	TaskDuration    int          `json:"taskDuration"`
	TaskOwnerRate   int          `json:"taskOwnerRate"`
	TaskDeliverRate int          `json:"taskDeliverRate"`
	Bidders         []BidderInfo `json:"bidders"`
}

type BidderInfo struct {
	TaskBidderId   string `json:"taskBidderId"`
	TaskBidderRate int    `json:"taskBidderRate"`
}

func GetOnGoingNewTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var getOnGoingNewTaskResponse OnGoingNewTaskResponse
	var err error
	var getAllRowsForTask *sql.Rows
	var getAllRowsForBidder *sql.Rows
	var status string
	fmt.Printf("request URI:%v\n", r.RequestURI)
	tasks := make([]OnGoingNewTaskResponse, 0)
	encoder := json.NewEncoder(w)
	userID := mux.Vars(r)["userID"]
	status = r.URL.Query().Get("status")
	if strings.TrimSpace(userID) == "" {
		getOnGoingNewTaskResponse.Status = "error"
		getOnGoingNewTaskResponse.Msg = "no userID"
		goto Label1
	}
	if strings.TrimSpace(status) == "" {
		getOnGoingNewTaskResponse.Status = "error"
		getOnGoingNewTaskResponse.Msg = "no status"
		goto Label1
	}
	if strings.TrimSpace(status) != "on-going" {
		getOnGoingNewTaskResponse.Status = "error"
		getOnGoingNewTaskResponse.Msg = "status not valid"
		goto Label1
	}
	fmt.Printf("userID:%v,status:%v\n", userID, status)
	getAllRowsForTask, err = db.Db.Query(SQLOngoingNewTask, userID)
	defer getAllRowsForTask.Close()
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	if getAllRowsForTask != nil {
		i := 1
		for getAllRowsForTask.Next() {
			var taskId string
			var taskTitle string
			var taskStart string
			var taskCategoryId int
			var taskOwnerId string
			var taskFrom string
			var taskTo string
			var taskDuration int
			var taskOwnerRate string
			err = getAllRowsForTask.Scan(&taskId, &taskTitle, &taskStart, &taskCategoryId, &taskOwnerId, &taskFrom, &taskTo, &taskDuration, &taskOwnerRate)
			if err != nil {
				fmt.Println(err)
				goto Label0
			}
			fmt.Printf("taskId:%v,taskTitle:%v,taskStart:%v,taskCategoryId:%v,taskOwnerId:%v,taskFrom:%v,taskTo:%v,taskDuration:%v,taskOwnerRate:%v\n", taskId, taskTitle, taskStart, taskCategoryId, taskOwnerId, taskFrom, taskTo, taskDuration, taskOwnerRate)
			getAllRowsForBidder, err = db.Db.Query(SQLOngoingNewTaskBidder, taskId)
			if err != nil {
				fmt.Println(err)
				goto Label0
			}
			if getAllRowsForBidder != nil {
				j := 1
				for getAllRowsForBidder.Next() {
					var taskBidderId string
					var taskBidderRate int
					err = getAllRowsForBidder.Scan(&taskBidderId, &taskBidderRate)
					if err != nil {
						fmt.Println(err)
						goto Label0
					}
					fmt.Printf("taskBidderId:%v,taskBidderRate:%v\n", taskBidderId, taskBidderRate)
				}
				j++
			}
			i++
		}
	}

	getOnGoingNewTaskResponse.Status = "success"
	getOnGoingNewTaskResponse.Msg = fmt.Sprintf("tasks:%v", tasks)
	//getOnGoingNewTaskResponse.Tasks = tasks
Label0:
	if getOnGoingNewTaskResponse.Status != "success" {
		getOnGoingNewTaskResponse.Status = "error"
		getOnGoingNewTaskResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(getOnGoingNewTaskResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
