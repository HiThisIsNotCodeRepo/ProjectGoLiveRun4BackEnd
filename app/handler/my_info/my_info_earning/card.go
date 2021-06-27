package my_info_earning

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/db"
	"strings"
)

const SQLEarningCard = `select task_complete ,task_title,task_category_id,task_owner_id,task_deliver_id ,task_from,task_to,task_deliver_rate,datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),DATE_FORMAT(curdate(), '%Y-%m-%d')) from task
where task_deliver_id=? AND datediff(DATE_FORMAT(task_complete, '%Y-%m-%d'),DATE_FORMAT(curdate(), '%Y-%m-%d')) > -11 order by task_complete DESC`

type CardResponse struct {
	Status            string `json:"status"`
	Msg               string `json:"msg"`
	PastTwoDaysTotal  int    `json:"pastTwoDaysTotal"`
	PastFiveDaysTotal int    `json:"pastFiveDaysTotal"`
	PastTenDaysTotal  int    `json:"pastTenDaysTotal"`
	PastTwoDays       []int  `json:"pastTwoDays"`
	PastFiveDays      []int  `json:"pastFiveDays"`
	PastTenDays       []int  `json:"pastTenDays"`
}

func Card(w http.ResponseWriter, r *http.Request) {
	var getEarningCardResponse CardResponse
	var err error
	var getAllRows *sql.Rows
	var pastTwoDays = make([]int, 2)
	var pastFiveDays = make([]int, 5)
	var pastTenDays = make([]int, 10)
	var pastTwoDaysTotal int
	var pastFiveDaysTotal int
	var pastTenDaysTotal int
	fmt.Printf("card->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	userID := mux.Vars(r)["userID"]
	if strings.TrimSpace(userID) == "" {
		getEarningCardResponse.Status = "error"
		getEarningCardResponse.Msg = "no userID"
		goto Label1
	}

	getAllRows, err = db.Db.Query(SQLEarningCard, userID)
	defer getAllRows.Close()
	if err != nil {
		fmt.Println(err)
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
				fmt.Println(err)
			}
			fmt.Printf("taskCompleteDate:%v,taskTitle:%v,taskCategoryId:%v,taskOwnerId:%v,taskDeliveredId:%v,taskFrom:%v,taskTo:%v,expense:%v,diff:%v\n", taskCompleteDate, taskTitle, taskCategoryId, taskOwnerId, taskDeliveredId, taskFrom, taskTo, taskDeliveryRate, diff)
			if diff < 0 {
				if diff >= -2 {
					pastTwoDays[2+diff] += taskDeliveryRate
					pastFiveDays[5+diff] += taskDeliveryRate
					pastTenDays[10+diff] += taskDeliveryRate
					pastTwoDaysTotal += taskDeliveryRate
					pastFiveDaysTotal += taskDeliveryRate
					pastTenDaysTotal += taskDeliveryRate
				} else if diff >= -5 {
					pastFiveDays[5+diff] += taskDeliveryRate
					pastTenDays[10+diff] += taskDeliveryRate
					pastFiveDaysTotal += taskDeliveryRate
					pastTenDaysTotal += taskDeliveryRate

				} else if diff >= -10 {
					pastTenDays[10+diff] += taskDeliveryRate
					pastTenDaysTotal += taskDeliveryRate
				}
			}
			i++
		}
	}
	getEarningCardResponse.Status = "success"
	getEarningCardResponse.Msg = fmt.Sprintf("pastTwoDays:%v,pastFiveDays:%v,pastTenDays:%v,pastTwoDaysTotal:%v,pastFiveDaysTotal:%v,pastTenDaysTotal:%v", pastTwoDays, pastFiveDays, pastTenDays, pastTwoDaysTotal, pastFiveDaysTotal, pastTenDaysTotal)
Label0:
	getEarningCardResponse.PastTwoDays = pastTwoDays
	getEarningCardResponse.PastFiveDays = pastFiveDays
	getEarningCardResponse.PastTenDays = pastTenDays
	getEarningCardResponse.PastTwoDaysTotal = pastTwoDaysTotal
	getEarningCardResponse.PastFiveDaysTotal = pastFiveDaysTotal
	getEarningCardResponse.PastTenDaysTotal = pastTenDaysTotal
	if getEarningCardResponse.Status != "success" {
		getEarningCardResponse.Status = "error"
		getEarningCardResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(getEarningCardResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
