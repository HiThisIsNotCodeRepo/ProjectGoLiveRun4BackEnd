package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"time"
)

func main() {
	var err error
	var db *sql.DB
	db, err = sql.Open("mysql", "user:password@tcp(:3306)/paotui?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("error in database connection")
	}

	// Category
	_, err = db.Exec("INSERT INTO category (cid,title) VALUES(?,?)", "0", "Buy Necessity")
	if err != nil {
		fmt.Println(`error in db.Exec("INSERT INTO category (cid,title) VALUES(?,?)", "0", "Buy Necessity")`)
	}
	_, err = db.Exec("INSERT INTO category (cid,title) VALUES(?,?)", "1", "Send Document")
	if err != nil {
		fmt.Println(`error in db.Exec("INSERT INTO category (cid,title) VALUES(?,?)", "1", "Send Document")`)
	}
	_, err = db.Exec("INSERT INTO category (cid,title) VALUES(?,?)", "2", "Food Delivery")
	if err != nil {
		fmt.Println(`error in db.Exec("INSERT INTO category (cid,title) VALUES(?,?)", "2", "Food Delivery")`)
	}
	_, err = db.Exec("INSERT INTO category (cid,title) VALUES(?,?)", "3", "Other")
	if err != nil {
		fmt.Println(`error in db.Exec("INSERT INTO category (cid,title) VALUES(?,?)", "3", "Other")`)
	}

	//	User
	for i := 1; i <= 5; i++ {
		uid := uuid.NewV4().String()
		name := fmt.Sprintf("user%d", i)
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("user%d", i)))
		password := base64.StdEncoding.EncodeToString(h.Sum(nil))
		email := fmt.Sprintf("user%d@email.com", i)
		mobileNumber := 84994075
		lastLogin := time.Now().Add(time.Hour * time.Duration(-i))
		fmt.Println(lastLogin)
		_, err = db.Exec("INSERT INTO user (uid,name,password,email,mobile_number,last_login) VALUES(?,?,?,?,?,?)", uid, name, password, email, mobileNumber, lastLogin)
	}
}
