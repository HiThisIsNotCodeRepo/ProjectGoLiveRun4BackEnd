package user

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"paotui.sg/app/db"
	"strings"
	"time"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	var registerRequest = RegisterRequest{}
	var registerResponse RegisterResponse
	var uid string
	var err error
	var lastLogin time.Time
	var name string
	var verifyCount int
	var email string
	var hash []byte
	fmt.Printf("register->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&registerRequest)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	fmt.Printf("Register Request:%v\n", registerRequest)
	if strings.TrimSpace(registerRequest.Name) != "" && strings.TrimSpace(registerRequest.Password) != "" && strings.TrimSpace(registerRequest.Email) != "" {
		name = registerRequest.Name
		err = db.Db.QueryRow("SELECT count(*) FROM user WHERE name = ?", name).Scan(&verifyCount)
		if err != nil {
			log.Println(err)
			goto Label0
		}
		if verifyCount != 0 {
			registerResponse.Status = "error"
			registerResponse.Msg = "user name duplicate"
			goto Label1
		}
		uid = uuid.NewV4().String()
		email = registerRequest.Email
		hash, err = bcrypt.GenerateFromPassword([]byte(name), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			goto Label0
		}
		password := string(hash)
		lastLogin = time.Now().Add(time.Hour)
		_, err = db.Db.Exec("INSERT INTO user (uid,name,password,email,last_login) VALUES(?,?,?,?,?)", uid, name, password, email, lastLogin)
		if err != nil {
			log.Println(err)
			goto Label0
		}

		registerResponse.Status = "success"
		registerResponse.Msg = "user register success"
	} else {
		registerResponse.Status = "error"
		registerResponse.Msg = "user register data error"
		goto Label1
	}

Label0:
	if registerResponse.Status != "success" {
		registerResponse.Status = "error"
		registerResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(registerResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
