package my_info_earning

import (
	"net/http"
	"strings"
)

func Earning(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	chartType := strings.TrimSpace(r.URL.Query().Get("chart-type"))
	if chartType == "card" {
		Card(w, r)
	} else if chartType == "radar" {
		Radar(w, r)
	} else if chartType == "datasource" {
		DataSource(w, r)
	}
}
