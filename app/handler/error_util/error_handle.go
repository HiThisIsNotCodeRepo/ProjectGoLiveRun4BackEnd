package error_util

import (
	"encoding/json"
	"log"
	"net/http"
)
type ErrorResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}
func ErrorHandle(w http.ResponseWriter) {
	if err := recover(); err != nil {
		log.Println(err)
		var errorResponse ErrorResponse
		encoder := json.NewEncoder(w)
		errorResponse.Status = "error"
		errorResponse.Msg = "server error"
		encodeErr := encoder.Encode(errorResponse)
		if encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
