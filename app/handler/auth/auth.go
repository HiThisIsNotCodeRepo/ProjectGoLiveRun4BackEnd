package auth

import "net/http"

func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	chartType := r.URL.Query().Get("option")
	if chartType == "login" {
		Login(w, r)
	} else if chartType == "token-verify" {
		TokenVerify(w, r)
	}
}
