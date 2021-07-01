package db

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

func cleanExpireTask() {
	defer autoRecover()
	var getAllRows *sql.Rows
	var err error
	var taskId string
	var tx *sql.Tx
	for {
		getAllRows, err = Db.Query("SELECT task_id FROM task WHERE task_start < NOW() AND task_step = 0")
		if err != nil {
			log.Printf("1%v\n", err)
		}
		if getAllRows != nil {
			for getAllRows.Next() {
				err = getAllRows.Scan(&taskId)
				if err != nil {
					log.Printf("2%v\n", err)
					goto Label0
				}
				tx, err = Db.Begin()
				if err != nil {
					log.Printf("3%v\n", err)
					goto Label0
				}
				_, err = tx.Exec("DELETE FROM task WHERE task_id = ? ", taskId)
				if err != nil {
					log.Printf("4%v\n", err)
					err = tx.Rollback()
					if err != nil {
						log.Printf("5%v\n", err)
					}
					goto Label0
				}
				_, err = tx.Exec("DELETE FROM task_bid WHERE task_id = ? ", taskId)
				if err != nil {
					log.Printf("6%v\n", err)
					err = tx.Rollback()
					if err != nil {
						log.Printf("7%v\n", err)
					}
					goto Label0
				}
				err = tx.Commit()
				if err != nil {
					log.Printf("8%v\n", err)
				}
			}

		}
	Label0:
		time.Sleep(5 * time.Second)
		if getAllRows != nil {
			err = getAllRows.Close()
			if err != nil {
				if !strings.Contains(err.Error(), "nil") {
					log.Printf("9%v\n", err)
				}
			}
		}
	}
}

func autoRecover() {
	if err := recover(); err != nil {
		log.Println(err)
		cleanExpireTask()
	}
}
