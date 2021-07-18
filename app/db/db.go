package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"time"
)

var Db *sql.DB

//mysqlPassWord123
func init() {
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	address := strings.TrimSpace(os.Getenv("DB_ADDRESS"))
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/paotui?charset=utf8&parseTime=True&loc=Local", user, password, address))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	Db = db
	go cleanExpireTask()
}
