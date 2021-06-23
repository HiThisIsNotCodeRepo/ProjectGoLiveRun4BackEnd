package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"math"
	"math/rand"
	"strings"
	"time"
)

func main() {
	var err error
	var db *sql.DB
	db, err = sql.Open("mysql", "user:password@tcp(:3306)/paotui?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("error in database connection")
	}
	clearDB(db, err)
	categoryInit(db, err)
	userInit(db, err)
	taskInit(db, err)
}
func clearDB(db *sql.DB, err error) {
	_, err = db.Exec("DELETE FROM category")
	if err != nil {
		fmt.Println(`error in db.Exec("DELETE FROM category")`)
	}
	_, err = db.Exec("DELETE FROM task")
	if err != nil {
		fmt.Println(`error in db.Exec("DELETE FROM task")`)
	}
	_, err = db.Exec("DELETE FROM task_bid")
	if err != nil {
		fmt.Println(`error in db.Exec("DELETE FROM task_bid")`)
	}
	_, err = db.Exec("DELETE FROM user")
	if err != nil {
		fmt.Println(`error in db.Exec("DELETE FROM user")`)
	}
}
func categoryInit(db *sql.DB, err error) {
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
}
func userInit(db *sql.DB, err error) {
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
		_, err = db.Exec("INSERT INTO user (uid,name,password,email,mobile_number,last_login) VALUES(?,?,?,?,?,?)", uid, name, password, email, mobileNumber, lastLogin)
		if err != nil {
			fmt.Println(`db.Exec("INSERT INTO user (uid,name,password,email,mobile_number,last_login) VALUES(?,?,?,?,?,?)", uid, name, password, email, mobileNumber, lastLogin)`)
		}
	}
}

type Category struct {
	cid   int
	title string
}
type Location struct {
	locationName string
	index        int
}
type User struct {
	uid  string
	name string
}

func taskInit(db *sql.DB, err error) {
	var userArr = new([]User)
	var categoryArr = new([]Category)
	var addressArr = new([]Location)
	populateAddress(addressArr)
	getAllUidRow, getAllUidErr := db.Query("SELECT uid,name FROM user")
	if getAllUidErr != nil {
		fmt.Println(`error in db.Query("SELECT uid FROM user")`)
	}
	if getAllUidRow == nil {
		fmt.Println(`getAllUidRow == nil`)
	} else {
		for getAllUidRow.Next() {
			var user User
			err = getAllUidRow.Scan(&user.uid, &user.name)
			if err != nil {
				fmt.Println(`error in getAllUidRow.Scan(&uid)`)
			}
			*userArr = append(*userArr, user)

		}
	}
	getAllCategoryRow, getAllCategoryErr := db.Query("SELECT cid,title FROM category")
	if getAllCategoryErr != nil {
		fmt.Println(`error in db.Query("SELECT cid,title FROM user")`)
	}
	if getAllCategoryRow == nil {
		fmt.Println(`getAllCateogryRow == nil`)
	} else {
		for getAllCategoryRow.Next() {
			var category Category
			err = getAllCategoryRow.Scan(&category.cid, &category.title)
			if err != nil {
				fmt.Println(`error in getAllCategoryRow.Scan(&category.cid, &category.title)`)
			}
			*categoryArr = append(*categoryArr, category)

		}
	}
	currentDayOfWeek := int(time.Now().Weekday())
	for i := -(currentDayOfWeek + 6); i < 0; i++ {
		tempTime := time.Now().AddDate(0, 0, i)
		for _, v := range *userArr {
			for j := 0; j < 5; j++ {
				tempUidArr := new([]string)
				for _, tempUser := range *userArr {
					if v.uid != tempUser.uid {
						*tempUidArr = append(*tempUidArr, tempUser.uid)
					}
				}
				taskId := uuid.NewV4().String()
				taskIdLast4Char := getLast4Char(taskId)
				rand.Seed(time.Now().UnixNano())
				categoryLen := len(*categoryArr)
				randIndexForCategory := rand.Intn(categoryLen)
				taskTitle := fmt.Sprintf("%s-%s", (*categoryArr)[randIndexForCategory].title, taskIdLast4Char)
				taskDescription := fmt.Sprintf("%s for %s taskId:%s", (*categoryArr)[randIndexForCategory].title, v.name, taskId)
				addressLen := len(*addressArr)
				indexForFrom := rand.Intn(addressLen)
				fromLocation := (*addressArr)[indexForFrom].locationName
				indexForTo := rand.Intn(addressLen)
				toLocation := (*addressArr)[indexForTo].locationName
				createTime := tempTime.Add(time.Hour * (-2))
				startTime := tempTime.Add(time.Hour * (-1))
				difference := (*addressArr)[indexForFrom].index - (*addressArr)[indexForTo].index
				if difference < 0 {
					difference = -1 * difference
				}
				duration := 15 + difference*3
				completeTime := startTime.Add(time.Minute * time.Duration(duration))
				taskStep := 2
				ownerRate := math.Ceil(float64(duration) * 0.7)
				indexForDeliverId := rand.Intn(len(*tempUidArr))
				taskDeliverId := (*tempUidArr)[indexForDeliverId]
				deliverRate := math.Ceil(ownerRate * 0.7)
				fmt.Printf("task_id:%v\ntask_title:%v\ntask_description:%v\ntask_category_id:%v\ntask_from:%v\ntask_to:%v\ntask_create:%v\ntask_start:%v\ntask_complete:%v\ntask_duration:%v\ntask_step:%v\ntask_ownder_id:%v\ntask_owner_rate:%v\ntask_deliver_id:%v\ntask_deliver_rate:%v\n", taskId, taskTitle, taskDescription, randIndexForCategory, fromLocation, toLocation, createTime, startTime, completeTime, duration, taskStep, v.uid, ownerRate, taskDeliverId, deliverRate)
				_, err = db.Exec("INSERT INTO task (task_id,task_title,task_description,task_category_id,task_from,task_to,task_create,task_start,task_complete,task_duration,task_step,task_owner_id,task_owner_rate,task_deliver_id,task_deliver_rate) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", taskId, taskTitle, taskDescription, randIndexForCategory, fromLocation, toLocation, createTime, startTime, completeTime, duration, taskStep, v.uid, ownerRate, taskDeliverId, deliverRate)
				if err != nil {
					fmt.Println(`error in db.Exec("INSERT INTO task (task_id,task_title,task_description,task_category_id,task_from,task_to,task_create,task_start,task_complete,task_duration,task_step,task_owner_id,task_owner_rate,task_deliver_id,task_deliver_rate) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", taskId, taskTitle, taskDescription, randIndexForCategory, fromLocation, toLocation,createTime,startTime,completeTime,taskStep,taskStep, v.uid,ownerRate,taskDeliverId,deliverRate)`)
				}
			}

		}
	}

	//for _, v := range *userArr {
	//	fmt.Printf("uid:%v,name:%v\n", v.uid, v.name)
	//}
	//for _, v := range *categoryArr {
	//	fmt.Println(v)
	//}
	//for _, v := range *addressArr {
	//	fmt.Println(v)
	//}
}

