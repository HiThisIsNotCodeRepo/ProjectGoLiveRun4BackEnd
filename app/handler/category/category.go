package category

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"paotui.sg/app/db"
)

type CategoriesResponse struct {
	Status     string     `json:"status"`
	Msg        string     `json:"msg"`
	Categories []Category `json:"categories"`
}
type Category struct {
	Cid   int    `json:"cid"`
	Title string `json:"title"`
}

func Categories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var getCategoriesResponse CategoriesResponse
	var getAllRows *sql.Rows
	categories := make([]Category, 0)
	var err error
	fmt.Printf("categories->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	getAllRows, err = db.Db.Query("SELECT cid,title from category")
	defer func() {
		if getAllRows != nil {
			getAllRows.Close()
		}
	}()
	if err != nil {
		log.Println(err)
		goto Label0
	}
	if getAllRows != nil {
		for getAllRows.Next() {
			var cid int
			var title string
			err = getAllRows.Scan(&cid, &title)
			if err != nil {
				log.Println(err)
				goto Label0
			}
			fmt.Printf("cid:%v,title:%v\n", cid, title)
			categories = append(categories, Category{
				Cid:   cid,
				Title: title,
			})
		}
	}

	getCategoriesResponse.Status = "success"
	getCategoriesResponse.Msg = fmt.Sprintf("categories:%v", categories)
	getCategoriesResponse.Categories = categories
Label0:
	if getCategoriesResponse.Status != "success" {
		getCategoriesResponse.Status = "error"
		getCategoriesResponse.Msg = "server error"
	}
	//Label1:
	encodeErr := encoder.Encode(getCategoriesResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
