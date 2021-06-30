package my_info_spending

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"strings"
	"time"
)

const SQLSpendingSummaryNonSun = `select task_complete , task_id,task_category_id,task_deliver_rate ,datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),subdate(curdate(),date_format(curdate(),'%w')-7)) 
from task where task_owner_id=? AND task_step = 2 AND datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),subdate(curdate(),date_format(curdate(),'%w')-7)) > -14`
const SQLSpendingSummarySun = `select task_complete , task_id,task_category_id,task_deliver_rate ,datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),curdate())
from task where task_owner_id=? AND task_step = 2 AND datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),curdate()) > -14`

type SummaryResponse struct {
	Status                    string `json:"status"`
	Msg                       string `json:"msg"`
	LineData                  []int  `json:"lineData"`
	ColumnData                []int  `json:"columnData"`
	TotalTasks                int    `json:"totalTasks"`
	DollarSpent               int    `json:"dollarSpent"`
	BuyNecessityCategoryCount int    `json:"buyNecessity"`
	FoodDeliveryCategoryCount int    `json:"foodDelivery"`
	SendDocumentCategoryCount int    `json:"sendDocument"`
	OtherCategoryCount        int    `json:"other"`
}

func Summary(w http.ResponseWriter, r *http.Request) {
	var getSummaryResponse SummaryResponse
	var err error
	var getAllRows *sql.Rows
	var lastWeekFlag bool
	var finalSQL string
	lineData := make([]int, 14)
	columnData := make([]int, 14)
	totalTask := make([]int, 2)
	dollarSpent := make([]int, 2)
	buyNecessityCategoryCount := make([]int, 2)
	foodDeliveryCategoryCount := make([]int, 2)
	sendDocumentCategoryCount := make([]int, 2)
	otherCategoryCount := make([]int, 2)
	encoder := json.NewEncoder(w)
	userID := strings.TrimSpace(mux.Vars(r)["userID"])
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	fmt.Printf("summary->URI:%v\n", r.RequestURI)
	fmt.Printf("summary->userID:%v,date:%v\n", userID, date)
	if strings.TrimSpace(userID) == "" {
		getSummaryResponse.Status = "error"
		getSummaryResponse.Msg = "no userID"
		goto Label1
	}

	if date == "last-week" {
		lastWeekFlag = true
	} else if date == "this-week" {
		lastWeekFlag = false
	} else {
		goto Label0
	}
	if time.Now().Weekday() == 0 {
		finalSQL = SQLSpendingSummarySun
	} else {
		finalSQL = SQLSpendingSummaryNonSun
	}
	getAllRows, err = db.Db.Query(finalSQL, userID)
	defer getAllRows.Close()
	if err != nil {
		log.Println(err)
		goto Label0
	}
	if getAllRows != nil {
		for getAllRows.Next() {
			var taskCompleteDate string
			var taskId string
			var taskCategoryId int
			var taskDeliveryRate int
			var diff int
			err = getAllRows.Scan(&taskCompleteDate, &taskId, &taskCategoryId, &taskDeliveryRate, &diff)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("taskCompleteDate:%v,taskID:%v,taskCategoryId:%v,expense:%v,diff:%v\n", taskCompleteDate, taskId, taskCategoryId, taskDeliveryRate, diff)
			lineData[13+diff] += taskDeliveryRate
			columnData[13+diff] += 1
			if 13+diff < 7 {
				totalTask[0] += 1
				dollarSpent[0] += taskDeliveryRate
				if taskCategoryId == 0 {
					buyNecessityCategoryCount[0] += 1
				} else if taskCategoryId == 1 {
					foodDeliveryCategoryCount[0] += 1
				} else if taskCategoryId == 2 {
					sendDocumentCategoryCount[0] += 1
				} else {
					otherCategoryCount[0] += 1
				}
			} else {
				totalTask[1] += 1
				dollarSpent[1] += taskDeliveryRate
				if taskCategoryId == 0 {
					buyNecessityCategoryCount[1] += 1
				} else if taskCategoryId == 1 {
					foodDeliveryCategoryCount[1] += 1
				} else if taskCategoryId == 2 {
					sendDocumentCategoryCount[1] += 1
				} else {
					otherCategoryCount[1] += 1
				}
			}
		}
	}

	if lastWeekFlag {
		lineData = lineData[0:7]
		columnData = columnData[0:7]
		getSummaryResponse.TotalTasks = totalTask[0]
		getSummaryResponse.DollarSpent = dollarSpent[0]
		getSummaryResponse.BuyNecessityCategoryCount = buyNecessityCategoryCount[0]
		getSummaryResponse.FoodDeliveryCategoryCount = foodDeliveryCategoryCount[0]
		getSummaryResponse.SendDocumentCategoryCount = sendDocumentCategoryCount[0]
		getSummaryResponse.OtherCategoryCount = otherCategoryCount[0]
	} else {
		lineData = lineData[7:]
		columnData = columnData[7:]
		getSummaryResponse.TotalTasks = totalTask[1]
		getSummaryResponse.DollarSpent = dollarSpent[1]
		getSummaryResponse.BuyNecessityCategoryCount = buyNecessityCategoryCount[1]
		getSummaryResponse.FoodDeliveryCategoryCount = foodDeliveryCategoryCount[1]
		getSummaryResponse.SendDocumentCategoryCount = sendDocumentCategoryCount[1]
		getSummaryResponse.OtherCategoryCount = otherCategoryCount[1]
	}
	getSummaryResponse.Status = "success"
	getSummaryResponse.Msg = fmt.Sprintf("line data:%v,column data:%v,totalTask:%v,dollarSpent:%v,buy necessity category count:%v,food delivery category count:%v,send document category count:%v,other count:%v\n", lineData, columnData, getSummaryResponse.TotalTasks, getSummaryResponse.DollarSpent, getSummaryResponse.BuyNecessityCategoryCount, getSummaryResponse.FoodDeliveryCategoryCount, getSummaryResponse.SendDocumentCategoryCount, getSummaryResponse.OtherCategoryCount)
	getSummaryResponse.LineData = lineData
	getSummaryResponse.ColumnData = columnData
Label0:
	if getSummaryResponse.Status != "success" {
		getSummaryResponse.Status = "error"
		getSummaryResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(getSummaryResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
