package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var Db *sql.DB
//mysqlPassWord123
func init() {
	db, err := sql.Open("mysql", "user:password@tcp(192.168.0.150:3306)/paotui?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	Db = db
	go cleanExpireTask()
}

