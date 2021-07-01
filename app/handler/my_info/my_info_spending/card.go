package my_info_spending

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"paotui.sg/app/handler/error_util"
	"strings"
)

const BuyNecessityFilterCondition = "AND task_category_id = 0"
const FoodDeliveryFilterCondition = "AND task_category_id = 1"
const SendDocumentFilterCondition = "AND task_category_id = 2"
const OtherFilterCondition = "AND task_category_id = 3"
const YesterdayFilterCondition = "AND TO_DAYS(NOW()) - TO_DAYS(task_complete) <= 1"
const TwoDaysAgoFilterCondition = "AND TO_DAYS(NOW()) -TO_DAYS(task_complete) > 1 AND TO_DAYS(NOW()) -TO_DAYS(task_complete) <= 2"
const ThreeDaysAgoFilterCondition = "AND TO_DAYS(NOW()) -TO_DAYS(task_complete) > 2 AND TO_DAYS(NOW()) -TO_DAYS(task_complete) <= 3"

type CardResponse struct {
	Status    string `json:"status"`
	Msg       string `json:"msg"`
	TaskCount int    `json:"taskCount"`
	TaskSpend int    `json:"taskSpend"`
}

func Card(w http.ResponseWriter, r *http.Request) {
	defer error_util.ErrorHandle(w)
	var baseSql string
	var finalSql string
	var getCardResponse CardResponse
	var err error
	var taskCount int
	var taskSpend int
	var categoryFilterCondition string
	var dateFilterCondition string
	encoder := json.NewEncoder(w)
	userID := strings.TrimSpace(mux.Vars(r)["userID"])
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	fmt.Printf("card->request URI:%v\n", r.RequestURI)
	fmt.Printf("card->userID:%v,date:%v,category:%v\n", userID, date, category)
	if strings.TrimSpace(userID) == "" {
		getCardResponse.Status = "error"
		getCardResponse.Msg = "no userID"
		goto Label1
	}

	if category == "buy-necessity" {
		categoryFilterCondition = BuyNecessityFilterCondition
	} else if category == "food-delivery" {
		categoryFilterCondition = FoodDeliveryFilterCondition
	} else if category == "send-document" {
		categoryFilterCondition = SendDocumentFilterCondition
	} else if category == "other" {
		categoryFilterCondition = OtherFilterCondition
	} else {
		goto Label0
	}

	if date == "yesterday" {
		dateFilterCondition = YesterdayFilterCondition
	} else if date == "two-days-ago" {
		dateFilterCondition = TwoDaysAgoFilterCondition
	} else if date == "three-days-ago" {
		dateFilterCondition = ThreeDaysAgoFilterCondition
	} else {
		goto Label0
	}
	baseSql = "SELECT count(*) FROM task WHERE task_owner_id = ? AND task_step = 2"
	finalSql = fmt.Sprintf("%s %s %s", baseSql, categoryFilterCondition, dateFilterCondition)
	fmt.Printf("card->count sql:%v\n", finalSql)
	err = db.Db.QueryRow(finalSql, userID).Scan(&taskCount)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	baseSql = "SELECT IFNULL(sum(task_deliver_rate),0) FROM task WHERE task_owner_id =?"
	finalSql = fmt.Sprintf("%s %s %s", baseSql, categoryFilterCondition, dateFilterCondition)
	err = db.Db.QueryRow(finalSql, userID).Scan(&taskSpend)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	fmt.Printf("card->sum sql:%v\n", finalSql)
	getCardResponse.Status = "success"
	getCardResponse.Msg = fmt.Sprintf("task count:%v,task spend:%v", taskCount, taskSpend)
	getCardResponse.TaskCount = taskCount
	getCardResponse.TaskSpend = taskSpend

Label0:
	if getCardResponse.Status != "success" {
		getCardResponse.Status = "error"
		getCardResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(getCardResponse)
	if encodeErr != nil {
w.WriteHeader(http.StatusInternalServerError)
	}
}
