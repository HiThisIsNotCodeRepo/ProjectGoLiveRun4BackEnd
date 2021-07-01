package bid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"paotui.sg/app/handler/error_util"
	"strings"
)

type NewBidRequest struct {
	TaskId         string `json:"taskId"`
	TaskBidderId   string `json:"taskBidderId"`
	TaskBidderRate int    `json:"taskBidderRate"`
}

type NewBidResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func TaskBid(w http.ResponseWriter, r *http.Request) {
	defer error_util.ErrorHandle(w)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var newBidRequest = NewBidRequest{}
	var newBidResponse NewBidResponse
	var err error
	fmt.Printf("newbid->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newBidRequest)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	fmt.Printf("newBidRequest:%v\n", newBidRequest)
	if strings.TrimSpace(newBidRequest.TaskId) != "" && strings.TrimSpace(newBidRequest.TaskBidderId) != "" && newBidRequest.TaskBidderRate >= 0 {
		_, err = db.Db.Exec("INSERT INTO task_bid (task_id,task_bidder_id,task_bidder_rate) VALUES(?,?,?)", newBidRequest.TaskId, newBidRequest.TaskBidderId, newBidRequest.TaskBidderRate)
		if err != nil {
			log.Println(err)
			goto Label0
		}
		newBidResponse.Status = "success"
		newBidResponse.Msg = "new bidder added success"
	} else {
		newBidResponse.Status = "error"
		newBidResponse.Msg = "new bidder added fail"
		goto Label1
	}

Label0:
	if newBidResponse.Status != "success" {
		newBidResponse.Status = "error"
		newBidResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(newBidResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
