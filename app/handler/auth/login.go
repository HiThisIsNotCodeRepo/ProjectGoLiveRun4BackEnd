package auth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"paotui.sg/app/db"
	"paotui.sg/app/jwt"
	"strings"
	"time"
)

type UserClaim struct {
	UserId         string
	ExpireDateTime int64
}

type UserLoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Status    string `json:"status"`
	Msg       string `json:"msg"`
	UserId    string `json:"userId"`
	Token     string `json:"token"`
	LastLogin string `json:"lastLogin"`
	Email     string `json:"email"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var userLoginRequest = UserLoginRequest{}
	var userLoginResponse UserLoginResponse
	var storedPassword string
	var uid string
	var err error
	var lastLogin string
	var userClaim []byte
	var encryptor jose.Encrypter
	var object *jose.JSONWebEncryption
	var token string
	var email string
	fmt.Printf("request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&userLoginRequest)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	fmt.Printf("userLoginRequest:%v\n", userLoginRequest)
	if strings.TrimSpace(userLoginRequest.Name) != "" && strings.TrimSpace(userLoginRequest.Password) != "" {
		err = db.Db.QueryRow("SELECT uid, password,last_Login,email FROM user WHERE name = ?", userLoginRequest.Name).Scan(&uid, &storedPassword, &lastLogin, &email)
		if err != nil {
			fmt.Println(err)
			goto Label0
		}
		if strings.TrimSpace(storedPassword) == "" || strings.TrimSpace(uid) == "" {
			userLoginResponse.Status = "error"
			userLoginResponse.Msg = "no record found for the user"
		}
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(userLoginRequest.Password))
		if err != nil {
			fmt.Println(err)
			userLoginResponse.Status = "error"
			userLoginResponse.Msg = "user password is not correct"
			goto Label1
		}
		publicKey := jwt.VerifyKey
		encryptor, err = jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: publicKey}, nil)
		if err != nil {
			goto Label0
		}
		userClaim, err = json.Marshal(&UserClaim{UserId: uid, ExpireDateTime: time.Now().Add(1 * time.Hour).Unix()})
		if err != nil {
			goto Label0
		}
		object, err = encryptor.Encrypt(userClaim)
		if err != nil {
			goto Label0
		}
		token, err = object.CompactSerialize()
		if err != nil {
			goto Label0
		}
		_, err = db.Db.Exec("UPDATE user SET last_login =? WHERE uid =? ", time.Now(), uid)
		if err != nil {
			goto Label0
		}
		fmt.Printf("token:%v\n", token)
		userLoginResponse.Status = "success"
		userLoginResponse.Msg = "user login success"
		userLoginResponse.Token = token
		userLoginResponse.UserId = uid
		userLoginResponse.LastLogin = lastLogin
		userLoginResponse.Email = email
	} else {
		userLoginResponse.Status = "error"
		userLoginResponse.Msg = "user login data error"
		goto Label1
	}

Label0:
	if userLoginResponse.Status != "success" {
		userLoginResponse.Status = "error"
		userLoginResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(userLoginResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
