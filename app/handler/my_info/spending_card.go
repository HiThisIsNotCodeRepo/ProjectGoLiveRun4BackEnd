package my_info

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"paotui.sg/app/db"
	"strings"
)

const BuyNecessityFilterCondition = "AND task_category_id = 0"
const FoodDeliveryFilterCondition = "AND task_category_id = 1"
const SendDocumentFilterCondition = "AND task_category_id = 2"
const OtherFilterCondition = "AND task_category_id = 3"
const YesterdayFilterCondition = "AND TO_DAYS(NOW()) - TO_DAYS(task_complete) <= 1"
const TwoDaysAgoFilterCondition = "AND TO_DAYS(NOW()) -TO_DAYS(task_complete) > 1 AND TO_DAYS(NOW()) -TO_DAYS(task_complete) <= 2"
const ThreeDaysAgoFilterCondition = "AND TO_DAYS(NOW()) -TO_DAYS(task_complete) > 2 AND TO_DAYS(NOW()) -TO_DAYS(task_complete) <= 3"

type SpendingCardResponse struct {
	Status    string `json:"status"`
	Msg       string `json:"msg"`
	TaskCount int    `json:"taskCount"`
	TaskSpend int    `json:"taskSpend"`
}

func GetSpendingCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var baseSql string
	var finalSql string
	var getBuyNecessitySpendingYesterdayResponse SpendingCardResponse
	var err error
	var taskCount int
	var taskSpend int
	var categoryFilterCondition string
	var dateFilterCondition string
	fmt.Printf("request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	userID := mux.Vars(r)["userID"]
	if strings.TrimSpace(userID) == "" {
		getBuyNecessitySpendingYesterdayResponse.Status = "error"
		getBuyNecessitySpendingYesterdayResponse.Msg = "no userID"
		goto Label1
	}

	if strings.Contains(r.RequestURI, "buy-necessity") {
		categoryFilterCondition = BuyNecessityFilterCondition
	} else if strings.Contains(r.RequestURI, "food-delivery") {
		categoryFilterCondition = FoodDeliveryFilterCondition
	} else if strings.Contains(r.RequestURI, "send-document") {
		categoryFilterCondition = SendDocumentFilterCondition
	} else if strings.Contains(r.RequestURI, "other") {
		categoryFilterCondition = OtherFilterCondition
	} else {
		goto Label0
	}

	if strings.Contains(r.RequestURI, "yesterday") {
		dateFilterCondition = YesterdayFilterCondition
	} else if strings.Contains(r.RequestURI, "two-days-ago") {
		dateFilterCondition = TwoDaysAgoFilterCondition
	} else if strings.Contains(r.RequestURI, "three-days-ago") {
		dateFilterCondition = ThreeDaysAgoFilterCondition
	} else {
		goto Label0
	}
	baseSql = "SELECT count(*) FROM task WHERE task_owner_id = ?"
	finalSql = fmt.Sprintf("%s %s %s", baseSql, categoryFilterCondition, dateFilterCondition)
	err = db.Db.QueryRow(finalSql, userID).Scan(&taskCount)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	baseSql = "SELECT IFNULL(sum(task_deliver_rate),0) FROM task WHERE task_owner_id =?"
	finalSql = fmt.Sprintf("%s %s %s", baseSql, categoryFilterCondition, dateFilterCondition)
	err = db.Db.QueryRow(finalSql, userID).Scan(&taskSpend)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	getBuyNecessitySpendingYesterdayResponse.Status = "success"
	getBuyNecessitySpendingYesterdayResponse.Msg = fmt.Sprintf("task count:%v,task spend:%v", taskCount, taskSpend)
	getBuyNecessitySpendingYesterdayResponse.TaskCount = taskCount
	getBuyNecessitySpendingYesterdayResponse.TaskSpend = taskSpend

Label0:
	if getBuyNecessitySpendingYesterdayResponse.Status != "success" {
		getBuyNecessitySpendingYesterdayResponse.Status = "error"
		getBuyNecessitySpendingYesterdayResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(getBuyNecessitySpendingYesterdayResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
