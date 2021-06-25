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

const SQLEarningRadar = `select task_complete , task_id,task_category_id,task_deliver_rate ,datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),subdate(curdate(),date_format(curdate(),'%w')-7)) 
from task where task_deliver_id=? AND datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),subdate(curdate(),date_format(curdate(),'%w')-7)) > -14`

type EarningRadarResponse struct {
	Status       string `json:"status"`
	Msg          string `json:"msg"`
	BuyNecessity int    `json:"buyNecessity"`
	FoodDelivery int    `json:"foodDelivery"`
	SendDocument int    `json:"sendDocument"`
	Other        int    `json:"other"`
}

func GetEarningRadar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var getEarningRadar EarningRadarResponse
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
	fmt.Printf("request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	userID := mux.Vars(r)["userID"]
	if strings.TrimSpace(userID) == "" {
		getEarningRadar.Status = "error"
		getEarningRadar.Msg = "no userID"
		goto Label1
	}

	if strings.Contains(r.RequestURI, "last-week") {
		lastWeekFlag = true
	} else if strings.Contains(r.RequestURI, "this-week") {
		lastWeekFlag = false
	} else {
		goto Label0
	}

	getAllRows, err = db.Db.Query(SQLEarningRadar, userID)
	if err != nil {
		fmt.Println(err)
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
				fmt.Println(err)
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
