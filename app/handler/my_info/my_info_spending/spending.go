package my_info_spending

import (
	"net/http"
	"strings"
)
func Spending(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	chartType := strings.TrimSpace(r.URL.Query().Get("chart-type"))
	if chartType == "card" {
		Card(w, r)
	} else if chartType == "summary" {
		Summary(w, r)
	} else if chartType == "datasource" {
		DataSource(w, r)
	}
}