func populateAddress(arr *[]Location) {
	*arr = append(*arr, Location{"Pasir Ris", 1})
	*arr = append(*arr, Location{"Tampines", 2})
	*arr = append(*arr, Location{"Simei", 3})
	*arr = append(*arr, Location{"Tanah Merah", 4})
	*arr = append(*arr, Location{"Bedok", 5})
	*arr = append(*arr, Location{"Kembangan", 6})
	*arr = append(*arr, Location{"Eunos", 7})
	*arr = append(*arr, Location{"Paya Lebar", 8})
	*arr = append(*arr, Location{"Aljunied", 9})
	*arr = append(*arr, Location{"Kallang", 10})
	*arr = append(*arr, Location{"Lavender", 11})
	*arr = append(*arr, Location{"Bugis", 12})
	*arr = append(*arr, Location{"City Hall", 13})
	*arr = append(*arr, Location{"Raffles Place", 14})
	*arr = append(*arr, Location{"Tanjong Pagar", 15})
	*arr = append(*arr, Location{"Outram Park", 16})
	*arr = append(*arr, Location{"Tiong Bahru", 17})
	*arr = append(*arr, Location{"Redhill", 18})
	*arr = append(*arr, Location{"Queenstown", 19})
	*arr = append(*arr, Location{"Commonwealth", 20})
	*arr = append(*arr, Location{"Buona Vista", 21})
	*arr = append(*arr, Location{"Dover", 22})
	*arr = append(*arr, Location{"Clementi", 23})
	*arr = append(*arr, Location{"Jurong East", 24})
	*arr = append(*arr, Location{"Chinese Garden", 25})
	*arr = append(*arr, Location{"Lakeside", 26})
	*arr = append(*arr, Location{"Boon Lay", 27})
	*arr = append(*arr, Location{"Pioneer", 28})
	*arr = append(*arr, Location{"Joo Koon", 29})
	*arr = append(*arr, Location{"Gul Circle", 30})
	*arr = append(*arr, Location{"Tuas Crescent", 31})
	*arr = append(*arr, Location{"Tuas West Road", 32})
	*arr = append(*arr, Location{"Tuas Link", 33})
}
func getLast4Char(str string) string {
	strArr := strings.Split(str, "-")
	lastStrByteArr := []byte(strArr[len(strArr)-1])
	var last4ByteArr []byte
	for i := len(lastStrByteArr) - 4; i < len(lastStrByteArr); i++ {
		last4ByteArr = append(last4ByteArr, lastStrByteArr[i])
	}
	return string(last4ByteArr)
}
