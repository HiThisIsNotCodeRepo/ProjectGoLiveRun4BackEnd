package tasks_enquiry

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

const SQLOngoingNewTaskIdentityUser = `select 
task_id,
task_title,
task_description,
task_category_id,
task_from,
task_to,
task_create,
task_start,
IFNULL(task_complete,''),
task_duration,
task_step,
task_owner_id,
task_owner_rate,
IFNULL(task_deliver_id,''),
IFNULL(task_deliver_rate,0) 
from task where 1=1 `

const SQLOngoingNewTaskBidder = ` select task_bidder_id,task_bidder_rate from task_bid where task_id = ? `

const OnlyMeFilterCondition = ` AND task_owner_id=? `
const ExcludeMeFilterCondition = ` AND task_owner_id<>? `
const OnGoingFilterCondition = ` AND task_step = 0 OR task_step = 1 order by task_start`

const SQLOngoingNewTaskIdentityTask = `select 
task_id,
task_title,
task_description,
task_category_id,
task_from,
task_to,
task_create,
task_start,
IFNULL(task_complete,''),
task_duration,
task_step,
task_owner_id,
task_owner_rate,
IFNULL(task_deliver_id,''),
IFNULL(task_deliver_rate,0) 
from task where task_id = ? `

type TaskResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Tasks  []Task `json:"tasks"`
}
type Task struct {
	No              int      `json:"no"`
	TaskId          string   `json:"taskId"`
	TaskTitle       string   `json:"taskTitle"`
	TaskDescription string   `json:"taskDescription"`
	TaskCategoryId  int      `json:"taskCategoryId"`
	TaskFrom        string   `json:"taskFrom"`
	TaskTo          string   `json:"taskTo"`
	TaskCreate      string   `json:"taskCreate"`
	TaskStart       string   `json:"taskStart"`
	TaskComplete    string   `json:"taskComplete"`
	TaskDuration    int      `json:"taskDuration"`
	TaskStep        int      `json:"taskStep"`
	TaskOwnerId     string   `json:"taskOwnerId"`
	TaskOwnerRate   int      `json:"taskOwnerRate"`
	TaskDeliverId   string   `json:"taskDeliverId"`
	TaskDeliverRate int      `json:"taskDeliverRate"`
	Bidders         []Bidder `json:"bidders"`
}

type Bidder struct {
	No             int    `json:"no"`
	TaskBidderId   string `json:"taskBidderId"`
	TaskBidderRate int    `json:"taskBidderRate"`
}

