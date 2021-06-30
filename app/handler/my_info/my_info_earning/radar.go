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
	"time"
)

const SQLEarningRadarNonSun = `select task_complete , task_id,task_category_id,task_deliver_rate ,datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),subdate(curdate(),date_format(curdate(),'%w')-7)) 
from task where task_deliver_id=? AND task_step = 2 AND datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),subdate(curdate(),date_format(curdate(),'%w')-7)) > -14`
const SQLEarningRadarSun = `select task_complete , task_id,task_category_id,task_deliver_rate ,datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),curdate())
from task where task_deliver_id=? AND task_step = 2 AND datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),curdate()) > -14`

type RadarResponse struct {
	Status       string `json:"status"`
	Msg          string `json:"msg"`
	BuyNecessity int    `json:"buyNecessity"`
	FoodDelivery int    `json:"foodDelivery"`
	SendDocument int    `json:"sendDocument"`
	Other        int    `json:"other"`
}

func Radar(w http.ResponseWriter, r *http.Request) {

	var getEarningRadar RadarResponse
	var err error
	var getAllRows *sql.Rows
	var lastWeekFlag bool
	var buyNecessity int
	var foodDelivery int
	var sendDocument int
	var other int
	var buyNecessityArr = make([]int, 2)
	var foodDeliveryArr = make([]int, 2)
	var sendDocumentArr = make([]int, 2)
	var otherArr = make([]int, 2)
	var finalSQL string
	fmt.Printf("radar->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	userID := strings.TrimSpace(mux.Vars(r)["userID"])
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	if strings.TrimSpace(userID) == "" {
		getEarningRadar.Status = "error"
		getEarningRadar.Msg = "no userID"
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
		finalSQL = SQLEarningRadarSun
	} else {
		finalSQL = SQLEarningRadarNonSun
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

			if 13+diff < 7 {
				if taskCategoryId == 0 {

					buyNecessityArr[0] += taskDeliveryRate
				} else if taskCategoryId == 1 {

					foodDeliveryArr[0] += taskDeliveryRate
				} else if taskCategoryId == 2 {

					sendDocumentArr[0] += taskDeliveryRate
				} else {

					otherArr[0] += taskDeliveryRate
				}
			} else {
				if taskCategoryId == 0 {

					buyNecessityArr[1] += taskDeliveryRate
				} else if taskCategoryId == 1 {

					foodDeliveryArr[1] += taskDeliveryRate
				} else if taskCategoryId == 2 {

					sendDocumentArr[1] += taskDeliveryRate
				} else {

					otherArr[1] += taskDeliveryRate
				}
			}
		}
	}

	if lastWeekFlag {
		buyNecessity = buyNecessityArr[0]
		foodDelivery = foodDeliveryArr[0]
		sendDocument = sendDocumentArr[0]
		other = otherArr[0]

	} else {
		buyNecessity = buyNecessityArr[1]
		foodDelivery = foodDeliveryArr[1]
		sendDocument = sendDocumentArr[1]
		other = otherArr[1]

	}
	getEarningRadar.Status = "success"
	getEarningRadar.Msg = fmt.Sprintf("BuyNecessity:%v,FoodDelivery:%v,SendDocument:%v,Other:%v\n", buyNecessity, foodDelivery, sendDocument, other)
	getEarningRadar.BuyNecessity = buyNecessity
	getEarningRadar.FoodDelivery = foodDelivery
	getEarningRadar.SendDocument = sendDocument
	getEarningRadar.Other = other

Label0:
	if getEarningRadar.Status != "success" {
		getEarningRadar.Status = "error"
		getEarningRadar.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(getEarningRadar)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