func TaskEnquiry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var taskEqnuiryResponse TaskResponse
	var err error
	var getAllRowsForTask *sql.Rows
	var getAllRowsForBidder *sql.Rows
	var option string
	var category string
	var finalSql string
	var identity string
	fmt.Printf("task enquiry->request URI:%v\n", r.RequestURI)
	tasks := make([]Task, 0)
	bidders := make([]Bidder, 0)
	encoder := json.NewEncoder(w)
	id := mux.Vars(r)["id"]
	option = r.URL.Query().Get("option")
	category = r.URL.Query().Get("category")
	identity = r.URL.Query().Get("identity")
	if strings.TrimSpace(id) == "" {
		taskEqnuiryResponse.Status = "error"
		taskEqnuiryResponse.Msg = "no id"
		goto Label1
	}
	if strings.TrimSpace(option) == "" {
		taskEqnuiryResponse.Status = "error"
		taskEqnuiryResponse.Msg = "no status"
		goto Label1
	}
	if strings.TrimSpace(identity) == "" {
		taskEqnuiryResponse.Status = "error"
		taskEqnuiryResponse.Msg = "no identity"
		goto Label1
	}
	fmt.Printf("id:%v,option:%v,category:%v,identity:%v\n", id, option, category, identity)

	if identity == "user" {
		if option == "on-going" && category == "only-me" {
			finalSql = fmt.Sprintf("%s%s%s", SQLOngoingNewTaskIdentityUser, OnlyMeFilterCondition, OnGoingFilterCondition)

		} else if option == "on-going" && category == "exclude-me" {
			finalSql = fmt.Sprintf("%s%s%s", SQLOngoingNewTaskIdentityUser, ExcludeMeFilterCondition, OnGoingFilterCondition)
		}
	} else if identity == "task" {
		finalSql = SQLOngoingNewTaskIdentityTask
	}
	getAllRowsForTask, err = db.Db.Query(finalSql, id)
	defer getAllRowsForTask.Close()
	if err != nil {
		log.Println(err)
		goto Label0
	}
	if getAllRowsForTask != nil {
		i := 1
		for getAllRowsForTask.Next() {
			var taskId string
			var taskTitle string
			var taskDescription string
			var taskCategoryId int
			var taskFrom string
			var taskTo string
			var taskCreate string
			var taskStart string
			var taskComplete string
			var taskDuration int
			var taskStep int
			var taskOwnerId string
			var taskOwnerRate int
			var deliverId string
			var deliverRate int

			err = getAllRowsForTask.Scan(&taskId, &taskTitle, &taskDescription, &taskCategoryId, &taskFrom, &taskTo, &taskCreate, &taskStart, &taskComplete, &taskDuration, &taskStep, &taskOwnerId, &taskOwnerRate, &deliverId, &deliverRate)
			if err != nil {
				log.Println(err)
				goto Label0
			}
			fmt.Printf("taskId:%v,taskTitle:%v,taskDescription:%v,taskCategoryId:%v,taskFrom:%v,taskTo:%v,taskCreate:%v,taskStart:%v,taskComplete:%v,taskDuration:%v,taskStep:%v,taskOwnerId:%v,taskOwnerRate:%v,deliverId:%v,deliverRate:%v\n", taskId, taskTitle, taskDescription, taskCategoryId, taskFrom, taskTo, taskCreate, taskStart, taskComplete, taskDuration, taskStep, taskOwnerId, taskOwnerRate, deliverId, deliverRate)
			getAllRowsForBidder, err = db.Db.Query(SQLOngoingNewTaskBidder, taskId)
			if err != nil {
				log.Println(err)
				goto Label0
			}

			tasks = append(tasks, Task{
				No:              i,
				TaskId:          taskId,
				TaskTitle:       taskTitle,
				TaskDescription: taskDescription,
				TaskCategoryId:  taskCategoryId,
				TaskFrom:        taskFrom,
				TaskTo:          taskTo,
				TaskCreate:      taskCreate,
				TaskStart:       taskStart,
				TaskComplete:    taskComplete,
				TaskDuration:    taskDuration,
				TaskStep:        taskStep,
				TaskOwnerId:     taskOwnerId,
				TaskOwnerRate:   taskOwnerRate,
				TaskDeliverId:   deliverId,
				TaskDeliverRate: deliverRate,
			})

			if getAllRowsForBidder != nil {
				j := 1
				for getAllRowsForBidder.Next() {
					var taskBidderId string
					var taskBidderRate int
					err = getAllRowsForBidder.Scan(&taskBidderId, &taskBidderRate)
					if err != nil {
						log.Println(err)
						goto Label0
					}
					fmt.Printf("taskBidderId:%v,taskBidderRate:%v\n", taskBidderId, taskBidderRate)
					bidders = append(bidders, Bidder{
						No:             j,
						TaskBidderId:   taskBidderId,
						TaskBidderRate: taskBidderRate,
					})
					j++
				}
				if j == 1 {
					fmt.Printf("no bidder for option:%v,category:%v\n", option, category)
				} else {
					tasks[i].Bidders = bidders
				}
			} else {
				fmt.Printf("no bidder for option:%v,category:%v\n", option, category)
			}
			i++
		}
		if i == 1 {
			fmt.Printf("no task for option:%v,category:%v\n", option, category)
		}
	} else {
		fmt.Printf("no task for option:%v,category:%v\n", option, category)
	}

	taskEqnuiryResponse.Status = "success"
	taskEqnuiryResponse.Msg = fmt.Sprintf("tasks:%v", tasks)
	taskEqnuiryResponse.Tasks = tasks
Label0:
	if taskEqnuiryResponse.Status != "success" {
		taskEqnuiryResponse.Status = "error"
		taskEqnuiryResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(taskEqnuiryResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
